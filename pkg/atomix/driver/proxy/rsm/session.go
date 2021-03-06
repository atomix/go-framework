// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rsm

import (
	"context"
	"github.com/atomix/atomix-go-framework/pkg/atomix/cluster"
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	"github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	streams "github.com/atomix/atomix-go-framework/pkg/atomix/stream"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"io"
	"math"
	"math/rand"
	"sync"
	"time"
)

// SessionOption implements a session option
type SessionOption interface {
	prepare(options *sessionOptions)
}

// WithSessionTimeout returns a session SessionOption to configure the session timeout
func WithSessionTimeout(timeout time.Duration) SessionOption {
	return sessionTimeoutOption{timeout: timeout}
}

type sessionTimeoutOption struct {
	timeout time.Duration
}

func (o sessionTimeoutOption) prepare(options *sessionOptions) {
	options.timeout = o.timeout
}

type sessionOptions struct {
	id      string
	timeout time.Duration
	retry   time.Duration
}

// NewSession creates a new Session for the given partition
// name is the name of the primitive
// handler is the primitive's session handler
func NewSession(partition cluster.Partition, log logging.Logger, opts ...SessionOption) *Session {
	options := &sessionOptions{
		id:      uuid.New().String(),
		retry:   15 * time.Second,
		timeout: time.Minute,
	}
	for i := range opts {
		opts[i].prepare(options)
	}
	return &Session{
		partition:     partition,
		Timeout:       options.timeout,
		retryInterval: options.retry,
		streams:       make(map[uint64]*StreamState),
		log:           log,
		mu:            sync.RWMutex{},
		ticker:        time.NewTicker(options.timeout / 4),
	}
}

// Session maintains the session for a primitive
type Session struct {
	partition     cluster.Partition
	Timeout       time.Duration
	SessionID     uint64
	lastIndex     uint64
	requestID     uint64
	responseID    uint64
	streams       map[uint64]*StreamState
	log           logging.Logger
	retryInterval time.Duration
	conn          *grpc.ClientConn
	leader        *cluster.Replica
	mu            sync.RWMutex
	ticker        *time.Ticker
}

// DoCommand submits a command to the service
func (s *Session) DoCommand(ctx context.Context, service rsm.ServiceId, name string, input []byte) ([]byte, error) {
	requestContext := s.nextCommandContext()
	response, responseStatus, responseContext, err := s.doCommand(ctx, name, input, service, requestContext)
	if err != nil {
		return nil, err
	}
	if responseStatus.Code != rsm.SessionResponseCode_OK {
		return nil, rsm.GetErrorFromStatus(responseStatus)
	}
	s.recordCommandResponse(requestContext, responseContext)
	return response, nil
}

// doCommand submits a command to the service
func (s *Session) doCommand(ctx context.Context, name string, input []byte, service rsm.ServiceId, command rsm.SessionCommandContext) ([]byte, rsm.SessionResponseStatus, rsm.SessionResponseContext, error) {
	request := &rsm.StorageRequest{
		PartitionID: uint32(s.partition.ID()),
		Request: &rsm.SessionRequest{
			Request: &rsm.SessionRequest_Command{
				Command: &rsm.SessionCommandRequest{
					Context: command,
					Command: rsm.ServiceCommandRequest{
						Service: service,
						Request: &rsm.ServiceCommandRequest_Operation{
							Operation: &rsm.ServiceOperationRequest{
								Method: name,
								Value:  input,
							},
						},
					},
				},
			},
		},
	}

	response, err := s.doRequest(ctx, request)
	if err != nil {
		return nil, rsm.SessionResponseStatus{}, rsm.SessionResponseContext{}, err
	}
	return response.Response.GetCommand().Response.GetOperation().Result, response.Response.Status, response.Response.GetCommand().GetContext(), err
}

// DoCommandStream submits a streaming command to the service
func (s *Session) DoCommandStream(ctx context.Context, service rsm.ServiceId, name string, input []byte, outStream streams.WriteStream) error {
	streamState, requestContext := s.nextStream()
	ch := make(chan streams.Result)
	inStream := streams.NewChannelStream(ch)
	err := s.doCommandStream(context.Background(), name, input, service, requestContext, inStream)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case result, ok := <-ch:
				if !ok {
					s.deleteStream(streamState.ID)
					return
				}

				if result.Failed() {
					s.deleteStream(streamState.ID)
					outStream.Error(result.Error)
					outStream.Close()
					return
				}

				response := result.Value.(PartitionOutput)
				switch response.Type {
				case rsm.SessionResponseType_OPEN_STREAM:
					if streamState.serialize(response.Context) {
						outStream.Send(response.Result)
					}
				case rsm.SessionResponseType_CLOSE_STREAM:
					if streamState.serialize(response.Context) {
						outStream.Close()
						streamState.Close()
						return
					}
				case rsm.SessionResponseType_RESPONSE:
					// Record the response
					s.recordCommandResponse(requestContext, response.Context)

					// Attempt to serialize the response to the stream and skip the response if serialization failed.
					if streamState.serialize(response.Context) {
						outStream.Send(response.Result)
					}
				}
			case <-ctx.Done():
				s.deleteStream(streamState.ID)
				outStream.Error(ctx.Err())
				outStream.Close()
				return
			}
		}
	}()
	return nil
}

// doCommandStream submits a streaming command to the service
func (s *Session) doCommandStream(ctx context.Context, name string, input []byte, service rsm.ServiceId, context rsm.SessionCommandContext, stream streams.WriteStream) error {
	request := &rsm.StorageRequest{
		PartitionID: uint32(s.partition.ID()),
		Request: &rsm.SessionRequest{
			Request: &rsm.SessionRequest_Command{
				Command: &rsm.SessionCommandRequest{
					Context: context,
					Command: rsm.ServiceCommandRequest{
						Service: service,
						Request: &rsm.ServiceCommandRequest_Operation{
							Operation: &rsm.ServiceOperationRequest{
								Method: name,
								Value:  input,
							},
						},
					},
				},
			},
		},
	}
	return s.doStream(ctx, request, streams.NewDecodingStream(stream, func(value interface{}, err error) (interface{}, error) {
		if err != nil {
			return nil, err
		}
		response := value.(*rsm.StorageResponse)
		commandResponse := response.Response.GetCommand()
		var result []byte
		if commandResponse.Response.GetOperation() != nil {
			result = commandResponse.Response.GetOperation().Result
		}
		return PartitionOutput{
			Type:    response.Response.Type,
			Status:  response.Response.Status,
			Context: commandResponse.Context,
			Result: streams.Result{
				Value: result,
			},
		}, nil
	}), true)
}

// DoQuery submits a query to the service
func (s *Session) DoQuery(ctx context.Context, service rsm.ServiceId, name string, input []byte, sync bool) ([]byte, error) {
	requestContext := s.getQueryContext(sync)
	response, responseStatus, responseContext, err := s.doQuery(ctx, name, input, service, requestContext)
	if err != nil {
		return nil, err
	}
	if responseStatus.Code != rsm.SessionResponseCode_OK {
		return nil, rsm.GetErrorFromStatus(responseStatus)
	}
	s.recordQueryResponse(requestContext, responseContext)
	return response, nil
}

// doQuery submits a query to the service
func (s *Session) doQuery(ctx context.Context, name string, input []byte, service rsm.ServiceId, query rsm.SessionQueryContext) ([]byte, rsm.SessionResponseStatus, rsm.SessionResponseContext, error) {
	request := &rsm.StorageRequest{
		PartitionID: uint32(s.partition.ID()),
		Request: &rsm.SessionRequest{
			Request: &rsm.SessionRequest_Query{
				Query: &rsm.SessionQueryRequest{
					Context: query,
					Query: rsm.ServiceQueryRequest{
						Service: &service,
						Request: &rsm.ServiceQueryRequest_Operation{
							Operation: &rsm.ServiceOperationRequest{
								Method: name,
								Value:  input,
							},
						},
					},
				},
			},
		},
	}

	response, err := s.doRequest(ctx, request)
	if err != nil {
		return nil, rsm.SessionResponseStatus{}, rsm.SessionResponseContext{}, err
	}
	return response.Response.GetQuery().Response.GetOperation().Result, response.Response.Status, response.Response.GetQuery().GetContext(), err
}

// DoQueryStream submits a streaming query to the service
func (s *Session) DoQueryStream(ctx context.Context, service rsm.ServiceId, name string, input []byte, outStream streams.WriteStream, sync bool) error {
	requestContext := s.getQueryContext(sync)
	ch := make(chan streams.Result)
	inStream := streams.NewChannelStream(ch)
	err := s.doQueryStream(context.Background(), name, input, service, requestContext, inStream)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case result, ok := <-ch:
				if !ok {
					return
				}

				if result.Failed() {
					outStream.Error(result.Error)
					outStream.Close()
					return
				}

				response := result.Value.(PartitionOutput)
				switch response.Type {
				case rsm.SessionResponseType_OPEN_STREAM:
				case rsm.SessionResponseType_CLOSE_STREAM:
					outStream.Close()
					return
				case rsm.SessionResponseType_RESPONSE:
					s.recordQueryResponse(requestContext, response.Context)
					outStream.Send(response.Result)
				}
			case <-ctx.Done():
				outStream.Error(ctx.Err())
				outStream.Close()
				return
			}
		}
	}()
	return nil
}

// doQueryStream submits a streaming query to the service
func (s *Session) doQueryStream(ctx context.Context, name string, input []byte, service rsm.ServiceId, context rsm.SessionQueryContext, stream streams.WriteStream) error {
	request := &rsm.StorageRequest{
		PartitionID: uint32(s.partition.ID()),
		Request: &rsm.SessionRequest{
			Request: &rsm.SessionRequest_Query{
				Query: &rsm.SessionQueryRequest{
					Context: context,
					Query: rsm.ServiceQueryRequest{
						Service: &service,
						Request: &rsm.ServiceQueryRequest_Operation{
							Operation: &rsm.ServiceOperationRequest{
								Method: name,
								Value:  input,
							},
						},
					},
				},
			},
		},
	}
	return s.doStream(ctx, request, streams.NewDecodingStream(stream, func(value interface{}, err error) (interface{}, error) {
		if err != nil {
			return nil, err
		}
		response := value.(*rsm.StorageResponse)
		queryResponse := response.Response.GetQuery()
		var result []byte
		if queryResponse.Response.GetOperation() != nil {
			result = queryResponse.Response.GetOperation().Result
		}
		return PartitionOutput{
			Type:    response.Response.Type,
			Status:  response.Response.Status,
			Context: queryResponse.Context,
			Result: streams.Result{
				Value: result,
			},
		}, nil
	}), false)
}

// doMetadata submits a metadata query to the service
func (s *Session) doMetadata(ctx context.Context, serviceType string, namespace string, context rsm.SessionQueryContext) ([]*rsm.ServiceId, rsm.SessionResponseStatus, rsm.SessionResponseContext, error) {
	request := &rsm.StorageRequest{
		PartitionID: uint32(s.partition.ID()),
		Request: &rsm.SessionRequest{
			Request: &rsm.SessionRequest_Query{
				Query: &rsm.SessionQueryRequest{
					Context: context,
					Query: rsm.ServiceQueryRequest{
						Request: &rsm.ServiceQueryRequest_Metadata{
							Metadata: &rsm.ServiceMetadataRequest{
								Type: serviceType,
							},
						},
					},
				},
			},
		},
	}
	response, err := s.doRequest(ctx, request)
	if err != nil {
		return nil, rsm.SessionResponseStatus{}, rsm.SessionResponseContext{}, err
	}
	return response.Response.GetQuery().Response.GetMetadata().Services, response.Response.Status, response.Response.GetQuery().Context, nil
}

// DoCreateService creates the service
func (s *Session) DoCreateService(ctx context.Context, service rsm.ServiceId) error {
	requestContext := s.nextCommandContext()
	responseStatus, responseContext, err := s.doCreateService(ctx, service, requestContext)
	if err != nil {
		return err
	}
	if responseStatus.Code != rsm.SessionResponseCode_OK {
		return rsm.GetErrorFromStatus(responseStatus)
	}
	s.recordCommandResponse(requestContext, responseContext)
	return nil
}

// doCreateService creates the service
func (s *Session) doCreateService(ctx context.Context, service rsm.ServiceId, context rsm.SessionCommandContext) (rsm.SessionResponseStatus, rsm.SessionResponseContext, error) {
	request := &rsm.StorageRequest{
		PartitionID: uint32(s.partition.ID()),
		Request: &rsm.SessionRequest{
			Request: &rsm.SessionRequest_Command{
				Command: &rsm.SessionCommandRequest{
					Context: context,
					Command: rsm.ServiceCommandRequest{
						Service: service,
						Request: &rsm.ServiceCommandRequest_Create{
							Create: &rsm.ServiceCreateRequest{},
						},
					},
				},
			},
		},
	}
	response, err := s.doRequest(ctx, request)
	if err != nil {
		return rsm.SessionResponseStatus{}, rsm.SessionResponseContext{}, err
	}
	return response.Response.Status, response.Response.GetCommand().Context, nil
}

// DoCloseService closes the service
func (s *Session) DoCloseService(ctx context.Context, service rsm.ServiceId) error {
	requestContext := s.nextCommandContext()
	responseStatus, responseContext, err := s.doCloseService(ctx, service, requestContext)
	if err != nil {
		return err
	}
	if responseStatus.Code != rsm.SessionResponseCode_OK {
		return rsm.GetErrorFromStatus(responseStatus)
	}
	s.recordCommandResponse(requestContext, responseContext)
	return nil
}

// doCloseService closes the service
func (s *Session) doCloseService(ctx context.Context, service rsm.ServiceId, context rsm.SessionCommandContext) (rsm.SessionResponseStatus, rsm.SessionResponseContext, error) {
	request := &rsm.StorageRequest{
		PartitionID: uint32(s.partition.ID()),
		Request: &rsm.SessionRequest{
			Request: &rsm.SessionRequest_Command{
				Command: &rsm.SessionCommandRequest{
					Context: context,
					Command: rsm.ServiceCommandRequest{
						Service: service,
						Request: &rsm.ServiceCommandRequest_Close{
							Close: &rsm.ServiceCloseRequest{},
						},
					},
				},
			},
		},
	}
	response, err := s.doRequest(ctx, request)
	if err != nil {
		return rsm.SessionResponseStatus{}, rsm.SessionResponseContext{}, err
	}
	return response.Response.Status, response.Response.GetCommand().Context, nil
}

// DoDeleteService deletes the service
func (s *Session) DoDeleteService(ctx context.Context, service rsm.ServiceId) error {
	requestContext := s.nextCommandContext()
	responseStatus, responseContext, err := s.doDeleteService(ctx, service, requestContext)
	if err != nil {
		return err
	}
	if responseStatus.Code != rsm.SessionResponseCode_OK {
		return rsm.GetErrorFromStatus(responseStatus)
	}
	s.recordCommandResponse(requestContext, responseContext)
	return nil
}

// doDeleteService deletes the service
func (s *Session) doDeleteService(ctx context.Context, service rsm.ServiceId, context rsm.SessionCommandContext) (rsm.SessionResponseStatus, rsm.SessionResponseContext, error) {
	request := &rsm.StorageRequest{
		PartitionID: uint32(s.partition.ID()),
		Request: &rsm.SessionRequest{
			Request: &rsm.SessionRequest_Command{
				Command: &rsm.SessionCommandRequest{
					Context: context,
					Command: rsm.ServiceCommandRequest{
						Service: service,
						Request: &rsm.ServiceCommandRequest_Delete{
							Delete: &rsm.ServiceDeleteRequest{},
						},
					},
				},
			},
		},
	}
	response, err := s.doRequest(ctx, request)
	if err != nil {
		return rsm.SessionResponseStatus{}, rsm.SessionResponseContext{}, err
	}
	return response.Response.Status, response.Response.GetCommand().Context, nil
}

// open creates the session and begins keep-alives
func (s *Session) open(ctx context.Context) error {
	requestContext, _ := s.getStateContexts()
	responseStatus, responseContext, err := s.doOpenSession(ctx, requestContext, &s.Timeout)
	if err != nil {
		return err
	}

	if responseStatus.Code != rsm.SessionResponseCode_OK {
		return rsm.GetErrorFromStatus(responseStatus)
	}

	s.mu.Lock()
	s.SessionID = responseContext.SessionID
	s.lastIndex = responseContext.SessionID
	s.mu.Unlock()

	go func() {
		for range s.ticker.C {
			go s.keepAlive(context.Background())
		}
	}()
	return nil
}

// doOpenSession opens a new session
func (s *Session) doOpenSession(ctx context.Context, context rsm.SessionCommandContext, timeout *time.Duration) (rsm.SessionResponseStatus, rsm.SessionResponseContext, error) {
	request := &rsm.StorageRequest{
		PartitionID: uint32(s.partition.ID()),
		Request: &rsm.SessionRequest{
			Request: &rsm.SessionRequest_OpenSession{
				OpenSession: &rsm.OpenSessionRequest{
					Timeout: timeout,
				},
			},
		},
	}
	response, err := s.doRequest(ctx, request)
	if err != nil {
		return rsm.SessionResponseStatus{}, rsm.SessionResponseContext{}, err
	}
	sessionID := response.Response.GetOpenSession().SessionID
	return response.Response.Status, rsm.SessionResponseContext{
		SessionID: sessionID,
		RequestID: sessionID,
		Index:     sessionID,
	}, nil
}

// keepAlive keeps the session alive
func (s *Session) keepAlive(ctx context.Context) error {
	requestContext, streamContexts := s.getStateContexts()
	responseStatus, responseContext, err := s.doKeepAliveSession(ctx, requestContext, streamContexts)
	if err != nil {
		return err
	}

	if responseStatus.Code != rsm.SessionResponseCode_OK {
		return rsm.GetErrorFromStatus(responseStatus)
	}

	s.recordCommandResponse(requestContext, responseContext)
	return nil
}

// doKeepAliveSession keeps a session alive
func (s *Session) doKeepAliveSession(ctx context.Context, context rsm.SessionCommandContext, streams []rsm.SessionStreamContext) (rsm.SessionResponseStatus, rsm.SessionResponseContext, error) {
	request := &rsm.StorageRequest{
		PartitionID: uint32(s.partition.ID()),
		Request: &rsm.SessionRequest{
			Request: &rsm.SessionRequest_KeepAlive{
				KeepAlive: &rsm.KeepAliveRequest{
					SessionID:    context.SessionID,
					AckRequestID: context.RequestID,
					Streams:      streams,
				},
			},
		},
	}
	response, err := s.doRequest(ctx, request)
	if err != nil {
		return rsm.SessionResponseStatus{}, rsm.SessionResponseContext{}, err
	}
	return response.Response.Status, rsm.SessionResponseContext{SessionID: context.SessionID}, nil
}

// doRequest submits a storage request
func (s *Session) doRequest(ctx context.Context, request *rsm.StorageRequest) (*rsm.StorageResponse, error) {
	i := 1
	for {
		s.log.Debugf("Sending StorageRequest %+v", request)
		requestCtx, _ := context.WithTimeout(ctx, s.retryInterval)
		response, err := s.tryRequest(requestCtx, request)
		if err == nil {
			s.log.Debugf("Received StorageResponse %+v", response)
			switch response.Response.Status.Code {
			case rsm.SessionResponseCode_OK:
				return response, err
			case rsm.SessionResponseCode_NOT_LEADER:
				if response.Response.Status.Leader != "" {
					s.log.Debugf("Reconnecting to leader %s", response.Response.Status.Leader)
					s.reconnect(cluster.ReplicaID(response.Response.Status.Leader))
				} else {
					s.log.Debug("Failed to locate leader, retrying...")
					select {
					case <-time.After(10 * time.Millisecond * time.Duration(math.Min(math.Pow(2, float64(i)), 1000))):
						i++
					case <-ctx.Done():
						return nil, ctx.Err()
					}
				}
			default:
				return response, rsm.GetErrorFromStatus(response.Response.Status)
			}
		} else if errors.IsTimeout(err) {
			s.log.Warnf("StorageRequest %+v timed out. Retrying...", request, err)
		} else if errors.IsCanceled(err) {
			return nil, err
		} else {
			s.log.Warnf("Sending StorageRequest %+v failed: %s", request, err)
			select {
			case <-time.After(10 * time.Millisecond * time.Duration(math.Min(math.Pow(2, float64(i)), 1000))):
				i++
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
}

// tryRequest submits a storage request
func (s *Session) tryRequest(ctx context.Context, request *rsm.StorageRequest) (*rsm.StorageResponse, error) {
	conn, err := s.connect()
	if err != nil {
		return nil, err
	}
	client := rsm.NewStorageServiceClient(conn)
	response, err := client.Request(ctx, request)
	if err != nil {
		return nil, errors.From(err)
	}
	return response, nil
}

// doStream submits a streaming request to the service
func (s *Session) doStream(ctx context.Context, request *rsm.StorageRequest, stream streams.WriteStream, idempotent bool) error {
	go s.tryStream(ctx, request, stream, idempotent, 0, false)
	return nil
}

// tryStream submits a stream request to the service recursively
func (s *Session) tryStream(ctx context.Context, request *rsm.StorageRequest, stream streams.WriteStream, idempotent bool, attempts int, open bool) {
	conn, err := s.connect()
	if err == context.Canceled {
		stream.Error(err)
		stream.Close()
		return
	} else if err != nil {
		s.log.Warnf("StorageRequest %+v failed", request, err)
		if idempotent || !open {
			go s.retryStream(ctx, request, stream, idempotent, attempts, open)
		} else {
			stream.Error(errors.From(err))
			stream.Close()
		}
		return
	}
	client := rsm.NewStorageServiceClient(conn)

	s.log.Debugf("Sending StorageRequest %+v", request)
	responseStream, err := client.Stream(ctx, request)
	if err != nil {
		s.log.Warnf("Sending StorageRequest %+v failed: %s", request, err)
		stream.Error(err)
		stream.Close()
		return
	}

	for {
		response, err := responseStream.Recv()
		if err == io.EOF {
			stream.Close()
			return
		} else if err == context.Canceled {
			stream.Error(err)
			stream.Close()
			return
		} else if err != nil {
			s.log.Warnf("StorageRequest %+v failed", request, err)
			if idempotent || !open {
				go s.retryStream(ctx, request, stream, idempotent, attempts, open)
			} else {
				stream.Error(errors.From(err))
				stream.Close()
			}
			return
		} else {
			s.log.Debugf("Received StorageResponse %+v", response)
			switch response.Response.Status.Code {
			case rsm.SessionResponseCode_OK:
				stream.Value(response)
				switch response.Response.Type {
				case rsm.SessionResponseType_OPEN_STREAM:
					if !open {
						open = true
					}
				case rsm.SessionResponseType_CLOSE_STREAM:
					stream.Close()
					return
				}
			case rsm.SessionResponseCode_NOT_LEADER:
				if response.Response.Status.Leader != "" {
					s.log.Debugf("Reconnecting to leader %s", response.Response.Status.Leader)
					s.reconnect(cluster.ReplicaID(response.Response.Status.Leader))
				}
				go s.retryStream(ctx, request, stream, idempotent, attempts, open)
				return
			default:
				stream.Error(rsm.GetErrorFromStatus(response.Response.Status))
				stream.Close()
				return
			}
			attempts = 0
		}
	}
}

func (s *Session) retryStream(ctx context.Context, request *rsm.StorageRequest, stream streams.WriteStream, idempotent bool, attempts int, open bool) {
	attempts++
	select {
	case <-time.After(10 * time.Millisecond * time.Duration(math.Min(math.Pow(2, float64(attempts)), 1000))):
		s.tryStream(ctx, request, stream, idempotent, attempts, open)
	case <-ctx.Done():
		stream.Error(errors.NewTimeout(ctx.Err().Error()))
		stream.Close()
	}
}

// connect gets the connection to the service
func (s *Session) connect() (*grpc.ClientConn, error) {
	s.mu.RLock()
	conn := s.conn
	s.mu.RUnlock()
	if conn != nil {
		return conn, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	conn = s.conn
	if conn != nil {
		return conn, nil
	}

	if s.leader == nil {
		replicas := make([]*cluster.Replica, 0)
		for _, replica := range s.partition.Replicas() {
			replicas = append(replicas, replica)
		}
		s.leader = replicas[rand.Intn(len(replicas))]
		s.conn = nil
	}

	s.log.Infof("Connecting to partition %d replica %s", s.partition.ID(), s.leader.ID)
	conn, err := s.leader.Connect(context.Background(), cluster.WithDialOption(grpc.WithInsecure()))
	if err != nil {
		s.log.Warnf("Connecting to partition %d replica %s failed", s.partition.ID(), s.leader.ID, err)
		return nil, err
	}
	s.conn = conn
	return conn, nil
}

// reconnect the connection to the given leader
func (s *Session) reconnect(replica cluster.ReplicaID) {
	if replica == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.leader.ID == replica {
		return
	}

	leader, ok := s.partition.Replica(replica)
	if !ok {
		return
	}

	s.leader = leader
	s.conn = nil
}

// disconnect closes the connections
func (s *Session) disconnect() error {
	s.mu.Lock()
	s.conn = nil
	s.mu.Unlock()
	return nil
}

// Close closes the session
func (s *Session) Close() error {
	err := s.close(context.TODO())
	s.ticker.Stop()
	return err
}

// close closes the session
func (s *Session) close(ctx context.Context) error {
	requestContext, _ := s.getStateContexts()
	responseStatus, responseContext, err := s.doCloseSession(ctx, requestContext)
	if err != nil {
		return err
	}

	if responseStatus.Code != rsm.SessionResponseCode_OK {
		return rsm.GetErrorFromStatus(responseStatus)
	}

	s.recordCommandResponse(requestContext, responseContext)
	return nil
}

// doCloseSession closes a session
func (s *Session) doCloseSession(ctx context.Context, context rsm.SessionCommandContext) (rsm.SessionResponseStatus, rsm.SessionResponseContext, error) {
	request := &rsm.StorageRequest{
		PartitionID: uint32(s.partition.ID()),
		Request: &rsm.SessionRequest{
			Request: &rsm.SessionRequest_CloseSession{
				CloseSession: &rsm.CloseSessionRequest{
					SessionID: context.SessionID,
				},
			},
		},
	}
	response, err := s.doRequest(ctx, request)
	if err != nil {
		return rsm.SessionResponseStatus{}, rsm.SessionResponseContext{}, err
	}
	return response.Response.Status, rsm.SessionResponseContext{SessionID: context.SessionID}, nil
}

// getStateContexts gets the header for the current state of the session
func (s *Session) getStateContexts() (rsm.SessionCommandContext, []rsm.SessionStreamContext) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return rsm.SessionCommandContext{
		SessionID: s.SessionID,
		RequestID: s.responseID,
	}, s.getStreamContexts()
}

// getQueryContext gets the current read header
func (s *Session) getQueryContext(sync bool) rsm.SessionQueryContext {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return rsm.SessionQueryContext{
		SessionID:     s.SessionID,
		LastRequestID: s.responseID,
		LastIndex:     s.lastIndex,
		Sync:          sync,
	}
}

// nextCommandContext returns the next write context
func (s *Session) nextCommandContext() rsm.SessionCommandContext {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requestID = s.requestID + 1
	return rsm.SessionCommandContext{
		SessionID: s.SessionID,
		RequestID: s.requestID,
	}
}

// nextStreamHeader returns the next write stream and header
func (s *Session) nextStream() (*StreamState, rsm.SessionCommandContext) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requestID = s.requestID + 1
	stream := &StreamState{
		ID:      s.requestID,
		session: s,
	}
	s.streams[s.requestID] = stream
	command := rsm.SessionCommandContext{
		SessionID: s.SessionID,
		RequestID: s.requestID,
	}
	return stream, command
}

// recordCommandResponse records the index in a response header
func (s *Session) recordCommandResponse(requestContext rsm.SessionCommandContext, responseContext rsm.SessionResponseContext) {
	// Use a double-checked lock to avoid locking when multiple responses are received for an index.
	s.mu.RLock()
	if responseContext.Index > s.lastIndex {
		s.mu.RUnlock()
		s.mu.Lock()

		// If the request ID is greater than the highest response ID, update the response ID.
		if requestContext.RequestID > s.responseID {
			s.responseID = requestContext.RequestID
		}

		// If the response index has increased, update the last received index
		if responseContext.Index > s.lastIndex {
			s.lastIndex = responseContext.Index
		}
		s.mu.Unlock()
	} else {
		s.mu.RUnlock()
	}
}

// recordQueryResponse records the index in a response header
func (s *Session) recordQueryResponse(requestContext rsm.SessionQueryContext, responseContext rsm.SessionResponseContext) {
	// Use a double-checked lock to avoid locking when multiple responses are received for an index.
	s.mu.RLock()
	if responseContext.Index > s.lastIndex {
		s.mu.RUnlock()
		s.mu.Lock()

		// If the response index has increased, update the last received index
		if responseContext.Index > s.lastIndex {
			s.lastIndex = responseContext.Index
		}
		s.mu.Unlock()
	} else {
		s.mu.RUnlock()
	}
}

// deleteStream deletes the given stream from the session
func (s *Session) deleteStream(streamID uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.streams, streamID)
}

// getStreamHeaders returns a slice of headers for all open streams
func (s *Session) getStreamContexts() []rsm.SessionStreamContext {
	result := make([]rsm.SessionStreamContext, 0, len(s.streams))
	for _, stream := range s.streams {
		if stream.ID <= s.responseID {
			result = append(result, stream.getHeader())
		}
	}
	return result
}

// StreamState manages the context for a single response stream within a session
type StreamState struct {
	ID         uint64
	session    *Session
	responseID uint64
	mu         sync.RWMutex
}

// getHeader returns the current header for the stream
func (s *StreamState) getHeader() rsm.SessionStreamContext {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return rsm.SessionStreamContext{
		RequestID:     s.ID,
		AckResponseID: s.responseID,
	}
}

// serialize updates the stream response metadata and returns whether the response was received in sequential order
func (s *StreamState) serialize(context rsm.SessionResponseContext) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if context.Sequence == s.responseID+1 {
		s.responseID++
		return true
	}
	return false
}

// Close closes the stream
func (s *StreamState) Close() {
	s.session.deleteStream(s.ID)
}

// SessionOutput is a result for session-supporting servers containing session header information
type SessionOutput struct {
	streams.Result
}

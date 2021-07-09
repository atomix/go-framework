// Code generated by atomix-go-framework. DO NOT EDIT.
package lock

import (
	"context"
	lock "github.com/atomix/atomix-api/go/atomix/primitive/lock"
	"github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm"
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	storage "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	streams "github.com/atomix/atomix-go-framework/pkg/atomix/stream"
	"github.com/golang/protobuf/proto"
)

const Type = "Lock"

const (
	lockOp    storage.OperationID = 1
	unlockOp  storage.OperationID = 2
	getLockOp storage.OperationID = 3
)

var log = logging.GetLogger("atomix", "proxy", "lock")

// NewProxyServer creates a new ProxyServer
func NewProxyServer(client *rsm.Client, readSync bool) lock.LockServiceServer {
	return &ProxyServer{
		Client:   client,
		readSync: readSync,
	}
}

type ProxyServer struct {
	*rsm.Client
	readSync bool
	log      logging.Logger
}

func (s *ProxyServer) Lock(ctx context.Context, request *lock.LockRequest) (*lock.LockResponse, error) {
	log.Debugf("Received LockRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		log.Errorf("Request LockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	service := storage.ServiceID{
		Type:      Type,
		Namespace: s.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	session, err := partition.GetSession(ctx, service)
	if err != nil {
		return nil, err
	}

	ch := make(chan streams.Result)
	stream := streams.NewChannelStream(ch)
	err = session.DoCommandStream(ctx, lockOp, input, stream)
	if err != nil {
		log.Warnf("Request LockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	result, ok := <-ch
	if !ok {
		return nil, context.Canceled
	}

	if result.Failed() {
		log.Warnf("Request LockRequest failed: %v", result.Error)
		return nil, errors.Proto(result.Error)
	}

	response := &lock.LockResponse{}
	err = proto.Unmarshal(result.Value.([]byte), response)
	if err != nil {
		log.Errorf("Request LockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	log.Debugf("Sending LockResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) Unlock(ctx context.Context, request *lock.UnlockRequest) (*lock.UnlockResponse, error) {
	log.Debugf("Received UnlockRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		log.Errorf("Request UnlockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	service := storage.ServiceID{
		Type:      Type,
		Namespace: s.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	session, err := partition.GetSession(ctx, service)
	if err != nil {
		log.Errorf("Request UnlockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	output, err := session.DoCommand(ctx, unlockOp, input)
	if err != nil {
		log.Warnf("Request UnlockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &lock.UnlockResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		log.Errorf("Request UnlockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	log.Debugf("Sending UnlockResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) GetLock(ctx context.Context, request *lock.GetLockRequest) (*lock.GetLockResponse, error) {
	log.Debugf("Received GetLockRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		log.Errorf("Request GetLockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	service := storage.ServiceID{
		Type:      Type,
		Namespace: s.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	session, err := partition.GetSession(ctx, service)
	if err != nil {
		log.Errorf("Request GetLockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	output, err := session.DoQuery(ctx, getLockOp, input, s.readSync)
	if err != nil {
		log.Warnf("Request GetLockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &lock.GetLockResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		log.Errorf("Request GetLockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	log.Debugf("Sending GetLockResponse %+v", response)
	return response, nil
}

// Code generated by atomix-go-framework. DO NOT EDIT.
package leader

import (
	"context"
	leader "github.com/atomix/atomix-api/go/atomix/primitive/leader"
	"github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm"
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	storage "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	streams "github.com/atomix/atomix-go-framework/pkg/atomix/stream"
	"github.com/golang/protobuf/proto"
)

const Type = "LeaderLatch"

const (
	latchOp  = "Latch"
	getOp    = "Get"
	eventsOp = "Events"
)

// NewProxyServer creates a new ProxyServer
func NewProxyServer(client *rsm.Client, readSync bool) leader.LeaderLatchServiceServer {
	return &ProxyServer{
		Client:   client,
		readSync: readSync,
		log:      logging.GetLogger("atomix", "proxy", "leaderlatch"),
	}
}

type ProxyServer struct {
	*rsm.Client
	readSync bool
	log      logging.Logger
}

func (s *ProxyServer) Latch(ctx context.Context, request *leader.LatchRequest) (*leader.LatchResponse, error) {
	s.log.Debugf("Received LatchRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request LatchRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	service := storage.ServiceId{
		Type:    Type,
		Cluster: request.Headers.ClusterKey,
		Name:    request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoCommand(ctx, service, latchOp, input)
	if err != nil {
		s.log.Warnf("Request LatchRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &leader.LatchResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request LatchRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending LatchResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) Get(ctx context.Context, request *leader.GetRequest) (*leader.GetResponse, error) {
	s.log.Debugf("Received GetRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request GetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	service := storage.ServiceId{
		Type:    Type,
		Cluster: request.Headers.ClusterKey,
		Name:    request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoQuery(ctx, service, getOp, input, s.readSync)
	if err != nil {
		s.log.Warnf("Request GetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &leader.GetResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request GetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending GetResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) Events(request *leader.EventsRequest, srv leader.LeaderLatchService_EventsServer) error {
	s.log.Debugf("Received EventsRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request EventsRequest failed: %v", err)
		return errors.Proto(err)
	}

	ch := make(chan streams.Result)
	stream := streams.NewChannelStream(ch)
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	service := storage.ServiceId{
		Type:    Type,
		Cluster: request.Headers.ClusterKey,
		Name:    request.Headers.PrimitiveID.Name,
	}
	err = partition.DoCommandStream(srv.Context(), service, eventsOp, input, stream)
	if err != nil {
		s.log.Warnf("Request EventsRequest failed: %v", err)
		return errors.Proto(err)
	}

	for result := range ch {
		if result.Failed() {
			if result.Error == context.Canceled {
				break
			}
			s.log.Warnf("Request EventsRequest failed: %v", result.Error)
			return errors.Proto(result.Error)
		}

		response := &leader.EventsResponse{}
		err = proto.Unmarshal(result.Value.([]byte), response)
		if err != nil {
			s.log.Errorf("Request EventsRequest failed: %v", err)
			return errors.Proto(err)
		}

		s.log.Debugf("Sending EventsResponse %+v", response)
		if err = srv.Send(response); err != nil {
			s.log.Warnf("Response EventsResponse failed: %v", err)
			return err
		}
	}
	s.log.Debugf("Finished EventsRequest %+v", request)
	return nil
}

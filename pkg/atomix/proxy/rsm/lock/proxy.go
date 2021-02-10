package lock

import (
	"context"
	lock "github.com/atomix/api/go/atomix/primitive/lock"
	"github.com/atomix/go-framework/pkg/atomix/errors"
	"github.com/atomix/go-framework/pkg/atomix/logging"
	protocol "github.com/atomix/go-framework/pkg/atomix/protocol/rsm"
	"github.com/atomix/go-framework/pkg/atomix/proxy/rsm"
	"github.com/golang/protobuf/proto"
)

const Type = "Lock"

const (
	lockOp    = "Lock"
	unlockOp  = "Unlock"
	getLockOp = "GetLock"
)

// RegisterProxy registers the primitive on the given node
func RegisterProxy(node *rsm.Node) {
	node.PrimitiveTypes().RegisterProxyFunc(Type, func() (interface{}, error) {
		return &Proxy{
			Proxy: rsm.NewProxy(node.Client),
			log:   logging.GetLogger("atomix", "lock"),
		}, nil
	})
}

type Proxy struct {
	*rsm.Proxy
	log logging.Logger
}

func (s *Proxy) Lock(ctx context.Context, request *lock.LockRequest) (*lock.LockResponse, error) {
	s.log.Debugf("Received LockRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request LockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition, err := s.PartitionFrom(ctx)
	if err != nil {
		return nil, errors.Proto(err)
	}

	service := protocol.ServiceId{
		Type: Type,
		Name: request.Headers.PrimitiveID,
	}
	output, err := partition.DoCommand(ctx, service, lockOp, input)
	if err != nil {
		s.log.Errorf("Request LockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &lock.LockResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request LockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending LockResponse %+v", response)
	return response, nil
}

func (s *Proxy) Unlock(ctx context.Context, request *lock.UnlockRequest) (*lock.UnlockResponse, error) {
	s.log.Debugf("Received UnlockRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request UnlockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition, err := s.PartitionFrom(ctx)
	if err != nil {
		return nil, errors.Proto(err)
	}

	service := protocol.ServiceId{
		Type: Type,
		Name: request.Headers.PrimitiveID,
	}
	output, err := partition.DoCommand(ctx, service, unlockOp, input)
	if err != nil {
		s.log.Errorf("Request UnlockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &lock.UnlockResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request UnlockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending UnlockResponse %+v", response)
	return response, nil
}

func (s *Proxy) GetLock(ctx context.Context, request *lock.GetLockRequest) (*lock.GetLockResponse, error) {
	s.log.Debugf("Received GetLockRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request GetLockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition, err := s.PartitionFrom(ctx)
	if err != nil {
		return nil, errors.Proto(err)
	}

	service := protocol.ServiceId{
		Type: Type,
		Name: request.Headers.PrimitiveID,
	}
	output, err := partition.DoQuery(ctx, service, getLockOp, input)
	if err != nil {
		s.log.Errorf("Request GetLockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &lock.GetLockResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request GetLockRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending GetLockResponse %+v", response)
	return response, nil
}

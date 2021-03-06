// Code generated by atomix-go-framework. DO NOT EDIT.
package _map

import (
	_map "github.com/atomix/atomix-api/go/atomix/primitive/map"
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	"github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	"github.com/golang/protobuf/proto"
	"io"
)

var log = logging.GetLogger("atomix", "map", "service")

const Type = "Map"

const (
	sizeOp    = "Size"
	putOp     = "Put"
	getOp     = "Get"
	removeOp  = "Remove"
	clearOp   = "Clear"
	eventsOp  = "Events"
	entriesOp = "Entries"
)

var newServiceFunc rsm.NewServiceFunc

func registerServiceFunc(rsmf NewServiceFunc) {
	newServiceFunc = func(scheduler rsm.Scheduler, context rsm.ServiceContext) rsm.Service {
		service := &ServiceAdaptor{
			Service: rsm.NewService(scheduler, context),
			rsm:     rsmf(newServiceContext(scheduler)),
		}
		service.init()
		return service
	}
}

type NewServiceFunc func(ServiceContext) Service

// RegisterService registers the election primitive service on the given node
func RegisterService(node *rsm.Node) {
	node.RegisterService(Type, newServiceFunc)
}

type ServiceAdaptor struct {
	rsm.Service
	rsm Service
}

func (s *ServiceAdaptor) init() {
	s.RegisterUnaryOperation(sizeOp, s.size)
	s.RegisterUnaryOperation(putOp, s.put)
	s.RegisterUnaryOperation(getOp, s.get)
	s.RegisterUnaryOperation(removeOp, s.remove)
	s.RegisterUnaryOperation(clearOp, s.clear)
	s.RegisterStreamOperation(eventsOp, s.events)
	s.RegisterStreamOperation(entriesOp, s.entries)
}
func (s *ServiceAdaptor) SessionOpen(rsmSession rsm.Session) {
	s.rsm.Sessions().open(newSession(rsmSession))
}

func (s *ServiceAdaptor) SessionExpired(session rsm.Session) {
	s.rsm.Sessions().expire(SessionID(session.ID()))
}

func (s *ServiceAdaptor) SessionClosed(session rsm.Session) {
	s.rsm.Sessions().close(SessionID(session.ID()))
}
func (s *ServiceAdaptor) Backup(writer io.Writer) error {
	err := s.rsm.Backup(newSnapshotWriter(writer))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *ServiceAdaptor) Restore(reader io.Reader) error {
	err := s.rsm.Restore(newSnapshotReader(reader))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
func (s *ServiceAdaptor) size(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &_map.SizeRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		log.Warn(err)
		return nil, err
	}

	proposal := newSizeProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().Size().register(proposal)
	session.Proposals().Size().register(proposal)

	defer func() {
		session.Proposals().Size().unregister(proposal.ID())
		s.rsm.Proposals().Size().unregister(proposal.ID())
	}()

	log.Debugf("Proposing SizeProposal %s", proposal)
	err = s.rsm.Size(proposal)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) put(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &_map.PutRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		log.Warn(err)
		return nil, err
	}

	proposal := newPutProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().Put().register(proposal)
	session.Proposals().Put().register(proposal)

	defer func() {
		session.Proposals().Put().unregister(proposal.ID())
		s.rsm.Proposals().Put().unregister(proposal.ID())
	}()

	log.Debugf("Proposing PutProposal %s", proposal)
	err = s.rsm.Put(proposal)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) get(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &_map.GetRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		log.Warn(err)
		return nil, err
	}

	proposal := newGetProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().Get().register(proposal)
	session.Proposals().Get().register(proposal)

	defer func() {
		session.Proposals().Get().unregister(proposal.ID())
		s.rsm.Proposals().Get().unregister(proposal.ID())
	}()

	log.Debugf("Proposing GetProposal %s", proposal)
	err = s.rsm.Get(proposal)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) remove(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &_map.RemoveRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		log.Warn(err)
		return nil, err
	}

	proposal := newRemoveProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().Remove().register(proposal)
	session.Proposals().Remove().register(proposal)

	defer func() {
		session.Proposals().Remove().unregister(proposal.ID())
		s.rsm.Proposals().Remove().unregister(proposal.ID())
	}()

	log.Debugf("Proposing RemoveProposal %s", proposal)
	err = s.rsm.Remove(proposal)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) clear(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &_map.ClearRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		log.Warn(err)
		return nil, err
	}

	proposal := newClearProposal(ProposalID(s.Index()), session, request)

	s.rsm.Proposals().Clear().register(proposal)
	session.Proposals().Clear().register(proposal)

	defer func() {
		session.Proposals().Clear().unregister(proposal.ID())
		s.rsm.Proposals().Clear().unregister(proposal.ID())
	}()

	log.Debugf("Proposing ClearProposal %s", proposal)
	err = s.rsm.Clear(proposal)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	output, err := proto.Marshal(proposal.response())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) events(input []byte, rsmSession rsm.Session, stream rsm.Stream) (rsm.StreamCloser, error) {
	request := &_map.EventsRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		log.Warn(err)
		return nil, err
	}

	proposal := newEventsProposal(ProposalID(s.Index()), session, request, stream)

	s.rsm.Proposals().Events().register(proposal)
	session.Proposals().Events().register(proposal)

	log.Debugf("Proposing EventsProposal %s", proposal)
	err = s.rsm.Events(proposal)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	return func() {
		session.Proposals().Events().unregister(proposal.ID())
		s.rsm.Proposals().Events().unregister(proposal.ID())
	}, nil
}

func (s *ServiceAdaptor) entries(input []byte, rsmSession rsm.Session, stream rsm.Stream) (rsm.StreamCloser, error) {
	request := &_map.EntriesRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		log.Warn(err)
		return nil, err
	}

	proposal := newEntriesProposal(ProposalID(s.Index()), session, request, stream)

	s.rsm.Proposals().Entries().register(proposal)
	session.Proposals().Entries().register(proposal)

	log.Debugf("Proposing EntriesProposal %s", proposal)
	err = s.rsm.Entries(proposal)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	return func() {
		session.Proposals().Entries().unregister(proposal.ID())
		s.rsm.Proposals().Entries().unregister(proposal.ID())
	}, nil
}

var _ rsm.Service = &ServiceAdaptor{}

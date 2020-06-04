package grpcServer

import (
	"context"
	"github.com/stetsd/monk-scheduler/internal/api"
	"google.golang.org/grpc"
	"net"
)

type ApiServer struct {
	server      *grpc.Server
	GrpcEmitter *GrpcEmitter
}

func NewGrpcServer(grpcEmitter *GrpcEmitter) (*ApiServer, error) {
	// TODO: check localhost
	tcpConn, err := net.Listen("tcp", "0.0.0.0:50051")
	apiServer := ApiServer{
		GrpcEmitter: grpcEmitter,
	}

	if err != nil {
		return nil, err
	}

	apiServer.server = grpc.NewServer()
	api.RegisterApiServer(apiServer.server, apiServer)
	_ = apiServer.server.Serve(tcpConn)

	return &apiServer, nil
}

func (s ApiServer) SendEvent(ctx context.Context, event *api.Event) (*api.EventResult, error) {
	s.GrpcEmitter.emitOnEventMsg(event)

	return &api.EventResult{
		EventId: 1,
	}, nil
}

type GrpcEmitter struct {
	OnEventMsgCbQueue []func(event *api.Event)
}

func (s *GrpcEmitter) emitOnEventMsg(event *api.Event) {
	for _, cb := range s.OnEventMsgCbQueue {
		cb(event)
	}
}

func (s *GrpcEmitter) OnEventMsg(cb func(event *api.Event)) {
	s.OnEventMsgCbQueue = append(s.OnEventMsgCbQueue, cb)
}

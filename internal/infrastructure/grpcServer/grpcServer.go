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
	id, err := s.GrpcEmitter.OnEventMsgHandler(event)

	if err != nil {
		return &api.EventResult{
			Status:     1,
			StatusText: err.Error(),
		}, err
	}

	return &api.EventResult{
		EventId: int64(id),
	}, nil
}

type GrpcEmitter struct {
	OnEventMsgHandler func(event *api.Event) (int, error)
}

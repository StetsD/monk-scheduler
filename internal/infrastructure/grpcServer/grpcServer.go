package grpcServer

import (
	"context"
	"fmt"
	"github.com/stetsd/monk-scheduler/internal/api"
	"google.golang.org/grpc"
	"net"
)

type ApiServer struct {
	server *grpc.Server
}

func NewGrpcServer() (*ApiServer, error) {
	// TODO: check localhost
	tcpConn, err := net.Listen("tcp", "0.0.0.0:50051")
	apiServer := ApiServer{}

	if err != nil {
		return nil, err
	}

	apiServer.server = grpc.NewServer()
	api.RegisterApiServer(apiServer.server, apiServer)
	_ = apiServer.server.Serve(tcpConn)

	return &apiServer, nil
}

func (s ApiServer) SendEvent(ctx context.Context, event *api.Event) (*api.EventResult, error) {

	fmt.Println(event)

	return &api.EventResult{
		EventId: 1,
	}, nil
}

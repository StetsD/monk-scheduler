package app

import (
	"fmt"
	config "github.com/stetsd/monk-conf"
	"github.com/stetsd/monk-scheduler/internal/infrastructure/grpcServer"
)

type Scheduler struct {
	config config.Config
}

func NewApp(config config.Config) *Scheduler {
	return &Scheduler{config: config}
}

func (scheduler *Scheduler) Start() {
	fmt.Println("MELLO")

	grpcSos, err := grpcServer.NewGrpcServer()
	if err != nil {
		panic(err)
	}
	fmt.Println(grpcSos)
}

func (scheduler *Scheduler) Stop() {
	fmt.Println("STOP")
}

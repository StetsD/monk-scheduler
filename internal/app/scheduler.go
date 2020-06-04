package app

import (
	"fmt"
	config "github.com/stetsd/monk-conf"
	monk_db_driver "github.com/stetsd/monk-db-driver"
	"github.com/stetsd/monk-scheduler/internal/api"
	"github.com/stetsd/monk-scheduler/internal/app/contracts"
	"github.com/stetsd/monk-scheduler/internal/infrastructure/grpcServer"
	"github.com/stetsd/monk-scheduler/internal/infrastructure/logger"
	"time"
)

type Scheduler struct {
	config      config.Config
	apiServer   *grpcServer.ApiServer
	db          contracts.PgDriver
	eventPicker *EventPicker
}

func NewApp(config config.Config) *Scheduler {
	return &Scheduler{config: config}
}

func (scheduler *Scheduler) Start() {
	dbDriver, err := monk_db_driver.NewDbDriver(scheduler.config)
	if err != nil {
		panic(err)
	}
	scheduler.db = dbDriver

	exitChan := make(chan int)
	onSend := make(chan []onSendMsg)
	eventPicker := NewEventPicker(&exitChan, &onSend, dbDriver)
	scheduler.eventPicker = eventPicker
	scheduler.eventPicker.Start()

	grpcEmitter := grpcServer.GrpcEmitter{
		OnEventMsgHandler: func(event *api.Event) (int, error) {
			id, err := scheduler.CreateEvent(event)
			return id, err
		},
	}

	apiServer, err := grpcServer.NewGrpcServer(&grpcEmitter)

	if err != nil {
		panic(err)
	}

	scheduler.apiServer = apiServer
}

func (scheduler *Scheduler) CreateEvent(event *api.Event) (int, error) {
	rows, err := scheduler.db.Query(`
		INSERT INTO "Event" (title, dateStart, dateEnd, description, userId)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`,
		event.Title,
		time.Unix(event.DateStart.Seconds, 0).UTC().Format(time.RFC3339),
		time.Unix(event.DateEnd.Seconds, 0).UTC().Format(time.RFC3339),
		event.Description, event.UserId,
	)

	if err != nil {
		logger.Log.Error(err.Error())
	}

	defer func() {
		if err := rows.Close(); err != nil {
			logger.Log.Error(err.Error())
		}
	}()

	var id int

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			logger.Log.Error(err.Error())
		}
	}

	return id, nil
}

func (scheduler *Scheduler) Stop() {
	fmt.Println("STOP")
}

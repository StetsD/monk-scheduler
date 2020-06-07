package app

import (
	"fmt"
	"github.com/Shopify/sarama"
	config "github.com/stetsd/monk-conf"
	monk_db_driver "github.com/stetsd/monk-db-driver"
	"github.com/stetsd/monk-scheduler/internal/api"
	"github.com/stetsd/monk-scheduler/internal/app/contracts"
	"github.com/stetsd/monk-scheduler/internal/infrastructure"
	"github.com/stetsd/monk-scheduler/internal/infrastructure/grpcServer"
	"github.com/stetsd/monk-scheduler/internal/infrastructure/logger"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Scheduler struct {
	config          config.Config
	apiServer       *grpcServer.ApiServer
	db              contracts.PgDriver
	eventPicker     *EventPicker
	transportClient contracts.TransportClient
}

func NewApp(config config.Config) *Scheduler {
	return &Scheduler{config: config}
}

func (scheduler *Scheduler) Start() {
	logger.Log.Info("service monk-scheduler is running")
	dbDriver, err := monk_db_driver.NewDbDriver(scheduler.config)
	if err != nil {
		panic(err)
	}
	scheduler.db = dbDriver

	onSend := make(chan onSendMsg)

	producer, err := scheduler.ConnectToTransportAsProducer()

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			logger.Log.Fatal(err.Error())
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	eventPicker := NewEventPicker(&signals, &onSend, dbDriver)
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

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			logger.Log.Error(err.Error())
		}
	}()

outer:
	for {
		select {
		case event := <-onSend:
			producer.Input() <- &sarama.ProducerMessage{Topic: "on_send", Value: sarama.StringEncoder(onSendMarshaling(&event))}
		case <-signals:
			apiServer.Stop()
			producer.AsyncClose()
			break outer
		}
	}

	wg.Wait()

	scheduler.apiServer = apiServer
}

func (scheduler *Scheduler) CreateEvent(event *api.Event) (int, error) {
	rows, err := scheduler.db.Query(`
		INSERT INTO "Event" (title, dateStart, dateEnd, description, userId, email)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
	`,
		event.Title,
		time.Unix(event.DateStart.Seconds, 0).UTC().Format(time.RFC3339),
		time.Unix(event.DateEnd.Seconds, 0).UTC().Format(time.RFC3339),
		event.Description, event.UserId, event.Email,
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

func (scheduler *Scheduler) ConnectToTransportAsProducer() (sarama.AsyncProducer, error) {
	scheduler.transportClient = infrastructure.NewKafkaClient(scheduler.config)
	producer, err := scheduler.transportClient.InitProducer()
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func (scheduler *Scheduler) Stop() {
	// TODO: implement norm stop mech
	fmt.Println("STOP")
}

// TODO:makefile create

package app

import (
	"fmt"
	"github.com/stetsd/monk-scheduler/internal/app/constants"
	"github.com/stetsd/monk-scheduler/internal/app/contracts"
	"github.com/stetsd/monk-scheduler/internal/infrastructure/logger"
	"time"
)

type EventPicker struct {
	onSend   *chan []onSendMsg
	exitChan *chan int
	ticker   *time.Ticker
	db       contracts.PgDriver
}

func NewEventPicker(exitChan *chan int, onSend *chan []onSendMsg, db contracts.PgDriver) *EventPicker {
	return &EventPicker{
		exitChan: exitChan,
		onSend:   onSend,
		db:       db,
	}
}

func (ep *EventPicker) Pick() {
	now := time.Now().Round(0).Format(constants.TimeFormat)

	rows, err := ep.db.Query(`
		SELECT id, title, description FROM "Event" WHERE datestart = $1;
	`, now)

	if err != nil {
		logger.Log.Error(err.Error())
	}

	defer func() {
		if err := rows.Close(); err != nil {
			logger.Log.Error(err.Error())
		}
	}()

	var sendColl []onSendMsg

	for rows.Next() {
		line := onSendMsg{}
		if err := rows.Scan(&line.Id, &line.Title, &line.Description); err != nil {
			logger.Log.Error(err.Error())
			return
		}
		sendColl = append(sendColl, line)
	}

	fmt.Println(sendColl)
}

func (ep *EventPicker) Start() {
	ep.ticker = time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ep.ticker.C:
				ep.Pick()
			case val := <-*ep.exitChan:
				if val == 1 {
					ep.Stop()
					return
				}
			}
		}
	}()
}

func (ep *EventPicker) Stop() {
	ep.ticker.Stop()
}

package app

import (
	"github.com/stetsd/monk-scheduler/internal/app/constants"
	"github.com/stetsd/monk-scheduler/internal/app/contracts"
	"github.com/stetsd/monk-scheduler/internal/infrastructure/logger"
	"os"
	"time"
)

type EventPicker struct {
	onSend   *chan onSendMsg
	exitChan *chan os.Signal
	ticker   *time.Ticker
	db       contracts.PgDriver
}

func NewEventPicker(exitChan *chan os.Signal, onSend *chan onSendMsg, db contracts.PgDriver) *EventPicker {
	return &EventPicker{
		exitChan: exitChan,
		onSend:   onSend,
		db:       db,
	}
}

func (ep *EventPicker) Pick() {
	now := time.Now().Round(0).Format(constants.TimeFormat)

	rows, err := ep.db.Query(`
		SELECT id, title, description, email FROM "Event" WHERE datestart = $1;
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
		if err := rows.Scan(&line.Id, &line.Title, &line.Description, &line.Email); err != nil {
			logger.Log.Error(err.Error())
			return
		}
		sendColl = append(sendColl, line)
	}

	if len(sendColl) != 0 {
		for _, event := range sendColl {
			*ep.onSend <- onSendMsg{
				Id:          event.Id,
				Description: event.Description,
				Title:       event.Title,
				Email:       event.Email,
			}
		}
	}
}

func (ep *EventPicker) Start() {
	ep.ticker = time.NewTicker(1 * time.Second)

	go func() {
	outer:
		for {
			select {
			case <-ep.ticker.C:
				ep.Pick()
			case <-*ep.exitChan:
				ep.Stop()
				break outer
			}
		}
	}()
}

func (ep *EventPicker) Stop() {
	ep.ticker.Stop()
}

package eventlog

import (
	cores "github.com/catalyst-zero/cores-go"

	"log"
	"os"
)

func NewCoresAmqpEventLog(amqpUrl string) *coresEventLog {
	bus, err := cores.NewAmqpEventBus(amqpUrl)
	if err != nil {
		panic(err)
	}

	return &coresEventLog{
		Bus:    bus,
		Logger: log.New(os.Stderr, "[cores] ", log.LstdFlags),
	}
}

type coresEventLog struct {
	Bus    cores.EventBus
	Logger *log.Logger
}

func (log *coresEventLog) NewGame(gameId string) {
	go log.Bus.Publish("tictactoed.started", gameId)
}

func (log *coresEventLog) Moved(gameId string) {
	go log.Bus.Publish("tictactoed.moved", gameId)
}

func (log *coresEventLog) Finished(gameId string) {
	go log.Bus.Publish("tictactoed.finished", gameId)
}

func (log *coresEventLog) publish(eventName, gameId string) {
	if err := log.Bus.Publish(eventName, gameId); err != nil {
		log.Logger.Fatalf("Failed to publish event to cores: %v", err)
	}
}

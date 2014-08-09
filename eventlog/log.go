package eventlog

import (
	"log"
	"os"
)

func NewLog() *logBus {
	return &logBus{
		logger: log.New(os.Stdout, "[event] ", log.LstdFlags),
	}
}

type logBus struct {
	logger *log.Logger
}

func (log *logBus) NewGame(gameId string) {
	log.logger.Printf("[%s] New Game started.", gameId)
}
func (log *logBus) Moved(gameId string) {
	log.logger.Printf("[%s] Move", gameId)
}
func (log *logBus) Finished(gameId string) {
	log.logger.Printf("[%s] Finished", gameId)
}

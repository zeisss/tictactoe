package service

import (
	"../game"

	"time"
)

type Storage interface {
	Get(gameId string) (game.GameState, error)
	Put(gameId string, game game.GameState) error
}

// EventLog notifies any interested messaging-bus about the events that occur.
// This is a write-only interface.
type EventLog interface {
	NewGame(gameId string)
	Moved(gameId string)
	Finished(gameId string)
}

// This is a write-only interface.
type LockManager interface {
	// Grab a lock
	LockGame(gameId string, timeout time.Duration)

	// Unlock
	UnlockGame(gameId string)
}

// IdFactory is used to generate IDs for new games.
type IdFactory func() (string, error)

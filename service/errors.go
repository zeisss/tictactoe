package service

// Wraps all 'good/known' backend errors in IsXXX() functions so clients don't need to know the backend errors.

import (
	"../game"
	"../storage"
)

func IsNotFoundError(err error) bool {
	return err == storage.NotFoundError
}

func IsNotActivePlayerError(err error) bool {
	return err == game.PlayerNotActive
}

func IsGameFinishedError(err error) bool {
	return err == game.GameIsFinished
}

func IsPositionAlreadyTakenError(err error) bool {
	return err == game.PositionAlreadyTaken
}

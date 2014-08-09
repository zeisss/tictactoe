package service

import (
	"../game"

	"time"
)

// TicTacToeService provides the service-composition logic which includes the game-logic with its backend-services.
type TicTacToeService struct {
	IdFactory   IdFactory
	Storage     Storage
	EventLog    EventLog
	LockManager LockManager
}

// New creates a new game and returns its id.
func (s *TicTacToeService) New() (string, error) {
	gameId, err := s.IdFactory()
	if err != nil {
		return gameId, err
	}

	state := game.New()

	if err := s.Storage.Put(gameId, state); err != nil {
		return gameId, err
	}

	s.EventLog.NewGame(gameId)
	return gameId, nil
}

// Move performs the move on the given gameId and stores it afterwards.
//
// Error Helpers:
//  IsNotActivePlayerError()
//  IsGameFinishedError()
//  IsPositionAlreadyTakenError()
//  IsNotFoundErrorError()
func (s *TicTacToeService) Move(gameId string, player game.Player, position game.Position) error {
	s.LockManager.LockGame(gameId, 10*time.Minute)
	defer s.LockManager.UnlockGame(gameId)

	state, err := s.Storage.Get(gameId)
	if err != nil {
		return err
	}

	if err := game.Move(&state, player, position); err != nil {
		return err
	}

	if err := s.Storage.Put(gameId, state); err != nil {
		return err
	}

	s.EventLog.Moved(gameId)
	if state.Status == game.STATUS_FINISHED {
		s.EventLog.Finished(gameId)
	}

	return nil
}

// Returns the current game state.
//
// Error Helpers:
//  IsNotFoundErrorError()
func (s *TicTacToeService) Get(gameId string) (game.GameState, error) {
	return s.Storage.Get(gameId)
}

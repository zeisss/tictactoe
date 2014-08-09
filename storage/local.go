package storage

import (
	"../game"
)

func New() *localStorage {
	return &localStorage{
		Games: make(map[string]game.GameState),
	}
}

type localStorage struct {
	Games map[string]game.GameState
}

func (s *localStorage) Get(gameId string) (game.GameState, error) {
	game, ok := s.Games[gameId]
	if !ok {
		return game, NotFoundError
	}
	return game, nil
}

func (s *localStorage) Put(gameId string, game game.GameState) error {
	s.Games[gameId] = game
	return nil
}

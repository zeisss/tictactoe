package storage

import (
	"errors"
)

var (
	NotFoundError = errors.New("No game found with the given game-id.")
)

package idfactory

import (
	"code.google.com/p/go-uuid/uuid"
)

func UUIDNextId() (string, error) {
	return uuid.New(), nil
}

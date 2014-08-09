package idfactory

import (
	"fmt"
)

func NewSequenceFactory(format string) func() (string, error) {
	if format == "" {
		format = "game%d"
	}

	var idSequence uint = 0
	return func() (string, error) {
		idSequence++
		return fmt.Sprintf(format, idSequence), nil
	}
}

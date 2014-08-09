package game

import (
	"encoding/json"
)

func (gs *GameState) TextMarshaler() ([]byte, error) {
	return json.Marshal(gs)
}

package game

import (
	"errors"
)

var (
	GameIsFinished       = errors.New("Game is finished.")
	PositionAlreadyTaken = errors.New("Invalid move - position already taken.")
	PlayerNotActive      = errors.New("Invalid player - only the active player can move.")
)

/// ------------------------------

const (
	NO_PLAYER Player = 0
	PLAYER_1  Player = 1
	PLAYER_2  Player = 2
)

type Player int

func (p Player) Valid() bool {
	return (p >= 0 && p <= PLAYER_2)
}

func (p Player) Validate() {
	if !p.Valid() {
		panic("Invalid player value.")
	}
}

/// ------------------------------

type Position struct {
	X, Y int
}

func (pos Position) Valid() bool {
	return !(pos.X < 0 || pos.Y < 0 || pos.X > 2 || pos.Y > 2)
}

func (pos Position) Validate() {
	if !pos.Valid() {
		panic("Invalid coordinates.")
	}
}

/// ------------------------------

type Status string

const (
	STATUS_NEW        Status = "new"
	STATUS_INPROGRESS Status = "inprogress"
	STATUS_FINISHED   Status = "finished"
)

/// ------------------------------

func New() GameState {
	return GameState{
		Turn:         0,
		ActivePlayer: PLAYER_1,
		Status:       STATUS_NEW,
		Fields: [][]Player{
			[]Player{NO_PLAYER, NO_PLAYER, NO_PLAYER},
			[]Player{NO_PLAYER, NO_PLAYER, NO_PLAYER},
			[]Player{NO_PLAYER, NO_PLAYER, NO_PLAYER},
		},
	}
}

type GameState struct {
	Turn         uint8      `json:"turn"`
	ActivePlayer Player     `json:"active_player"`
	Status       Status     `json:"status"`
	Fields       [][]Player `json:"fields"`
}

func (gs *GameState) Get(x, y int) Player {
	return gs.Fields[y][x]
}
func (gs *GameState) Set(x, y int, player Player) {
	gs.Fields[y][x] = player
}

// Move performs a move onto the given position by the given player.
// If this move is invalid, an error is returned.
func Move(game *GameState, player Player, pos Position) error {
	// fmt.Printf("Move(%v, %v, %v)\n", game, player, pos)

	// Validation
	player.Validate()
	pos.Validate()

	if game.Status == STATUS_FINISHED {
		return GameIsFinished
	}

	if game.ActivePlayer != player {
		return PlayerNotActive
	}

	// Action

	owner := game.Get(pos.X, pos.Y)
	if owner != NO_PLAYER {
		return PositionAlreadyTaken
	}

	game.Set(pos.X, pos.Y, player)

	finished, _ := IsFinished(game)
	if finished {
		game.Status = STATUS_FINISHED
		game.ActivePlayer = NO_PLAYER
	} else {
		game.Status = STATUS_INPROGRESS
		if game.ActivePlayer == PLAYER_1 {
			game.ActivePlayer = PLAYER_2
		} else {
			game.ActivePlayer = PLAYER_1
		}
	}
	game.Turn++

	return nil
}

// IsFinished returns the player that won the game or (false, NO_PLAYER) if no player won yet.
func IsFinished(game *GameState) (bool, Player) {
	var p1, p2, p3 Player

	// Vertical checks (3x)
	for x := 0; x < 3; x++ {
		p1 = game.Get(x, 0)
		p2 = game.Get(x, 1)
		p3 = game.Get(x, 2)

		if p1 != NO_PLAYER && p1 == p2 && p2 == p3 {
			return true, p1
		}
	}

	// Horizontal checks (3x)
	for y := 0; y < 3; y++ {
		p1 = game.Get(0, y)
		p2 = game.Get(1, y)
		p3 = game.Get(2, y)

		if p1 != NO_PLAYER && p1 == p2 && p2 == p3 {
			return true, p1
		}
	}

	// Diagonal checks (2x)
	p1 = game.Get(0, 0)
	p2 = game.Get(1, 1)
	p3 = game.Get(2, 2)

	if p1 != NO_PLAYER && p1 == p2 && p2 == p3 {
		return true, p1
	}

	p1 = game.Get(2, 0)
	p2 = game.Get(1, 1)
	p3 = game.Get(0, 2)

	if p1 != NO_PLAYER && p1 == p2 && p2 == p3 {
		return true, p1
	}

	return false, NO_PLAYER
}

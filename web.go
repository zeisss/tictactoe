package main

import (
	"./game"
	"./service"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type BaseHandler struct {
	Service *service.TicTacToeService
}

func (h *BaseHandler) setNotFoundError(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusNotFound)
}

func (h *BaseHandler) setBadRequestError(resp http.ResponseWriter, msg string) {
	resp.WriteHeader(http.StatusBadRequest)
	resp.Write([]byte(msg))
}

func (h *BaseHandler) setProcessingError(resp http.ResponseWriter, err error) {
	resp.WriteHeader(http.StatusInternalServerError)

	log.Printf("Internal error: %#v\n", err)
}

func (h *BaseHandler) GameID(resp http.ResponseWriter, req *http.Request) (string, bool) {
	gameID := req.FormValue("game")
	if gameID == "" {
		h.setBadRequestError(resp, "No 'game' parameter given.")
		return "", false
	}
	return gameID, true
}

// ------------------------

type NewGameHandler struct {
	BaseHandler
}

func (h *NewGameHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	gameID, err := h.Service.New()
	if err != nil {
		h.setProcessingError(resp, err)
	} else {
		resp.Header().Add("location", "/game/get?game="+gameID)
		resp.WriteHeader(http.StatusCreated)
		resp.Write([]byte(gameID))
	}
}

// ------------------------

type MoveHandler struct {
	BaseHandler
}

func (h *MoveHandler) Player(resp http.ResponseWriter, req *http.Request) (game.Player, bool) {
	playerParam := req.FormValue("player")
	var player game.Player = game.NO_PLAYER
	if playerParam == "player1" {
		player = game.PLAYER_1
	} else if playerParam == "player2" {
		player = game.PLAYER_2
	} else {
		h.setBadRequestError(resp, "Invalid value for parameter 'player'.")
		return game.NO_PLAYER, false
	}
	return player, true
}

func (h *MoveHandler) Position(resp http.ResponseWriter, req *http.Request) (game.Position, bool) {
	var result game.Position

	positionParam := req.FormValue("position")
	positionSplit := strings.Split(positionParam, ",")
	if len(positionSplit) != 2 {
		h.setBadRequestError(resp, "Invalid value for parameter 'position'.")
		return result, false
	}

	x, err := strconv.ParseInt(positionSplit[0], 10, 0)
	if err != nil {
		h.setBadRequestError(resp, "Invalid value for parameter 'position'.")
		return result, false
	}
	y, err := strconv.ParseInt(positionSplit[1], 10, 0)
	if err != nil {
		h.setBadRequestError(resp, "Invalid value for parameter 'position'.")
		return result, false
	}

	result = game.Position{int(x), int(y)}
	if !result.Valid() {
		h.setBadRequestError(resp, "Invalid value for parameter 'position'.")
		return result, false
	}
	return result, true
}

func (h *MoveHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get GameID (can be taken as is)
	gameID, ok := h.GameID(resp, req)
	if !ok {
		return
	}

	// Get player
	player, ok := h.Player(resp, req)
	if !ok {
		return
	}

	// Get position parameter (..&position=2,2)
	position, ok := h.Position(resp, req)
	if !ok {
		return
	}

	// Perform move
	if err := h.Service.Move(gameID, player, position); err != nil {
		if service.IsNotFoundError(err) {
			h.setNotFoundError(resp)
		} else if service.IsGameFinishedError(err) {
			h.setBadRequestError(resp, "Game is finished.")
		} else if service.IsPositionAlreadyTakenError(err) {
			h.setBadRequestError(resp, "Position is already taken. Invalid move.")
		} else if service.IsNotActivePlayerError(err) {
			h.setBadRequestError(resp, fmt.Sprintf("Invalid move. '%s' is not the active player.", player))
		} else {
			h.setProcessingError(resp, err)
		}
	} else {
		resp.WriteHeader(http.StatusNoContent)
	}
}

// ------------------------

type GetGameHandler struct {
	BaseHandler
}

func (h *GetGameHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get GameID (can be taken as is)
	gameID, ok := h.GameID(resp, req)
	if !ok {
		return
	}

	state, err := h.Service.Get(gameID)
	if err != nil {
		if service.IsNotFoundError(err) {
			h.setNotFoundError(resp)
		} else {
			h.setProcessingError(resp, err)
		}
		return
	}

	result := map[string]interface{}{}
	result["turn"] = state.Turn
	if state.Status != game.STATUS_FINISHED {
		if state.ActivePlayer == game.PLAYER_1 {
			result["active_player"] = "player1"
		} else {
			result["active_player"] = "player2"
		}
	}
	result["fields"] = state.Fields
	result["status"] = state.Status

	if err := json.NewEncoder(resp).Encode(result); err != nil {
		panic(err)
	}
}

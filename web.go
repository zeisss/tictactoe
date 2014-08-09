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

type Handlers struct {
	Service *service.TicTacToeService
}

func (h *Handlers) NewGameHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	gameId, err := h.Service.New()
	if err != nil {
		h.setProcessingError(resp, err)
	} else {
		resp.Header().Add("location", "/game/get?game="+gameId)
		resp.WriteHeader(http.StatusCreated)
		resp.Write([]byte(gameId))
	}
}

func (h *Handlers) MoveHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get GameId (can be taken as is)
	gameId := req.FormValue("game")
	if gameId == "" {
		h.setBadRequestError(resp, "No game parameter given.")
		return
	}

	// Get player
	playerParam := req.FormValue("player")
	var player game.Player = game.NO_PLAYER
	if playerParam == "player1" {
		player = game.PLAYER_1
	} else if playerParam == "player2" {
		player = game.PLAYER_2
	} else {
		h.setBadRequestError(resp, "Invalid value for parameter 'player'.")
		return
	}

	// Get position parameter (..&position=2,2)
	positionParam := req.FormValue("position")
	positionSplit := strings.Split(positionParam, ",")
	if len(positionSplit) != 2 {
		h.setBadRequestError(resp, "Invalid value for parameter 'position'.")
		return
	}

	x, err := strconv.ParseInt(positionSplit[0], 10, 0)
	if err != nil {
		h.setBadRequestError(resp, "Invalid value for parameter 'position'.")
		return
	}
	y, err := strconv.ParseInt(positionSplit[1], 10, 0)
	if err != nil {
		h.setBadRequestError(resp, "Invalid value for parameter 'position'.")
		return
	}

	position := game.Position{int(x), int(y)}
	if !position.Valid() {
		h.setBadRequestError(resp, "Invalid value for parameter 'position'.")
	}

	// Perform move
	if err := h.Service.Move(gameId, player, position); err != nil {
		if service.IsNotFoundError(err) {
			h.setNotFoundError(resp)
		} else if service.IsGameFinishedError(err) {
			h.setBadRequestError(resp, "Game is finished.")
		} else if service.IsPositionAlreadyTakenError(err) {
			h.setBadRequestError(resp, "Position is already taken. Invalid move.")
		} else if service.IsNotActivePlayerError(err) {
			h.setBadRequestError(resp, fmt.Sprintf("Invalid move. '%s' is not the active player.", playerParam))
		} else {
			h.setProcessingError(resp, err)
		}
	} else {
		resp.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handlers) GetGameHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get GameId (can be taken as is)
	gameId := req.FormValue("game")
	if gameId == "" {
		h.setBadRequestError(resp, "Invalid value for parameter 'game'.")
		return
	}

	state, err := h.Service.Get(gameId)
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

func (h *Handlers) setNotFoundError(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusNotFound)
}

func (h *Handlers) setBadRequestError(resp http.ResponseWriter, msg string) {
	resp.WriteHeader(http.StatusBadRequest)
	resp.Write([]byte(msg))
}

func (h *Handlers) setProcessingError(resp http.ResponseWriter, err error) {
	resp.WriteHeader(http.StatusInternalServerError)

	log.Printf("Internal error: %#v\n", err)
}

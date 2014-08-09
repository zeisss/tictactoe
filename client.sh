#!/bin/bash

########
#
# Yeah, horrible. I know :D Still works though *runs*
#
#
########

LAST_REQUEST=
LAST_OUTPUT=

function request() {
	local curl_call="curl -s -X ${METHOD:=GET} $*"
	#echo "> ${curl_call}"
	LAST_REQUEST="${curl_call}"
	LAST_OUTPUT=$(${curl_call})

	#echo
	#echo "$ ${curl_call}"
	#echo "> ${LAST_OUTPUT}"
	#echo
}

function assert_request_output() {
	local expected_output=$1
	assert_equals "$expected_output" "${LAST_OUTPUT}"
}

function assert_equals() {
	local expected=$1
	local actual=$2

	if [ "${expected}" != "${actual}" ]; then
		echo 
		echo "Expected: ${expected_output}"
		echo "Actual:   ${actual}"
		echo
		echo "==> ✕"
		exit 1
	else
		echo -n "✓"
	fi
}

function givenNewGame() {
	# Start new game
	METHOD="POST" request "localhost:8080/game/new"
	GAME=$LAST_OUTPUT

	export GAME
}


function test_new_game() {
	givenNewGame

	if [ -z "${GAME}" ]; then
		echo "No game id received for new game!"
		exit 1
	fi

	# Assert game is in 'new'
	request "localhost:8080/game/get?game=${GAME}"
	assert_request_output \
			'{"active_player":"player1","fields":[[0,0,0],[0,0,0],[0,0,0]],"status":"new","turn":0}'
}

function test_players_move() {
	# Setup new game
	givenNewGame

	# Perform some moves
	METHOD="POST" request "localhost:8080/game/move?game=${GAME}&player=player1&position=0,0"

	request "localhost:8080/game/get?game=${GAME}"
	assert_request_output \
			'{"active_player":"player2","fields":[[1,0,0],[0,0,0],[0,0,0]],"status":"inprogress","turn":1}'

	METHOD="POST" request "localhost:8080/game/move?game=${GAME}&player=player2&position=1,0"

	METHOD="POST" request "localhost:8080/game/move?game=${GAME}&player=player1&position=0,1"

	METHOD="POST" request "localhost:8080/game/move?game=${GAME}&player=player2&position=1,1"

	request "localhost:8080/game/get?game=${GAME}"
	assert_request_output '{"active_player":"player1","fields":[[1,2,0],[1,2,0],[0,0,0]],"status":"inprogress","turn":4}'

	METHOD="POST" request "localhost:8080/game/move?game=${GAME}&player=player1&position=0,2"

	# Assert 'finished' state
	request "localhost:8080/game/get?game=${GAME}"
	assert_request_output '{"fields":[[1,2,0],[1,2,0],[1,0,0]],"status":"finished","turn":5}'			


	# Assert no more move available
	METHOD="POST" request "localhost:8080/game/move?game=${GAME}&player=player2&position=1,2"
	assert_request_output 'Game is finished.'	
}

function test_coordinate_already_taken() {
	givenNewGame

	METHOD="POST" request "localhost:8080/game/move?game=${GAME}&player=player1&position=0,0"
	METHOD="POST" request "localhost:8080/game/move?game=${GAME}&player=player2&position=1,0"

	METHOD="POST" request "localhost:8080/game/move?game=${GAME}&player=player1&position=1,0"
	assert_request_output 'Position is already taken. Invalid move.'
}

function test_same_player_cannot_move_twice() {
	givenNewGame

	METHOD="POST" request "localhost:8080/game/move?game=${GAME}&player=player1&position=0,0"
	METHOD="POST" request "localhost:8080/game/move?game=${GAME}&player=player1&position=0,1"
	assert_request_output "Invalid move. 'player1' is not the active player."

	request "localhost:8080/game/get?game=${GAME}"
	assert_request_output '{"active_player":"player2","fields":[[1,0,0],[0,0,0],[0,0,0]],"status":"inprogress","turn":1}'
}

function main() {
	test_new_game
	test_players_move
	test_coordinate_already_taken
	test_same_player_cannot_move_twice
}

set -e
main
echo

package main

import (
	"net/http"
	"strconv"
)

var game [3][3]int // game bord
var player int     // current player turn
var turn int       // the turn
var res string     // the winner
var scoreJ1 int    // score of the Player 1
var scoreJ2 int    // score of the Player 2

func play(writer http.ResponseWriter, request *http.Request) {
	if turn == 0 {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				game[i][j] = -1
			}
		}
	}
	x, _ := strconv.Atoi(request.FormValue("x"))
	y, _ := strconv.Atoi(request.FormValue("y"))
	game[x][y] = player
	str := "X"
	turn++
	if game[x][y] == 1 {
		str = "O"
	}
	for i := 0; i < 3; i++ {
		if game[i][0] == game[i][1] && game[i][0] == game[i][2] && game[i][0] != -1 {
			res = strconv.Itoa(player)
		}
		if game[0][i] == game[1][i] && game[0][i] == game[2][i] && game[0][i] != -1 {
			res = strconv.Itoa(player)
		}
	}
	if (game[0][0] == game[1][1] && game[0][0] == game[2][2] && game[0][0] != -1) || (game[0][2] == game[1][1] && game[0][2] == game[2][0] && game[0][2] != -1) {
		res = strconv.Itoa(player)
	} else if turn == 9 {
		res = "Draw"
	}
	str += res
	player = player%2 + 1
	writer.Header().Set("Cache-Control", "no-cache")
	_, _ = writer.Write([]byte(str))
}

func getVictory(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Cache-Control", "no-cache")
	_, _ = writer.Write([]byte(res))
}

func getScore(writer http.ResponseWriter, _ *http.Request) {
	score := "<p>Player 1 : " + strconv.Itoa(scoreJ1) + "</p><p>Player 2 : " + strconv.Itoa(scoreJ2) + "</p>"
	writer.Header().Set("Cache-Control", "no-cache")
	_, _ = writer.Write([]byte(score))
}

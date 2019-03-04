package main

import (
	"github.com/zserge/webview"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"net"
)

var dir string                           // current directory
var windowWidth, windowHeight = 400, 300 // width and height of the window
var game [3][3]int						 // game bord
var player = 0							 // current player turn
var turn = 0
var res = ""							 // the winner

func init() {
	// getting the current directory to access resources
	var err error
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
}

// main function
func main() {
	// channel to get the web prefix
	prefixChannel := make(chan string)
	// run the web server in a separate goroutine
	go app(prefixChannel)
	prefix := <-prefixChannel
	// create a web view
	err := webview.Open("TicTacGo", prefix+"/public/html/index.html",
		windowWidth, windowHeight, false)
	if err != nil {
		log.Fatal(err)
	}
}

// web app
func app(prefixChannel chan string) {
	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(dir+"/public"))))
	mux.HandleFunc("/start", start)
	mux.HandleFunc("/play", play)
	mux.HandleFunc("/victory", getVictory)
	mux.HandleFunc("/res", victory)

	// get an ephemeral port, so we're guaranteed not to conflict with anything else
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	portAddress := listener.Addr().String()
	prefixChannel <- "http://" + portAddress
	listener.Close()
	server := &http.Server{
		Addr:    portAddress,
		Handler: mux,
	}
	server.ListenAndServe()
}

func play(writer http.ResponseWriter, request *http.Request) {
	x, _ := strconv.Atoi(request.FormValue("x"))
	y, _ := strconv.Atoi(request.FormValue("y"))
	game[x][y] = player
	str := "X"
	turn++
	if game[x][y] == 0 {
		str = "O"
	}
	winner := player+1
	for i := 0; i < 3; i++ {
		if game[i][0] == game[i][1] && game[i][0] == game[i][2] && game[i][0] != -1 {
			res = "Winner : Player " + strconv.Itoa(winner)
		}
		if game[0][i] == game[1][i] && game[0][i] == game[2][i] && game[0][i] != -1 {
			res = "Winner : Player " + strconv.Itoa(winner)
		}
	}
	if game[0][0] == game[1][1] && game[0][0] == game[2][2] && game[0][0] != -1 {
		res = "Winner : Player " + strconv.Itoa(winner)
	} else if game[0][2] == game[1][1] && game[0][2] == game[2][0] && game[0][2] != -1 {
		res = "Winner : Player " + strconv.Itoa(winner)
	} else if turn == 9 {
		res = "Draw"
	}
	str += res
	player = (player + 1) % 2
	writer.Header().Set("Cache-Control", "no-cache")
	_, _ = writer.Write([]byte(str))
}

func getVictory(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Cache-Control", "no-cache")
	_, _ = writer.Write([]byte(res))
}

func victory(writer http.ResponseWriter, _ *http.Request) {
	t, _ := template.ParseFiles(dir + "/public/html/victory.html")
	_ = t.Execute(writer, nil)
}

// start the game
func start(w http.ResponseWriter, _ *http.Request) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			game[i][j] = -1
		}
	}
	res = ""
	turn = 0
	player = 0
	t, _ := template.ParseFiles(dir + "/public/html/tictactoe.html")
	// start generating frames in a new goroutine
	_ = t.Execute(w, nil)
}

package main

import (
	"github.com/zserge/webview"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var dir string                           // current directory
var windowWidth, windowHeight = 400, 300 // width and height of the window

func init() {
	// getting the current directory to access resources
	var err error
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			game[i][j] = -1
		}
	}
	res = ""
	turn = 0
	player = 1
	scoreJ1 = 0
	scoreJ2 = 0
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
	mux.HandleFunc("/restart", restart)
	mux.HandleFunc("/play", play)
	mux.HandleFunc("/victory", getVictory)
	mux.HandleFunc("/score", getScore)

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

// start the game
func start(w http.ResponseWriter, r *http.Request) {
	restart(w, r)
	t, _ := template.ParseFiles(dir + "/public/html/tictactoe.html")
	_ = t.Execute(w, nil)
}

// start the game
func restart(w http.ResponseWriter, _ *http.Request) {
	if res != "" {
		if res == "1" {
			scoreJ1 += 1
		} else if res == "2" {
			scoreJ2 += 1
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				game[i][j] = -1
			}
		}
		res = ""
		turn = 0
		player = 1
		w.Header().Set("Cache-Control", "no-cache")
	}
}


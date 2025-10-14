package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Results struct {
	OpenGames    int
	JaneWins     int
	GwenWins     int
	Draws        int
	JaneBadMoves int
	GwenBadMoves int
}

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {

	var memoryLock sync.RWMutex

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(2 * time.Second))
	r.Use(middleware.Throttle(100))

	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)
	filePath := filepath.Join(exeDir, "tic.log")
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	var results Results

	games := make(map[string]Game)
	for i := 0; i < 100; i++ {
		results.OpenGames++
		game := NewGame()
		games[game.Id] = game
	}

	r.Get("/tic", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("yes i am alive")
	})
	r.Get("/tic/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(1760400306)
	})
	r.Get("/tic/games", func(w http.ResponseWriter, r *http.Request) {
		memoryLock.RLock()
		defer memoryLock.RUnlock()
		var gamesSlice []Game
		for _, v := range games {
			gamesSlice = append(gamesSlice, v)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gamesSlice)
	})
	r.Get("/tic/results", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})
	r.Post("/tic/move", func(w http.ResponseWriter, r *http.Request) {
		var move Move
		if err := json.NewDecoder(r.Body).Decode(&move); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else if game, ok := games[move.Id]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "thats not an ID")
			switch move.User {
			case "jane":
				results.JaneBadMoves++
			case "gwen":
				results.GwenBadMoves++
			}
		} else if err := MoveResult(&game, &move); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad move")
			switch move.User {
			case "jane":
				results.JaneBadMoves++
			case "gwen":
				results.GwenBadMoves++
			}
		} else if res, err := json.Marshal(game); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			memoryLock.Lock()
			defer memoryLock.Unlock()
			games[move.Id] = game
			writeResult := false
			if game.Result == "win" && game.Winner == "jane" {
				results.JaneWins++
			} else if game.Result == "win" && game.Winner == "gwen" {
				results.GwenWins++
			} else if game.Result == "draw" {
				results.Draws++
			}
			if writeResult {
				file.Write(res)
				fmt.Fprintf(file, "\000")
				delete(games, move.Id)
				newGame := NewGame()
				games[newGame.Id] = newGame
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
		}
	})
	fs := http.FileServer(http.Dir(filepath.Join(exeDir, "dashboard")))
	r.Handle("/tic/dashboard", http.StripPrefix("/tic/dashboard", fs))
	log.Println("listening on :3333")
	return http.ListenAndServe(":3333", r)
}

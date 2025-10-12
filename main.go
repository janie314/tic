package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Game struct {
	Id      string     `json:"id"`
	Turn    string     `json:"turn"`
	Piece   string     `json:"piece"`
	Players []string   `json:"players"`
	State   [][]string `json:"state"`
	Winner  string     `json:"winner"`
}

type Move struct {
	Id   string `json:"id"`
	User string `json:"user"`
	Row  int    `json:"row"`
	Col  int    `json:"col"`
}

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {

	var lock sync.RWMutex

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	games := make(map[string]Game)
	for i := 0; i < 100; i++ {
		game := NewGame()
		games[game.Id] = game
	}

	r.Get("/tic", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("yes i am alive"))
	})
	r.Get("/tic/games", func(w http.ResponseWriter, r *http.Request) {
		lock.RLock()
		defer lock.RUnlock()
		var gamesSlice []Game
		for _, v := range games {
			gamesSlice = append(gamesSlice, v)
		}
		res, err := json.Marshal(gamesSlice)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
		}
	})
	r.Post("/tic/move", func(w http.ResponseWriter, r *http.Request) {
		var move Move
		if err := json.NewDecoder(r.Body).Decode(&move); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		game, ok := games[move.Id]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "thats not an ID")
			return
		}
		res := MoveResult(&game, &move)
		if res == Illegal {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	return http.ListenAndServe(":3333", r)
}

func NewGame() Game {
	return Game{
		Id:      uuid.NewString(),
		Turn:    randomPlayer(),
		Piece:   "X",
		Players: []string{"jane", "gwen"},
		State:   [][]string{{"_", "_", "_"}, {"_", "_", "_"}, {"_", "_", "_"}},
		Winner:  "",
	}
}

type MoveRes int

const (
	Win = 0
	Lose
	Draw
	Illegal
	OK
)

func MoveResult(game *Game, move *Move) MoveRes {
	if move.Row < 0 || move.Row > 2 || move.Col < 0 || move.Col > 2 {
		return Illegal
	}
	if game.State[move.Row][move.Col] != "_" {
		return Illegal
	}
	for _, winner := range []string{"X", "O"} {
		for i := 0; i < 3; i++ {
			if (game.State[i][0] == winner && game.State[i][1] == winner && game.State[i][2] == winner) || (game.State[0][i] == winner && game.State[1][i] == winner && game.State[2][i] == winner) {
				if winner == game.Turn {
					return Win
				} else {
					return Lose
				}
			}
		}
		if (game.State[0][0] == winner && game.State[1][1] == winner && game.State[2][2] == winner) || (game.State[0][2] == winner && game.State[1][1] == winner && game.State[2][0] == winner) {
			if winner == game.Turn {
				return Win
			} else {
				return Lose
			}
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if game.State[i][j] == "_" {
					return OK
				}
			}
		}
	}
	// TODO: could figure a draw a move earlier sometimes
	return Draw
}

func randomPlayer() string {
	return []string{"gwen", "jane"}[rand.Intn(2)]
}

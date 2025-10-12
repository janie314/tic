package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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
	for i := 0; i < 1; i++ {
		game := NewGame()
		games[game.Id] = game
	}

	r.Get("/tic", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("yes i am alive"))
	})
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("1760301245"))
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
			fmt.Fprint(w, err.Error())
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
		} else if game, ok := games[move.Id]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "thats not an ID")
			return
		} else if err := MoveResult(&game, &move); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else if res, err := json.Marshal(game); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			lock.Lock()
			defer lock.Unlock()
			games[move.Id] = game
			w.Write(res)
		}
	})
	return http.ListenAndServe(":3333", r)
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Game struct {
	Id    string     `json:"id"`
	Turn  string     `json:"turn"`
	Piece string     `json:"piece"`
	State [][]string `json:"state"`
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
	games["2ca973fc-b37d-4f3b-9e25-81360f430a84"] = Game{
		Id:    "2ca973fc-b37d-4f3b-9e25-81360f430a84",
		Turn:  "gwen",
		Piece: "X",
		State: [][]string{[]string{"_", "_", "_"}, []string{"_", "_", "_"}, []string{"_", "_", "_"}},
	}

	r.Get("/tic", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("yes i am alive"))
	})
	r.Get("/tic/challenges", func(w http.ResponseWriter, r *http.Request) {
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
			if _, err := w.Write(res); err != nil {
				log.Println("error /tic/challenges", err)
			}
		}
	})

	return http.ListenAndServe(":3333", r)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	r.Use(middleware.Timeout(2 * time.Second))
	r.Use(middleware.Throttle(100))

	filePath := "tic.log"
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
		w.Write([]byte("yes i am alive"))
	})
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("1760305657"))
	})
	r.Get("/tic/games", func(w http.ResponseWriter, r *http.Request) {
		memoryLock.RLock()
		defer memoryLock.RUnlock()
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
	r.Get("/tic/results", func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal(results)
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
			switch move.User {
			case "jane":
				results.JaneBadMoves++
			case "gwen":
				results.GwenBadMoves++
			}
			return
		} else if err := MoveResult(&game, &move); err != nil {
			switch move.User {
			case "jane":
				results.JaneBadMoves++
			case "gwen":
				results.GwenBadMoves++
			}
			w.WriteHeader(http.StatusBadRequest)
		} else if res, err := json.Marshal(game); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			memoryLock.Lock()
			defer memoryLock.Unlock()
			games[move.Id] = game
			if game.Result == "win" && game.Winner == "jane" {
				file.Write(res)
				fmt.Fprintf(file, "\000")
				delete(games, move.Id)
				newGame := NewGame()
				games[newGame.Id] = newGame
				results.JaneWins++
			} else if game.Result == "win" && game.Winner == "gwen" {
				file.Write(res)
				fmt.Fprintf(file, "\000")
				delete(games, move.Id)
				newGame := NewGame()
				games[newGame.Id] = newGame
				results.GwenWins++
			} else if game.Result == "draw" {
				file.Write(res)
				fmt.Fprintf(file, "\000")
				delete(games, move.Id)
				newGame := NewGame()
				games[newGame.Id] = newGame
				results.Draws++
			}
			w.Write(res)
		}
	})
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "dashboard"))
	FileServer(r, "/dashboard", filesDir)
	return http.ListenAndServe(":3333", r)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

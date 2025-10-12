package main

import (
	"errors"
	"math/rand"

	"github.com/google/uuid"
)

type Game struct {
	Id      string     `json:"id"`
	Turn    string     `json:"turn"`
	Piece   string     `json:"piece"`
	Players []string   `json:"players"`
	State   [][]string `json:"state"`
	Result  string     `json:"result"`
	Winner  string     `json:"winner"`
}

type Move struct {
	Id   string `json:"id"`
	User string `json:"user"`
	Row  int    `json:"row"`
	Col  int    `json:"col"`
}

type MoveRes int

const (
	Win = 0
	Lose
	Draw
	Illegal
	OK
)

func MoveResult(game *Game, move *Move) error {
	if game.Turn != move.User {
		return errors.New("not your turn")
	}
	if move.Row < 0 || move.Row > 2 || move.Col < 0 || move.Col > 2 {
		return errors.New("illegal move")
	}
	if game.State[move.Row][move.Col] != "_" {
		return errors.New("illegal move")
	}
	game.State[move.Row][move.Col] = game.Piece
	if err := switchTurn(game); err != nil {
		return nil
	}
	for _, winner := range []string{"X", "O"} {
		for i := 0; i < 3; i++ {
			if (game.State[i][0] == winner && game.State[i][1] == winner && game.State[i][2] == winner) || (game.State[0][i] == winner && game.State[1][i] == winner && game.State[2][i] == winner) {
				if winnerId, err := winnerID(game, winner == game.Piece); err != nil {
					return err
				} else {
					game.Winner = winnerId
					game.Result = "win"
					return nil
				}
			}
		}
		if (game.State[0][0] == winner && game.State[1][1] == winner && game.State[2][2] == winner) || (game.State[0][2] == winner && game.State[1][1] == winner && game.State[2][0] == winner) {
			if winnerId, err := winnerID(game, winner == game.Piece); err != nil {
				return err
			} else {
				game.Winner = winnerId
				game.Result = "win"
				return nil
			}
		}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if game.State[i][j] == "_" {
					return nil
				}
			}
		}
	}
	// TODO: could figure a draw a move earlier sometimes
	game.Result = "draw"
	return nil
}

func randomPlayer() string {
	return []string{"gwen", "jane"}[rand.Intn(2)]
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

func winnerID(game *Game, win bool) (string, error) {
	if win {
		return game.Turn, nil
	}
	for _, player := range game.Players {
		if player != game.Turn {
			return player, nil
		}
	}
	return "", nil
}

func switchTurn(game *Game) error {
	if game.Piece == "X" {
		game.Piece = "O"
	} else {
		game.Piece = "X"
	}
	for _, player := range game.Players {
		if player != game.Turn {
			game.Turn = player
			return nil
		}
	}
	return errors.New("bad bad")
}

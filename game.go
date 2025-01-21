package main

import (
	"fmt"
)

var boardChar = []rune{' ', 'X', 'O'}

type Game struct {
	Board      [boardSize]int `json:"board"`
	PlayerTurn int            `json:"playerTurn"`
}

func NewGame() *Game {
	return &Game{Board: [boardSize]int{}, PlayerTurn: player1}
}

func (g *Game) CheckMove(x int) bool {
	return g.Board[x - 1] == 0
}

func (g *Game) Play(x int) {

	if g.Board[x - 1] == 0 {
		g.Board[x - 1] = g.PlayerTurn
		g.PlayerTurn = 3 - g.PlayerTurn
	}

	g.PrintBoard()
}

func (g *Game) CheckWin() int {

	solutions := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // columns
		{0, 4, 8}, {2, 4, 6},            // diagonals
	}

	for _, solution := range solutions {
		if g.Board[solution[0]] != 0 && g.Board[solution[0]] == g.Board[solution[1]] && g.Board[solution[0]] == g.Board[solution[2]] {
			return g.Board[solution[0]]
		}
	}

	return 0
}

func (g *Game) PrintBoard() {
	fmt.Printf("\nBoard :")
	fmt.Printf("\n-------------\n")
	for i := 0; i < 9; i += 3 {
		fmt.Printf("| %c | %c | %c |\n", boardChar[g.Board[i]], boardChar[g.Board[i+1]], boardChar[g.Board[i+2]])
		fmt.Printf("-------------\n")
	}
}
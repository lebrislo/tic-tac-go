package main

import (
	"fmt"
)

type Game struct {
	Board      [9]int `json:"board"`
	PlayerTurn int `json:"playerTurn"`
}

func NewGame() *Game {
	return &Game{[9]int{0, 0, 0, 0, 0, 0,0, 0, 0}, 1}
}

func (g *Game) Play(x int) {
	fmt.Println("cell played", x - 1)
	if g.Board[x - 1] == 0 {
		g.Board[x - 1] = g.PlayerTurn
		g.PlayerTurn = 3 - g.PlayerTurn
	}

	g.PrintBoard()
}

func (g *Game) CheckWin() int {
	return 0
}

func (g *Game) PrintBoard() {
	fmt.Printf("\nBoard:")
	fmt.Printf("\n-------------\n")
	for i := 0; i < 9; i += 3 {
		fmt.Printf("| %d | %d | %d |\n", g.Board[i], g.Board[i+1], g.Board[i+2])
		fmt.Printf("-------------\n")
	}
}
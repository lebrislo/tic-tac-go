package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	port         = "8080"
	address      = "0.0.0.0:" + port
	retryDelay   = 1 * time.Second
	boardSize    = 9
	player1      = 1
	player2      = 2
	emptyCell    = 0
	messageDelim = '\n'
)

var (
	playerNumber   int
	playerIsServer bool
	game           Game
	connection     net.Conn
)

func main() {
	setupLogger()

	playerNumber = getPlayerNumber()
	playerIsServer = (playerNumber == player1)

	game = *NewGame()

	if playerIsServer {
		log.Println("Player 1 selected. Starting server...")
		startServerPlayer()
	} else {
		log.Println("Player 2 selected. Connecting to server...")
		startClientPlayer()
	}
}

// setupLogger initializes log settings.
func setupLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

// getPlayerNumber prompts the user to select a player number.
func getPlayerNumber() int {
	for {
		fmt.Print("Choose player number (1 or 2): ")
		var input int
		if _, err := fmt.Fscanln(os.Stdin, &input); err == nil && (input == player1 || input == player2) {
			return input
		}
		log.Println("Invalid input. Please enter 1 or 2.")
	}
}

// startServerPlayer sets up the server for Player 1.
func startServerPlayer() {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Println("Waiting for Player 2 to connect...")
	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("Failed to accept connection: %v", err)
	}
	connection = conn
	log.Println("Player 2 connected!")
	startGame()
}

// startClientPlayer connects to the server as Player 2.
func startClientPlayer() {
	for {
		conn, err := net.Dial("tcp", address)
		if err == nil {
			connection = conn
			log.Println("Connected to Player 1!")
			startGame()
			return
		}
		log.Printf("Connection failed. Retrying in %v...\n", retryDelay)
		time.Sleep(retryDelay)
	}
}

// startGame handles the main game loop.
func startGame() {
	defer connection.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go receiveMessages(ctx)

	if playerIsServer {
		askPlayerMove()
	}

	for {
		winner := game.CheckWin()
		if winner != emptyCell {
			log.Printf("Player %d wins!\n", winner)
			return
		}
	}
}

// receiveMessages listens for messages from the other player.
func receiveMessages(ctx context.Context) {
	reader := bufio.NewReader(connection)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			message, err := reader.ReadString(messageDelim)
			if err != nil {
				log.Printf("Connection closed: %v", err)
				return
			}
			handleReceivedMessage(message)
		}
	}
}

// handleReceivedMessage processes the received game state.
func handleReceivedMessage(message string) {
	if err := json.Unmarshal([]byte(message), &game); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}
	log.Printf("Player %d played.\n", 3-game.PlayerTurn)
	game.PrintBoard()
	askPlayerMove()
}

// askPlayerMove prompts the player for their move.
func askPlayerMove() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Type board cell for your move [1-%d]: ", boardSize)
		if !scanner.Scan() {
			log.Println("Input error. Please try again.")
			continue
		}
		input := scanner.Text()
		cell, err := strconv.Atoi(input)
		if err != nil || cell < 1 || cell > boardSize || !game.CheckMove(cell) {
			log.Println("Invalid move. Choose a valid cell.")
			continue
		}
		game.Play(cell)
		sendGameState()
		return
	}
}

// sendGameState sends the updated game state to the other player.
func sendGameState() {
	gameJSON, err := json.Marshal(&game)
	if err != nil {
		log.Printf("Failed to marshal game state: %v", err)
		return
	}
	if _, err := connection.Write(append(gameJSON, messageDelim)); err != nil {
		log.Printf("Failed to send game state: %v", err)
	}
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var playerNumber int
var playerIsServer bool

var game Game
var connection net.Conn

func main() {
	
	// Ask user to be player 1 or player 2
	fmt.Printf("Choose player number: 1 | 2 : ")

	_, err := fmt.Fscanln(os.Stdin, &playerNumber);
	if err != nil {
		fmt.Fprint(os.Stderr, "Error reading input\n")
		os.Exit(1)
	}

	playerIsServer = playerNumber == 1

	game = *NewGame()

	if playerIsServer {
		fmt.Println("You choose player 1")
		startServerPlayer()
	} else {
		fmt.Println("You choose player 2")
		startClientPlayer()
	}
}

func startServerPlayer() {
	// Écoute sur le port 8080
	ln, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Waiting for Player 2 to connect...")

	// Accepte une connexion entrante
	conn, err := ln.Accept()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accepting connection: %s\n", err.Error())
		return
	}

	connection = conn

	fmt.Println("Player 2 connected !")
	startGame()
}

func startClientPlayer() {
	var conn net.Conn
	var err error

	// Tentative de connexion au serveur (Player 1)
	for {
		conn, err = net.Dial("tcp", "0.0.0.0:8080")
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
		fmt.Println("Retrying connection...")
	}

	connection = conn
	fmt.Println("Connected to Player 1 !")
	startGame()
}

func startGame() {
	defer connection.Close()

	// Lancement d'une goroutine pour lire les messages reçus
	go func() {
		reader := bufio.NewReader(connection)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading: %s\n", err.Error())
				return
			}
			fmt.Printf("Message received: %s", message)
		}
	}()

	if playerIsServer {
		askPlayerMove()
	}

	for game.CheckWin() == 0{

	}
}

func askPlayerMove() {
	var cell int

	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Printf("\nType board cell next play [1-9]: ")
		scanner.Scan()
		input := scanner.Text()

		var err error
		cell, err = strconv.Atoi(input)

		if err != nil {
			fmt.Println("Invalid input. Please enter a number between 1 and 9.")
			continue
		}

		if cell >= 1 && cell <= 9 {
			break
		} else {
			fmt.Println("Invalid cell number. Please enter a number between 1 and 9.")
		}
	}

	game.Play(cell)

	gameJson, _ := json.Marshal(&game)
	message := string(gameJson) + "\n"
	_, err := connection.Write([]byte(message))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending message: %s\n", err.Error())
	}
}
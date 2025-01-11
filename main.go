package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

var playerNumber int
var playerIsServer bool

var myTurn bool
var board [3][3]int

func main() {
	
	// Ask user to be player 1 or player 2
	fmt.Printf("Choose player number: 1 | 2 : ")

	_, err := fmt.Fscanln(os.Stdin, &playerNumber);
	if err != nil {
		fmt.Fprint(os.Stderr, "Error reading input\n")
		os.Exit(1)
	}

	playerIsServer = playerNumber == 1

	board = [3][3]int{ {0, 0, 0}, {0, 0, 0}, {0, 0, 0} }

	if playerIsServer {
		fmt.Println("You choose player 1")
		myTurn = true
		startServerPlayer()
	} else {
		fmt.Println("You choose player 2")
		myTurn = false
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

	fmt.Println("Player 2 connected !")
	handleBidirectionalCommunication(conn)
}

func startClientPlayer() {
	var conn net.Conn
	var err error

	// Tentative de connexion au serveur (Player 1)
	for {
		conn, err = net.Dial("tcp", "0.0.0.0:8080")
		if err == nil {
			fmt.Println("Connected to Player 1 !")
			break
		}
		time.Sleep(1 * time.Second)
		fmt.Println("Retrying connection...")
	}

	handleBidirectionalCommunication(conn)
}

func handleBidirectionalCommunication(conn net.Conn) {
	defer conn.Close()

	// Lancement d'une goroutine pour lire les messages reçus
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading: %s\n", err.Error())
				return
			}
			fmt.Printf("Message received: %s", message)
			printBoard()
		}
	}()

	// Boucle principale pour envoyer des messages
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type a message to send (type 'quit' to exit):")
	for scanner.Scan() {
		message := scanner.Text()
		if message == "quit" {
			fmt.Println("Closing connection...")
			return
		}
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %s\n", err.Error())
			return
		}
	}
}

func printBoard() {
	fmt.Println("Board:")
	fmt.Printf("\n-------------\n")
	for i := 0; i < 3; i++ {
		fmt.Printf("|")
		for j := 0; j < 3; j++ {
			if board[i][j] == 1 {
				fmt.Printf(" X ")
			} else if board[i][j] == 2 {
				fmt.Printf(" O ")
			} else {
				fmt.Printf("   ")
			}
			fmt.Printf("|")
		}
		fmt.Printf("\n-------------\n")
	}
}
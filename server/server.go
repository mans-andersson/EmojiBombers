package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"
)

var gameState state
var id int

func main() {
	id = 0
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	rand.Seed(time.Now().UnixNano())
	gameState.Players = append(gameState.Players, player{XPos: rand.Intn(1000), YPos: rand.Intn(600)})
	playerID := id
	id = id + 1
	for {
		command := action{}
		netData, err := bufio.NewReader(c).ReadString('}')
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal([]byte(netData), &command)
		processCommand(playerID, &command)
		s, _ := json.Marshal(&gameState)
		// fmt.Println(gameState)
		// fmt.Println(string(s))
		c.Write(s)
	}
}

func processCommand(id int, command *action) {
	if command.Action == "up" {
		gameState.Players[id].YPos--
	} else if command.Action == "left" {
		gameState.Players[id].XPos--
	} else if command.Action == "down" {
		gameState.Players[id].YPos++
	} else if command.Action == "right" {
		gameState.Players[id].XPos++
	}
}

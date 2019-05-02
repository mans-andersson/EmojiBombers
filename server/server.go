package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
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
	// rand.Seed(time.Now().UnixNano())
	playerID := id
	id = id + 1
	gameState.Players = append(gameState.Players, player{ID: playerID, XPos: startingXPos[playerID], YPos: startingYPos[playerID]})
	go sendState(c)
	for {
		command := action{}
		netData, err := bufio.NewReader(c).ReadString(';')
		if err != nil {
			// fmt.Println(err)
		}
		if len(netData) > 0 {
			netData = netData[:len(netData)-1]
		}
		err = json.Unmarshal([]byte(netData), &command)
		go processCommand(playerID, &command)
	}
}

func sendState(c net.Conn) {
	for {
		time.Sleep(16 * time.Millisecond)
		c.Write(stateToJsonTransmission())
	}
}

func stateToJsonTransmission() []byte {
	result := ([]byte)(";")
	s, err := json.Marshal(&gameState)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	result = append(result, s...)
	result = append(result, ([]byte)(";")...)
	return result
}

func processCommand(id int, command *action) {
	if command.Action == "up" {
		movePlayer(id, "up")
	} else if command.Action == "left" {
		movePlayer(id, "left")
	} else if command.Action == "down" {
		movePlayer(id, "down")
	} else if command.Action == "right" {
		movePlayer(id, "right")
	}
}

func movePlayer(id int, direction string) {
	newX, newY := newPosition(id, direction)
	for _, player := range gameState.Players {
		if(distance(newX, newY, player.XPos, player.YPos) <= 100) && player.ID != id {
			return
		}
	}
	gameState.Players[id].XPos = newX
	gameState.Players[id].YPos = newY
}


func newPosition(id int, direction string) (int, int){
	oldX, oldY := gameState.Players[id].XPos, gameState.Players[id].YPos
	if direction == "up" {
		return oldX, oldY - 2
	} else if direction == "left" {
		return oldX - 2, oldY
	} else if direction == "down" {
		return oldX, oldY + 2
	} else {
		return oldX + 2, oldY
	}
}

func distance(x1 int, y1 int, x2 int, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x1 - x2), 2) + math.Pow(float64(y1 - y2), 2))
}

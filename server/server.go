package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net"
	"time"
)

var gameState state
var id int

func main() {
	gameState.Winner = -1
	addBlockades()
	initializePlayers()
	initializeBombs()
	initializePlacedBombs()
	initializeExplosions()
	go spawnBombs()
	go checkDamage()
	go checkVictory()
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
	if id > 3 {
		return
	}
	playerID := id
	gameState.Players[playerID].Spawned = true
	id++
	go sendState(c)
	for {
		time.Sleep(10 * time.Millisecond)
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

func addBlockades() {
	for i := 0; i <= 800; i += 200 {
		for j := 0; j <= 800; j += 200 {
			gameState.Blockades = append(gameState.Blockades, blockade{XPos: i, YPos: j})
		}
	}
}

func initializePlayers() {
	for i := 0; i < 4; i++ {
		gameState.Players = append(gameState.Players, player{ID: i, Lives: 3})
	}
	initializeStartingPositions()
}

func initializeStartingPositions() {
	startingXPos := []int{100, 700, 100, 700}
	startingYPos := []int{100, 700, 700, 100}
	for i := 0; i < 4; i++ {
		gameState.Players[i].XPos = startingXPos[i]
		gameState.Players[i].YPos = startingYPos[i]
	}
}

func initializeBombs() {
	for i := 100; i < 800; i += 200 {
		for j := 100; j < 800; j += 200 {
			gameState.Bombs = append(gameState.Bombs, bomb{Spawned: false, XPos: i, YPos: j})
		}
	}
}

func initializePlacedBombs() {
	for i := 0; i < 4; i++ {
		gameState.PlacedBombs = append(gameState.PlacedBombs, bomb{Spawned: false, XPos: 0, YPos: 0})
	}
}

func initializeExplosions() {
	for i := 0; i < 4; i++ {
		gameState.Explosions = append(gameState.Explosions, bomb{Spawned: false, XPos: 0, YPos: 0})
	}
}

func spawnBombs() {
	for {
		time.Sleep(3000 * time.Millisecond)
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(len(gameState.Bombs))
		gameState.Bombs[i].Spawned = true
	}
}

func checkDamage() {
	for {
		time.Sleep(16 * time.Millisecond)
		for id, player := range gameState.Players {
			for _, exp := range gameState.Explosions {
				if(distance(player.XPos, player.YPos, exp.XPos, exp.YPos) <= 150 && !player.DamageTaken && exp.Spawned) {
					gameState.Players[id].DamageTaken = true
					go playerDamaged(id)
				}
			}
		}
	}
}

func checkVictory() {
	for {
		time.Sleep(1000 * time.Millisecond)
		if(countSpawnedPlayers() > 1) {
			winner := getWinner()
			if(winner != -1) {
				time.Sleep(1000 * time.Millisecond)
				gameState.Winner = winner
			}
		}
	}
}

func getWinner() int {
	aliveCount := 0
	potentialWinner := 0
	for i, p := range gameState.Players {
		if(!p.Dead && p.Spawned) {
			aliveCount++
			potentialWinner = i
		}
	}
	if aliveCount == 1 {
		return potentialWinner
	}
	return -1
}

func countSpawnedPlayers() int {
	counter := 0
	for _, p := range gameState.Players {
		if(p.Spawned) {
			counter++
		}
	}
	return counter
}

func playerDamaged(id int) {
	gameState.Players[id].Lives--
	if(gameState.Players[id].Lives == 0) {
		gameState.Players[id].Dead = true
	}
	time.Sleep(1500 * time.Millisecond)
	gameState.Players[id].DamageTaken = false
}

func stateToJsonTransmission() []byte {
	result := ([]byte)(";")
	s, err := json.Marshal(&gameState)
	if err != nil {
		return []byte{}
	}
	result = append(result, s...)
	result = append(result, ([]byte)(";")...)
	return result
}

func processCommand(id int, command *action) {
	if gameState.Players[id].Dead {
		return
	}
	if command.Action == "up" {
		movePlayer(id, "up")
	} else if command.Action == "left" {
		movePlayer(id, "left")
	} else if command.Action == "down" {
		movePlayer(id, "down")
	} else if command.Action == "right" {
		movePlayer(id, "right")
	} else if command.Action == "space" {
		placeBomb(id, gameState.Players[id].XPos, gameState.Players[id].YPos)
	}
}

func placeBomb(id int, x int, y int) {
	if gameState.PlacedBombs[id].Spawned || gameState.Players[id].BombCount == 0 {
		return
	}
	gameState.PlacedBombs[id].Spawned = true
	gameState.PlacedBombs[id].XPos = x
	gameState.PlacedBombs[id].YPos = y
	gameState.Players[id].BombCount--
	go explodeBomb(id)
}

func explodeBomb(id int) {
	time.Sleep(2000 * time.Millisecond)
	gameState.PlacedBombs[id].Spawned = false
	gameState.Explosions[id].XPos = gameState.PlacedBombs[id].XPos
	gameState.Explosions[id].YPos = gameState.PlacedBombs[id].YPos
	gameState.Explosions[id].Spawned = true
	time.Sleep(2000 * time.Millisecond)
	gameState.Explosions[id].Spawned = false

}

func movePlayer(id int, direction string) {
	newX, newY := newPosition(id, direction)
	for _, player := range gameState.Players {
		if(distance(newX, newY, player.XPos, player.YPos) <= 100) && player.ID != id && player.Spawned {
			return
		}
	}
	for _, blockade := range gameState.Blockades {
		if(distance(newX, newY, blockade.XPos, blockade.YPos) <= 75) {
			return
		}
	}
	for i, bomb := range gameState.Bombs {
		if(bomb.Spawned) {
			if(distance(newX, newY, bomb.XPos, bomb.YPos) <= 75) {
				gameState.Bombs[i].Spawned = false
				gameState.Players[id].BombCount++
			}
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

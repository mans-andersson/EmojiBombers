// Written by: Måns Andersson

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

// Holds the state of the game
var gameState state

// Variable used to give players unique ID's
var id int

func main() {
	addBlockades()
	initializePlayers()
	resetGame()
	go waitForStart()
	id = 0
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		// Accept new connections and handle them in parallel with handleConnection
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
		// This means that the maximum number of players has been reached
		return
	}
	playerID := id
	// Set player as spawned
	gameState.Players[playerID].Spawned = true
	id++
	go sendState(c)
	for {
		// Contiously read data sent from the client and process it
		time.Sleep(10 * time.Millisecond)
		command := action{}
		netData, _ := bufio.NewReader(c).ReadString(';')
		if len(netData) > 0 {
			netData = netData[:len(netData)-1]
		}
		json.Unmarshal([]byte(netData), &command)
		go processCommand(playerID, &command)
	}
}

// Continously send the game state on the connection
func sendState(c net.Conn) {
	for {
		time.Sleep(16 * time.Millisecond)
		c.Write(stateToJsonTransmission())
	}
}

// Add all the blockades to the game state
func addBlockades() {
	for i := 0; i <= 800; i += 200 {
		for j := 0; j <= 800; j += 200 {
			gameState.Blockades = append(gameState.Blockades, blockade{XPos: i, YPos: j})
		}
	}
}

// Wait until at least two players are connected
// After that bombs and snowflakes start spawning
// We also start to check for player damage and the victory condition
func waitForStart() {
	for {
		time.Sleep(1000 * time.Millisecond)
		if(countSpawnedPlayers() > 1) {
			go spawnBombs()
			go spawnSnowflakes()
			go checkDamage()
			go checkVictory()
			return
		}
	}
}

// Fills up the list with four players
// Spawned will be false at this state so they won't be rendered
func initializePlayers() {
	for i := 0; i < 4; i++ {
		gameState.Players = append(gameState.Players, player{ID: i})
	}
	initializeStartingPositions()
}

func initializePlayerLives() {
	for i := 0; i < 4; i++ {
		gameState.Players[i].Lives = 3
		gameState.Players[i].Dead = false
	}
}

// Sets the four players to the starting positions in the four corners
func initializeStartingPositions() {
	startingXPos := []int{100, 700, 100, 700}
	startingYPos := []int{100, 700, 700, 100}
	for i := 0; i < 4; i++ {
		gameState.Players[i].XPos = startingXPos[i]
		gameState.Players[i].YPos = startingYPos[i]
	}
}

// Adds all the default bomb positions to the game state
// Spawned will be false for all of them at this stage so they won't be rendered
func initializeBombs() {
	bombs := []bomb{}
	for i := 100; i < 800; i += 200 {
		for j := 100; j < 800; j += 200 {
			bombs = append(bombs, bomb{Spawned: false, XPos: i, YPos: j})
		}
	}
	gameState.Bombs = bombs
}

// Adds all the default snowflake positions to the game state
// Spawned will be false for all of them at this stage so they won't be rendered
func initializeSnowflakes() {
	snowflakes := []bomb{}
	for i := 100; i < 800; i += 200 {
		for j := 100; j < 800; j += 200 {
			snowflakes = append(snowflakes, bomb{Spawned: false, XPos: i, YPos: j})
		}
	}
	gameState.Snowflakes = snowflakes
}

// Initializes four placed bombs
// Spawned will be set to true and their positions will be changed when they are
// actually placed
func initializePlacedBombs() {
	placedBombs := []bomb{}
	for i := 0; i < 4; i++ {
		placedBombs = append(placedBombs, bomb{Spawned: false, XPos: 0, YPos: 0})
	}
	gameState.PlacedBombs = placedBombs
}


// Initializes four explosions
// Spawned will be set to true and their positions will be changed when they are
// a player blows a bomb
func initializeExplosions() {
	explosions := []bomb{}
	for i := 0; i < 4; i++ {
		explosions = append(explosions, bomb{Spawned: false, XPos: 0, YPos: 0})
	}
	gameState.Explosions = explosions
}


// Spawn a bomb in a random default bomb position every 4 seconds
func spawnBombs() {
	for {
		time.Sleep(4000 * time.Millisecond)
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(len(gameState.Bombs))
		gameState.Bombs[i].Spawned = true
	}
}

// Spawn a snowflake in a random bomb position every 4 seconds
func spawnSnowflakes() {
	for {
		time.Sleep(10000 * time.Millisecond)
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(len(gameState.Snowflakes))
		gameState.Snowflakes[i].Spawned = true
	}
}

// Check if a player collides with an explosion
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


// Check if there is only one player alive, in that case they are the winner
// Wait for 5 seconds then reset the game
func checkVictory() {
	for {
		time.Sleep(1000 * time.Millisecond)
		if(countSpawnedPlayers() > 1) {
			winner := getWinner()
			if(winner != -1) {
				time.Sleep(1000 * time.Millisecond)
				gameState.Winner = winner
				time.Sleep(5000 * time.Millisecond)
				resetGame()
			}
		}
	}
}

// Initialize the default state of the game
func resetGame() {
	gameState.Winner = -1
	initializeStartingPositions()
	initializePlayerLives()
	initializeExplosions()
	initializeBombs()
	initializeSnowflakes()
	initializePlacedBombs()
}

// Look to see if there is only one player alive
// In that case return their ID. Otherwise return -1
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

// If a player has been damaged decrement their lives and check if they're dead
// If they are not dead they will be invulnerable for 1.5 seconds before
// they can take damage again
func playerDamaged(id int) {
	gameState.Players[id].Lives--
	if(gameState.Players[id].Lives == 0) {
		gameState.Players[id].Dead = true
	}
	time.Sleep(1500 * time.Millisecond)
	gameState.Players[id].DamageTaken = false
}

// Convert the current game state to a json byte array
func stateToJsonTransmission() []byte {
	result := ([]byte)(";")
	s, err := json.Marshal(&gameState)
	if err != nil {
		return []byte{}
	}
	result = append(s, ([]byte)(";")...)
	return result
}


// Process a command received from the client and act accordingly
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


// Place a bomb at the specified position
// The id is the id of the player who is placing it
// A player can only have one placed bomb at a time
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

// Wait for a short while before a specific placed bomb gets converted to an
// explosion
func explodeBomb(id int) {
	time.Sleep(2000 * time.Millisecond)
	gameState.PlacedBombs[id].Spawned = false
	gameState.Explosions[id].XPos = gameState.PlacedBombs[id].XPos
	gameState.Explosions[id].YPos = gameState.PlacedBombs[id].YPos
	gameState.Explosions[id].Spawned = true
	time.Sleep(2000 * time.Millisecond)
	gameState.Explosions[id].Spawned = false
}


// Stun all players except the one with the specified id
func stunPlayers(id int) {
	for i, p := range gameState.Players {
		if(p.ID != id) {
			gameState.Players[i].Stunned = true
		}
	}
	time.Sleep(5000 * time.Millisecond)
	for i, p := range gameState.Players {
		if(p.ID != id) {
			gameState.Players[i].Stunned = false
		}
	}
}

func movePlayer(id int, direction string) {
	if(gameState.Players[id].Stunned) {
		// The player is stunned, do nothing
		return
	}
	newX, newY := newPosition(id, direction)
	if(newX > 800 || newX < 0 || newY > 800 || newY < 0) {
		// The player is trying to move outside of the screen, do nothing
		return
	}
	// Check for collisions with other players
	for _, player := range gameState.Players {
		if(distance(newX, newY, player.XPos, player.YPos) <= 100) && player.ID != id && player.Spawned {
			// If colliding - do nothing
			return
		}
	}
	// Check for collisions with blockades
	for _, blockade := range gameState.Blockades {
		if(distance(newX, newY, blockade.XPos, blockade.YPos) <= 75) {
			// If colliding - do nothing
			return
		}
	}
	// Check for collisions with bombs, if collided - pick it up
	for i, bomb := range gameState.Bombs {
		if(bomb.Spawned) {
			if(distance(newX, newY, bomb.XPos, bomb.YPos) <= 75) {
				gameState.Bombs[i].Spawned = false
				gameState.Players[id].BombCount++
			}
		}
	}
	// Check for collisions with snowflakes, if collided - stun the other players
	for i, snowflake := range gameState.Snowflakes {
		if(snowflake.Spawned) {
			if(distance(newX, newY, snowflake.XPos, snowflake.YPos) <= 75) {
				gameState.Snowflakes[i].Spawned = false
				go stunPlayers(id)
			}
		}
	}
	// Confirm the movement
	gameState.Players[id].XPos = newX
	gameState.Players[id].YPos = newY
}


// Calculate what might become the new position of the player based on
// the direction of movement
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


// The distance formula
func distance(x1 int, y1 int, x2 int, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x1 - x2), 2) + math.Pow(float64(y1 - y2), 2))
}

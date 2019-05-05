package main

type state struct {
	Players []player `json:"players"`
	Blockades []blockade `json:"blockades"`
	Bombs []bomb  `json:"bombs"`
	PlacedBombs []bomb `json:"placed_bombs"`
	Explosions []bomb `json:"explosions"`
	GameStarted bool `json:"game_started"`
	Winner player `json:"winner"`
}

type player struct {
	ID int `json:"id"`
	XPos int `json:"x_pos"`
	YPos int `json:"y_pos"`
	Spawned bool `json:"spawned"`
	BombCount int `json:"bomb_count"`
	Lives int `json:"lives"`
	DamageTaken bool `json:"damage_taken"`
	Dead bool `json:"dead"`
}

type bomb struct {
	Spawned bool `json:"spawned"`
	XPos int `json:"x_pos"`
	YPos int `json:"y_pos"`
}

type blockade struct {
	XPos int `json:"x_pos"`
	YPos int `json:"y_pos"`
}

type action struct {
	Action string `json:"action"`
}

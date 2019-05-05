package main

type state struct {
	Players []player `json:"players"`
	Blockades []blockade `json:"blockades"`
	Bombs []bomb  `json:"bombs"`
}

type player struct {
	ID int `json:"id"`
	XPos int `json:"x_pos"`
	YPos int `json:"y_pos"`
	BombCount int `json:"bomb_count"`
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

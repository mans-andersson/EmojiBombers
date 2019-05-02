package main

type state struct {
	Players []player `json:"players"`
}

type player struct {
	ID int `json:"id"`
	XPos int `json:"x_pos"`
	YPos int `json:"y_pos"`
}

type action struct {
	Action string `json:"action"`
}

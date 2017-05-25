package main

import (
	"log"

	"github.com/zenja/mario/game"
)

var G *game.Game

func quit() {
	G.Quit()
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer quit()

	G = game.NewGame()
	G.LoadLevel("assets/levels/level1.txt")
	G.StartGameLoop()
}

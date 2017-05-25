package main

import (
	"log"

	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/level"
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

	G.ClearScreen()
	l.Draw(G, 0, 0)
	G.ShowScreen()

	G.Delay(3000)
}

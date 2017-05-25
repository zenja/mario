package main

import (
	"log"

	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/level"
)

var G *graphic.Graphic

func quit() {
	G.DestroyAndQuit()
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer quit()

	G = graphic.New()
	l := level.ParseLevelFromFile("assets/levels/level1.txt")

	G.ClearScreen()
	l.Draw(G, 0, 0)
	G.ShowScreen()

	G.Delay(3000)
}

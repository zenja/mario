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
	levelArr := [][]byte{
		[]byte{'o', 'G', 'o'},
		[]byte{'G', 'o', 'G'},
		[]byte{'o', 'G', 'o'},
	}
	l := level.ParseLevel(levelArr, G.TileRegistry)

	G.ClearScreen()
	l.Draw(G, 0, 0)
	G.ShowScreen()

	G.Delay(3000)
}

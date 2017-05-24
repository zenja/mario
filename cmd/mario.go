package main

import (
	"log"

	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/object"
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
	o := object.NewSingleTileObject(graphic.TILE_TYPE_GROUD, 100, 100)
	G.ClearScreen()
	o.Draw(G, 230, 230)
	G.ShowScreen()

	G.Delay(5000)
}

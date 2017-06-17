package main

import (
	"log"

	"runtime"

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

	// this will prevent window not responding
	runtime.LockOSThread()

	G = game.NewGame()
	G.Init()
	G.StartGameLoop()
}

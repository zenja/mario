package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// this will prevent window not responding
	runtime.LockOSThread()

	// Ask for hero user ID
	fmt.Println("Please enter your NT username: ")
	reader := bufio.NewReader(os.Stdin)
	uid, err := reader.ReadString('\n')
	uid = strings.TrimSuffix(uid, "\n")
	if err != nil {
		log.Fatal(err)
	}

	G = game.NewGame()
	G.Init(uid)
	G.StartGameLoop()
}

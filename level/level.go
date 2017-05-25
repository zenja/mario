package level

import (
	"bufio"
	"log"
	"os"

	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/object"
)

type Level struct {
	Objects []object.Object
}

func ParseLevel(arr [][]byte) *Level {
	var objs []object.Object
	var currentX, currentY int32
	for i, arrRow := range arr {
		currentX = 0
		for j, _ := range arrRow {
			switch arr[i][j] {
			// Ground
			case 'G':
				objs = append(objs, object.NewSingleTileObject(graphic.TILE_TYPE_GROUD, currentX, currentY))
			// Hero
			case 'H':
				objs = append(objs, object.NewHero(currentX, currentY))
			}
			currentX += graphic.TILE_SIZE
		}
		currentY += graphic.TILE_SIZE
	}
	return &Level{Objects: objs}
}

func ParseLevelFromFile(filename string) *Level {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file %s", filename)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var arr [][]byte
	for scanner.Scan() {
		arr = append(arr, []byte(scanner.Text()))
	}
	return ParseLevel(arr)
}

func (l *Level) Draw(g *graphic.Graphic, xCamStart, yCamStart int32) {
	for _, o := range l.Objects {
		o.Draw(g, xCamStart, yCamStart)
	}
}

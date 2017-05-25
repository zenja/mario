package level

import (
	"bufio"
	"log"
	"os"

	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/object"
)

type Level struct {
	Objects [][]object.Object
}

func ParseLevel(arr [][]byte) *Level {
	var objRows [][]object.Object
	var currentX, currentY int32
	for i, arrRow := range arr {
		currentX = 0
		var objRow []object.Object
		for j, _ := range arrRow {
			switch arr[i][j] {
			// Ground
			case 'G':
				objRow = append(objRow, object.NewSingleTileObject(graphic.TILE_TYPE_GROUD, currentX, currentY))
			}
			currentX += graphic.TILE_SIZE
		}
		objRows = append(objRows, objRow)
		currentY += graphic.TILE_SIZE
	}
	return &Level{Objects: objRows}
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
	for i := range l.Objects {
		for j := range l.Objects[i] {
			l.Objects[i][j].Draw(g, xCamStart, yCamStart)
		}
	}
}

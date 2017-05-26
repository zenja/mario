package level

import (
	"bufio"
	"log"
	"os"

	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/object"
	"github.com/zenja/mario/vector"
)

type Level struct {
	Objects  []object.Object
	ObstMngr *ObstacleManager
}

func ParseLevel(arr [][]byte, resourceRegistry map[graphic.ResourceID]graphic.Resource) *Level {
	var objs []object.Object
	obstMngr := NewObstacleManager(len(arr[0]), len(arr))

	var currentX, currentY int32
	for i, arrRow := range arr {
		currentX = 0
		for j := range arrRow {
			switch arr[i][j] {
			// Ground
			case 'G':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD]
				objs = append(objs, object.NewSingleTileObject(resource, currentX, currentY))
				// ground is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})
			// Hero
			case 'H':
				objs = append(objs, object.NewHero(currentX, currentY, resourceRegistry))
			}
			currentX += graphic.TILE_SIZE
		}
		currentY += graphic.TILE_SIZE
	}
	return &Level{
		Objects:  objs,
		ObstMngr: obstMngr,
	}
}

func ParseLevelFromFile(filename string, resourceRegistry map[graphic.ResourceID]graphic.Resource) *Level {
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
	return ParseLevel(arr, resourceRegistry)
}

func (l *Level) Draw(g *graphic.Graphic, xCamStart, yCamStart int32) {
	for _, o := range l.Objects {
		o.Draw(g, xCamStart, yCamStart)
	}
}

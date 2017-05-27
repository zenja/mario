package level

import (
	"bufio"
	"log"
	"os"

	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type Level struct {
	Objects  []Object
	ObstMngr *ObstacleManager
	Hero     Object
	numTiles vector.Vec2D
}

func ParseLevel(arr [][]byte, resourceRegistry map[graphic.ResourceID]graphic.Resource) *Level {
	var objs []Object
	numTiles := vector.Vec2D{int32(len(arr[0])), int32(len(arr))}
	obstMngr := NewObstacleManager(len(arr[0]), len(arr))
	var hero Object

	var currentPos vector.Pos
	for i, arrRow := range arr {
		currentPos.X = 0
		for j := range arrRow {
			switch arr[i][j] {
			// Ground
			case 'G':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD]
				objs = append(objs, NewSingleTileObject(resource, currentPos, ZINDEX_0))
				// ground is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})
			// Hero
			case 'H':
				if hero != nil {
					log.Fatal("more than one hero found")
				}
				hero = NewHero(currentPos, resourceRegistry)
				objs = append(objs, hero)
			}
			currentPos.X += graphic.TILE_SIZE
		}
		currentPos.Y += graphic.TILE_SIZE
	}

	if hero == nil {
		log.Fatal("no hero found when parsing level")
	}

	return &Level{
		Objects:  objs,
		ObstMngr: obstMngr,
		Hero:     hero,
		numTiles: numTiles,
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

func (l *Level) Draw(g *graphic.Graphic, camPos vector.Pos) {
	var zIndexObjs [ZINDEX_NUM][]Object
	// draw lowest z-index, bookkeeping higher z-index for later rendering
	for _, o := range l.Objects {
		z := o.GetZIndex()
		if z == ZINDEX_0 {
			o.Draw(g, camPos)
		} else {
			zIndexObjs[z] = append(zIndexObjs[z], o)
		}
	}
	// render higher z-index one-by-one
	for _, objs := range zIndexObjs {
		if len(objs) > 0 {
			for _, o := range objs {
				o.Draw(g, camPos)
			}
		}
	}
}

func (l *Level) GetLevelWidth() int32 {
	return l.numTiles.X * graphic.TILE_SIZE
}

func (l *Level) GetLevelHeight() int32 {
	return l.numTiles.Y * graphic.TILE_SIZE
}

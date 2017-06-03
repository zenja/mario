package level

import (
	"bufio"
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type Level struct {
	Objects  []Object
	ObstMngr *ObstacleManager
	Hero     Object
	BGColor  sdl.Color
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
			// Brick
			case 'B':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_BRICK]
				objs = append(objs, NewSingleTileObject(resource, currentPos, ZINDEX_0))
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})

			// Ground with left grass
			case 'L':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD_GRASS_LEFT]
				objs = append(objs, NewSingleTileObject(resource, currentPos, ZINDEX_0))
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})

			// Ground with mid grass
			case 'G':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD_GRASS_MID]
				objs = append(objs, NewSingleTileObject(resource, currentPos, ZINDEX_0))
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})

			// Ground with right grass
			case 'R':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD_GRASS_RIGHT]
				objs = append(objs, NewSingleTileObject(resource, currentPos, ZINDEX_0))
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})

			// Inner ground in middle
			case 'I':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD_INNER_MID]
				objs = append(objs, NewSingleTileObject(resource, currentPos, ZINDEX_0))
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})

			// Myth box
			case 'M':
				objs = append(objs, NewMythBox(currentPos, 1, resourceRegistry))
				// this is obstacle
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
		BGColor:  sdl.Color{204, 237, 255, 255},
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

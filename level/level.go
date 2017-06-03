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
	TileObjects [][]Object
	ObstMngr    *ObstacleManager
	Hero        Object
	BGColor     sdl.Color
	NumTiles    vector.Vec2D
}

func ParseLevel(arr [][]byte, resourceRegistry map[graphic.ResourceID]graphic.Resource) *Level {
	var tileObjs [][]Object
	numTiles := vector.Vec2D{int32(len(arr[0])), int32(len(arr))}
	obstMngr := NewObstacleManager(len(arr[0]), len(arr))
	var hero Object

	// init tileObjs array
	for i := 0; i < int(numTiles.Y); i++ {
		tileObjs = append(tileObjs, make([]Object, numTiles.X))
	}

	var currentPos vector.Pos
	for i, arrRow := range arr {
		currentPos.X = 0
		for j := range arrRow {
			switch arr[i][j] {
			// Brick
			case 'B':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_BRICK]
				tileObjs[i][j] = NewSingleTileObject(resource, currentPos, ZINDEX_0)
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})

			// Ground with left grass
			case 'L':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD_GRASS_LEFT]
				tileObjs[i][j] = NewSingleTileObject(resource, currentPos, ZINDEX_0)
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})

			// Ground with mid grass
			case 'G':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD_GRASS_MID]
				tileObjs[i][j] = NewSingleTileObject(resource, currentPos, ZINDEX_0)
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})

			// Ground with right grass
			case 'R':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD_GRASS_RIGHT]
				tileObjs[i][j] = NewSingleTileObject(resource, currentPos, ZINDEX_0)
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})

			// Inner ground in middle
			case 'I':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD_INNER_MID]
				tileObjs[i][j] = NewSingleTileObject(resource, currentPos, ZINDEX_0)
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})

			// Myth box
			case 'M':
				tileObjs[i][j] = NewMythBox(currentPos, 1, resourceRegistry)
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(j), int32(i)})

			// Hero
			case 'H':
				if hero != nil {
					log.Fatal("more than one hero found")
				}
				hero = NewHero(currentPos, resourceRegistry)
			}
			currentPos.X += graphic.TILE_SIZE
		}
		currentPos.Y += graphic.TILE_SIZE
	}

	if hero == nil {
		log.Fatal("no hero found when parsing level")
	}

	return &Level{
		TileObjects: tileObjs,
		ObstMngr:    obstMngr,
		Hero:        hero,
		BGColor:     sdl.Color{204, 237, 255, 255},
		NumTiles:    numTiles,
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
	for i := 0; i < int(l.NumTiles.Y); i++ {
		for j := 0; j < int(l.NumTiles.X); j++ {
			o := l.TileObjects[i][j]
			if o == nil {
				continue
			}

			z := o.GetZIndex()
			if z == ZINDEX_0 {
				o.Draw(g, camPos)
			} else {
				zIndexObjs[z] = append(zIndexObjs[z], o)
			}
		}
	}

	// put non-tile objects into render queue
	if l.Hero.GetZIndex() == ZINDEX_0 {
		log.Fatal("hero's z-index cannot be lowest")
	}
	zIndexObjs[l.Hero.GetZIndex()] = append(zIndexObjs[l.Hero.GetZIndex()], l.Hero)

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
	return l.NumTiles.X * graphic.TILE_SIZE
}

func (l *Level) GetLevelHeight() int32 {
	return l.NumTiles.Y * graphic.TILE_SIZE
}

package level

import (
	"bufio"
	"log"
	"os"

	"container/list"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type Level struct {
	// Public
	TileObjects      [][]Object
	ObstMngr         *ObstacleManager
	Hero             Object
	BGColor          sdl.Color
	NumTiles         vector.Vec2D // NOTE: X, Y is TID
	ResourceRegistry map[graphic.ResourceID]graphic.Resource

	// Private
	effects *list.List
}

func ParseLevel(arr [][]byte, resourceRegistry map[graphic.ResourceID]graphic.Resource) *Level {
	// NOTE: index is tid.X, tid.Y
	var tileObjs [][]Object

	numTiles := vector.Vec2D{int32(len(arr[0])), int32(len(arr))}
	obstMngr := NewObstacleManager(len(arr[0]), len(arr))
	var hero Object

	// init tileObjs array
	for i := 0; i < int(numTiles.X); i++ {
		tileObjs = append(tileObjs, make([]Object, numTiles.Y))
	}

	var currentPos vector.Pos
	for tidY := 0; tidY < int(numTiles.Y); tidY++ {
		currentPos.X = 0
		for tidX := 0; tidX < int(numTiles.X); tidX++ {
			// note that arr's index is not TID, need reverse
			switch arr[tidY][tidX] {
			// Brick
			case 'B':
				mainRes := resourceRegistry[graphic.RESOURCE_TYPE_BRICK]
				pieceRes := resourceRegistry[graphic.RESOURCE_TYPE_BRICK_PIECE]
				tileObjs[tidX][tidY] = NewBreakableTileObject(mainRes, pieceRes, currentPos, ZINDEX_0)
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(tidX), int32(tidY)})

			// Ground with left grass
			case 'L':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD_GRASS_LEFT]
				tileObjs[tidX][tidY] = NewSingleTileObject(resource, currentPos, ZINDEX_0)
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(tidX), int32(tidY)})

			// Ground with mid grass
			case 'G':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD_GRASS_MID]
				tileObjs[tidX][tidY] = NewSingleTileObject(resource, currentPos, ZINDEX_0)
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(tidX), int32(tidY)})

			// Ground with right grass
			case 'R':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD_GRASS_RIGHT]
				tileObjs[tidX][tidY] = NewSingleTileObject(resource, currentPos, ZINDEX_0)
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(tidX), int32(tidY)})

			// Inner ground in middle
			case 'I':
				resource := resourceRegistry[graphic.RESOURCE_TYPE_GROUD_INNER_MID]
				tileObjs[tidX][tidY] = NewSingleTileObject(resource, currentPos, ZINDEX_0)
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(tidX), int32(tidY)})

			// Myth box
			case 'M':
				tileObjs[tidX][tidY] = NewMythBox(currentPos, 1, resourceRegistry)
				// this is obstacle
				obstMngr.AddTileObst(vector.TileID{int32(tidX), int32(tidY)})

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
		TileObjects:      tileObjs,
		ObstMngr:         obstMngr,
		Hero:             hero,
		BGColor:          sdl.Color{204, 237, 255, 255},
		NumTiles:         numTiles,
		ResourceRegistry: resourceRegistry,
		effects:          list.New(),
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
	for i := 0; i < int(l.NumTiles.X); i++ {
		for j := 0; j < int(l.NumTiles.Y); j++ {
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

	// render effects and remove finished effects
	var finishedEffs []*list.Element
	for e := l.effects.Front(); e != nil; e = e.Next() {
		eff, ok := e.Value.(Effect)
		if !ok {
			log.Fatalf("eff is not an effect object: %T", e.Value)
		}
		eff.UpdateAndDraw(g, camPos, sdl.GetTicks())

		if eff.Finished() {
			finishedEffs = append(finishedEffs, e)
		}
	}
	for _, e := range finishedEffs {
		l.effects.Remove(e)
	}
}

func (l *Level) GetLevelWidth() int32 {
	return l.NumTiles.X * graphic.TILE_SIZE
}

func (l *Level) GetLevelHeight() int32 {
	return l.NumTiles.Y * graphic.TILE_SIZE
}

func (l *Level) AddEffect(e Effect) {
	l.effects.PushFront(e)
}

func (l *Level) RemoveObstacleTileObject(tid vector.TileID) {
	l.TileObjects[tid.X][tid.Y] = nil
	l.ObstMngr.RemoveTileObst(tid)
}

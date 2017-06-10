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
	BGRes            graphic.Resource
	Decorations      []Object
	TileObjects      [][]Object
	Enemies          []Enemy
	VolatileObjs     *list.List // a list of volatileObject objects
	ObstMngr         *ObstacleManager
	EnemyObstMngr    *ObstacleManager // obstacle manager for enemies
	TheHero          *Hero
	InitHeroRect     sdl.Rect
	BGColor          sdl.Color
	NumTiles         vector.Vec2D // NOTE: X, Y is TID
	ResourceRegistry map[graphic.ResourceID]graphic.Resource

	// Private
	effects *list.List
}

func ParseLevel(bgFilename string, gra *graphic.Graphic, levelArr [][]byte, decArr [][]byte) *Level {
	gra.RegisterBackgroundResource(bgFilename, graphic.RESOURCE_TYPE_BG_0, len(levelArr))
	bgRes := gra.ResourceRegistry[graphic.RESOURCE_TYPE_BG_0]

	// NOTE: index is tid.X, tid.Y
	var tileObjs [][]Object

	var enemies []Enemy

	numTiles := vector.Vec2D{int32(len(levelArr[0])), int32(len(levelArr))}
	obstMngr := NewObstacleManager(len(levelArr[0]), len(levelArr))
	enemyObstMngr := NewObstacleManager(len(levelArr[0]), len(levelArr))
	var hero *Hero

	// init tileObjs array
	for i := 0; i < int(numTiles.X); i++ {
		tileObjs = append(tileObjs, make([]Object, numTiles.Y))
	}

	addAsFullObstTile := func(tid vector.TileID, o Object) {
		tileObjs[tid.X][tid.Y] = o
		obstMngr.AddTileObst(tid)
		enemyObstMngr.AddTileObst(tid)
	}

	addAsEnemyOnlyObstTile := func(tid vector.TileID, o Object) {
		tileObjs[tid.X][tid.Y] = o
		enemyObstMngr.AddTileObst(tid)
	}

	var decorations []Object
	addDecoration := func(d *decoration) {
		decorations = append(decorations, d)
	}

	// parse level
	var currentPos vector.Pos
	for tidY := 0; tidY < int(numTiles.Y); tidY++ {
		currentPos.X = 0
		for tidX := 0; tidX < int(numTiles.X); tidX++ {
			tid := vector.TileID{int32(tidX), int32(tidY)}
			// note that levelArr's index is not TID, need reverse
			switch levelArr[tidY][tidX] {
			// Invisible block
			case '#':
				addAsFullObstTile(tid, NewInvisibleTileObject(tid))

			// Invisible block only to enemies
			case '"':
				addAsEnemyOnlyObstTile(tid, NewInvisibleTileObject(tid))

			// Brick
			case 'B':
				mainRes := gra.ResourceRegistry[graphic.RESOURCE_TYPE_BRICK]
				pieceRes := gra.ResourceRegistry[graphic.RESOURCE_TYPE_BRICK_PIECE]
				o := NewBreakableTileObject(mainRes, pieceRes, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// Ground with left grass
			case 'L':
				res := gra.ResourceRegistry[graphic.RESOURCE_TYPE_GROUD_GRASS_LEFT]
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// Ground with mid grass
			case 'G':
				res := gra.ResourceRegistry[graphic.RESOURCE_TYPE_GROUD_GRASS_MID]
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// Ground with right grass
			case 'R':
				res := gra.ResourceRegistry[graphic.RESOURCE_TYPE_GROUD_GRASS_RIGHT]
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// Inner ground in middle
			case 'I':
				res := gra.ResourceRegistry[graphic.RESOURCE_TYPE_GROUD_INNER_MID]
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// Myth box for coins
			case 'C':
				addAsFullObstTile(tid, NewCoinMythBox(currentPos, 3, gra.ResourceRegistry))

			// Myth box for mushrooms
			case 'M':
				addAsFullObstTile(tid, NewMushroomMythBox(currentPos, gra.ResourceRegistry))

			// left middle of pipe
			case '[':
				res := gra.ResourceRegistry[graphic.RESOURCE_TYPE_PIPE_LEFT_MID]
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// right middle of pipe
			case ']':
				res := gra.ResourceRegistry[graphic.RESOURCE_TYPE_PIPE_RIGHT_MID]
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// left top of pipe
			case '{':
				res := gra.ResourceRegistry[graphic.RESOURCE_TYPE_PIPE_LEFT_TOP]
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// right top of pipe
			case '}':
				res := gra.ResourceRegistry[graphic.RESOURCE_TYPE_PIPE_RIGHT_TOP]
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// left bottom of pipe
			case '<':
				res := gra.ResourceRegistry[graphic.RESOURCE_TYPE_PIPE_LEFT_BOTTOM]
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// right bottom of pipe
			case '>':
				res := gra.ResourceRegistry[graphic.RESOURCE_TYPE_PIPE_RIGHT_BOTTOM]
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// Enemy 1: mushroom enemy
			case '1':
				enemies = append(enemies, NewMushroomEnemy(currentPos, gra.ResourceRegistry))

			// Enemy 2: tortoise enemy
			case '2':
				enemies = append(enemies, NewTortoiseEnemy(currentPos, gra.ResourceRegistry))

			// Hero
			case 'H':
				if hero != nil {
					log.Fatal("more than one hero found")
				}
				hero = NewHero(currentPos, 8, 2, -16, -4, gra.ResourceRegistry)
			}
			currentPos.X += graphic.TILE_SIZE
		}
		currentPos.Y += graphic.TILE_SIZE
	}

	// parse decorations
	currentPos = vector.Pos{}
	for tidY := 0; tidY < int(numTiles.Y); tidY++ {
		currentPos.X = 0
		for tidX := 0; tidX < int(numTiles.X); tidX++ {
			tid := vector.TileID{int32(tidX), int32(tidY)}
			// note that decArr's index is not TID, need reverse
			switch decArr[tidY][tidX] {
			case '1':
				resIds := []graphic.ResourceID{
					graphic.RESOURCE_TYPE_DEC_GRASS_0,
					graphic.RESOURCE_TYPE_DEC_GRASS_1,
				}
				addDecoration(NewDecoration(tid, resIds, gra.ResourceRegistry, 800))
			}
		}
	}

	if hero == nil {
		log.Fatal("no hero found when parsing level")
	}

	return &Level{
		BGRes:            bgRes,
		Decorations:      decorations,
		TileObjects:      tileObjs,
		Enemies:          enemies,
		VolatileObjs:     list.New(),
		ObstMngr:         obstMngr,
		EnemyObstMngr:    enemyObstMngr,
		TheHero:          hero,
		InitHeroRect:     hero.levelRect,
		BGColor:          sdl.Color{204, 237, 255, 255},
		NumTiles:         numTiles,
		ResourceRegistry: gra.ResourceRegistry,
		effects:          list.New(),
	}
}

func ParseLevelFromFile(filename string, gra *graphic.Graphic) *Level {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file %s", filename)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// first line is background pic filename
	scanner.Scan()
	bgFilename := scanner.Text()

	// the second line is a "~"
	scanner.Scan()
	if scanner.Text() != "~" {
		log.Fatal("should meet \"~\"")
	}

	// parse level
	var levelArr [][]byte
	for scanner.Scan() {
		// stops when meet "~" line
		if scanner.Text() == "~" {
			break
		}

		levelArr = append(levelArr, []byte(scanner.Text()))
	}

	// parse decorations
	var decArr [][]byte
	for scanner.Scan() {
		decArr = append(decArr, []byte(scanner.Text()))
	}

	if len(levelArr) != len(decArr) {
		log.Fatal("level arr and decoration arr should have same height")
	}

	if len(levelArr[0]) != len(decArr[0]) {
		log.Fatal("level arr and decoration arr should have same width")
	}

	return ParseLevel(bgFilename, gra, levelArr, decArr)
}

func (l *Level) UpdateAndDraw(g *graphic.Graphic, camPos vector.Pos) {
	var ticks = sdl.GetTicks()

	// render background
	bgLevelRect := sdl.Rect{
		camPos.X * 80 / 100,
		0,
		l.BGRes.GetW(),
		l.BGRes.GetH(),
	}
	g.DrawResource(l.BGRes, bgLevelRect, camPos)

	// render decorations
	for _, d := range l.Decorations {
		d.Update(nil, ticks, l)
		d.Draw(g, camPos)
	}

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
	if l.TheHero.GetZIndex() == ZINDEX_0 {
		log.Fatal("hero's z-index cannot be lowest")
	}
	zIndexObjs[l.TheHero.GetZIndex()] = append(zIndexObjs[l.TheHero.GetZIndex()], l.TheHero)

	// render higher z-index one-by-one
	for _, objs := range zIndexObjs {
		if len(objs) > 0 {
			for _, o := range objs {
				o.Draw(g, camPos)
			}
		}
	}

	// update and render live enemies
	for _, e := range l.Enemies {
		if e.IsDead() {
			continue
		}

		e.Update(nil, ticks, l)
		e.Draw(g, camPos)
	}

	// update and render volatile objects
	var deadVolatileObjs []*list.Element
	for e := l.VolatileObjs.Front(); e != nil; e = e.Next() {
		vo, ok := e.Value.(volatileObject)
		if !ok {
			log.Fatalf("eff is not an volatile object: %T", e.Value)
		}
		vo.Update(nil, ticks, l)

		if vo.IsDead() {
			deadVolatileObjs = append(deadVolatileObjs, e)
		} else {
			vo.Draw(g, camPos)
		}
	}
	for _, e := range deadVolatileObjs {
		l.VolatileObjs.Remove(e)
	}

	// render effects and remove finished effects
	var finishedEffs []*list.Element
	for e := l.effects.Front(); e != nil; e = e.Next() {
		eff, ok := e.Value.(Effect)
		if !ok {
			log.Fatalf("eff is not an effect object: %T", e.Value)
		}
		eff.UpdateAndDraw(g, camPos, ticks)

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
	l.EnemyObstMngr.RemoveTileObst(tid)
}

func (l *Level) AddVolatileObject(vo volatileObject) {
	l.VolatileObjs.PushBack(vo)
}

func (l *Level) AddEnemy(e Enemy) {
	l.Enemies = append(l.Enemies, e)
}

func (l *Level) Restart() {
	l.TheHero.levelRect = l.InitHeroRect
}

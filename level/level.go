package level

import (
	"bufio"
	"log"
	"os"

	"container/list"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/event"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

type Level struct {
	// Public
	BGRes         graphic.Resource
	Decorations   []Object
	TileObjects   [][]Object
	NumTiles      vector.Vec2D // NOTE: X, Y is TID
	Enemies       []Enemy
	VolatileObjs  *list.List // a list of volatileObject objects
	ObstMngr      *ObstacleManager
	EnemyObstMngr *ObstacleManager // obstacle manager for enemies
	TheHero       *Hero
	InitHeroRect  sdl.Rect
	BGColor       sdl.Color

	// Private
	effects *list.List

	// ticks when hero died, used to wait for hero die effect to finish
	lastHeroDieTicks uint32

	// if should restart
	shouldRestart bool
}

func ParseLevel(bgFilename string, levelArr [][]byte, decArr [][]byte) *Level {
	graphic.RegisterBackgroundResource(bgFilename, graphic.RESOURCE_TYPE_BG_0, len(levelArr))
	bgRes := graphic.Res(graphic.RESOURCE_TYPE_BG_0)

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
				mainRes := graphic.Res(graphic.RESOURCE_TYPE_BRICK)
				pieceRes := graphic.Res(graphic.RESOURCE_TYPE_BRICK_PIECE)
				o := NewBreakableTileObject(mainRes, pieceRes, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// Ground with left grass
			case 'L':
				res := graphic.Res(graphic.RESOURCE_TYPE_GROUD_GRASS_LEFT)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// Ground with mid grass
			case 'G':
				res := graphic.Res(graphic.RESOURCE_TYPE_GROUD_GRASS_MID)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// Ground with right grass
			case 'R':
				res := graphic.Res(graphic.RESOURCE_TYPE_GROUD_GRASS_RIGHT)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// Inner ground in middle
			case 'I':
				res := graphic.Res(graphic.RESOURCE_TYPE_GROUD_INNER_MID)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// Myth box for coins
			case 'C':
				addAsFullObstTile(tid, NewCoinMythBox(currentPos, 3))

			// Myth box for mushrooms
			case 'M':
				addAsFullObstTile(tid, NewMushroomMythBox(currentPos))

			// left middle of pipe
			case '[':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_MID)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// right middle of pipe
			case ']':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_MID)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// left top of pipe
			case '{':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_TOP)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// right top of pipe
			case '}':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_TOP)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// left bottom of pipe
			case '<':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_BOTTOM)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// right bottom of pipe
			case '>':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_BOTTOM)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

			// Enemy 1: mushroom enemy
			case '1':
				enemies = append(enemies, NewMushroomEnemy(currentPos))

			// Enemy 2: tortoise enemy
			case '2':
				enemies = append(enemies, NewTortoiseEnemy(currentPos))

			// Hero
			case 'H':
				if hero != nil {
					log.Fatal("more than one hero found")
				}
				hero = NewHero(currentPos, 0.2, 0.1)
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
				addDecoration(NewDecoration(tid, resIds, 800))
			}
		}
	}

	if hero == nil {
		log.Fatal("no hero found when parsing level")
	}

	return &Level{
		BGRes:         bgRes,
		Decorations:   decorations,
		TileObjects:   tileObjs,
		Enemies:       enemies,
		VolatileObjs:  list.New(),
		ObstMngr:      obstMngr,
		EnemyObstMngr: enemyObstMngr,
		TheHero:       hero,
		InitHeroRect:  hero.levelRect,
		BGColor:       sdl.Color{204, 237, 255, 255},
		NumTiles:      numTiles,
		effects:       list.New(),
	}
}

func ParseLevelFromFile(filename string) *Level {
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

	return ParseLevel(bgFilename, levelArr, decArr)
}

func (l *Level) Init() {
	l.fadeIn()
}

func (l *Level) HandleEvents(events *intsets.Sparse) {
	if events.Has(int(event.EVENT_KEYDOWN_F4)) {
		l.fadeIn()
	}
}

func (l *Level) Update(events *intsets.Sparse, ticks uint32) {
	if l.shouldRestart {
		l.shouldRestart = false
		l.Restart()
	}

	// wait for a while after hero died
	if l.lastHeroDieTicks > 0 && ticks-l.lastHeroDieTicks > 1500 {
		l.shouldRestart = true
		l.lastHeroDieTicks = 0
	}

	// update tile objects
	for i := 0; i < int(l.NumTiles.X); i++ {
		for j := 0; j < int(l.NumTiles.Y); j++ {
			o := l.TileObjects[i][j]
			if o == nil {
				continue
			}
			o.Update(ticks, l)
		}
	}

	if l.TheHero.diedTicks > 0 {
		// if hero just died, show hero die effects
		var deadRes graphic.Resource
		if l.TheHero.isFacingRight {
			deadRes = l.TheHero.currResStandRight
		} else {
			deadRes = l.TheHero.currResStandLeft
		}
		l.AddEffect(NewStraightDeadDownEffect(deadRes, l.TheHero.getRenderRect(), l.TheHero.diedTicks))

		// and reset hero's diedTicks so that the effect will only be added once
		l.TheHero.diedTicks = 0

		// set lastHeroDieTicks, so we can know we need to wait for a while
		l.lastHeroDieTicks = ticks
	} else if !l.TheHero.IsDead() {
		// update hero with events
		l.TheHero.HandleEvents(events, l)
		l.TheHero.Update(ticks, l)

		// if hero is out of level, kills it
		if l.isOutOfLevel(l.TheHero.GetRect()) {
			l.TheHero.Kill()
		}
	}

	// update live enemies
	for _, e := range l.Enemies {
		if e.IsDead() {
			continue
		}

		// if enemy is out of level, kills it
		if l.isOutOfLevel(e.GetRect()) {
			e.Kill()
		}

		e.Update(ticks, l)
	}

	// update volatile objects
	var deadVolatileObjs []*list.Element
	for e := l.VolatileObjs.Front(); e != nil; e = e.Next() {
		vo, ok := e.Value.(volatileObject)
		if !ok {
			log.Fatalf("eff is not an volatile object: %T", e.Value)
		}
		vo.Update(ticks, l)

		if vo.IsDead() {
			deadVolatileObjs = append(deadVolatileObjs, e)
		}
	}
	for _, e := range deadVolatileObjs {
		l.VolatileObjs.Remove(e)
	}

	// update effects and remove finished effects
	var finishedEffs []*list.Element
	for e := l.effects.Front(); e != nil; e = e.Next() {
		eff, ok := e.Value.(Effect)
		if !ok {
			log.Fatalf("eff is not an effect object: %T", e.Value)
		}
		eff.Update(ticks)

		if eff.Finished() {
			finishedEffs = append(finishedEffs, e)
		}
	}
	for _, e := range finishedEffs {
		l.effects.Remove(e)
	}
}

func (l *Level) Draw(camPos vector.Pos, ticks uint32) {
	// render background
	bgLevelRect := sdl.Rect{
		camPos.X * 80 / 100,
		0,
		l.BGRes.GetW(),
		l.BGRes.GetH(),
	}
	graphic.DrawResource(l.BGRes, bgLevelRect, camPos)

	// render decorations
	for _, d := range l.Decorations {
		d.Update(ticks, l)
		d.Draw(camPos)
	}

	// put all tile objects in render list
	var zIndexObjs [ZINDEX_NUM][]Object
	for i := 0; i < int(l.NumTiles.X); i++ {
		for j := 0; j < int(l.NumTiles.Y); j++ {
			o := l.TileObjects[i][j]
			if o == nil {
				continue
			}

			z := o.GetZIndex()
			zIndexObjs[z] = append(zIndexObjs[z], o)
		}
	}

	// put non-tile objects into render queue
	if l.TheHero.GetZIndex() == ZINDEX_0 {
		log.Fatal("hero's z-index cannot be lowest")
	}
	zIndexObjs[l.TheHero.GetZIndex()] = append(zIndexObjs[l.TheHero.GetZIndex()], l.TheHero)

	// render live enemies
	for _, e := range l.Enemies {
		if e.IsDead() {
			continue
		}

		z := e.GetZIndex()
		zIndexObjs[z] = append(zIndexObjs[z], e)
	}

	// render higher z-index one-by-one
	for _, objs := range zIndexObjs {
		if len(objs) > 0 {
			for _, o := range objs {
				o.Draw(camPos)
			}
		}
	}

	// render volatile objects (don't care z-index)
	for e := l.VolatileObjs.Front(); e != nil; e = e.Next() {
		vo, _ := e.Value.(volatileObject)

		if !vo.IsDead() {
			vo.Draw(camPos)
		}
	}

	// render effects (don't care z-index)
	for e := l.effects.Front(); e != nil; e = e.Next() {
		eff, _ := e.Value.(Effect)

		if !eff.Finished() {
			eff.Draw(camPos, ticks)
		}
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
	l.fadeIn()
	l.TheHero.Reborn(l.InitHeroRect)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Private helpers
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (l *Level) fadeIn() {
	l.AddEffect(NewScreenFadeEffect(true, 1000, sdl.GetTicks()))
}

func (l *Level) isOutOfLevel(rect sdl.Rect) bool {
	levelWidth := l.NumTiles.X * graphic.TILE_SIZE
	levelHeight := l.NumTiles.Y * graphic.TILE_SIZE
	if rect.X > levelWidth || rect.X+rect.W < 0 || rect.Y > levelHeight || rect.Y+rect.H < 0 {
		return true
	} else {
		return false
	}
}

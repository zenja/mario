package level

import (
	"log"

	"container/list"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/event"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

type Level struct {
	// Public
	Spec          *LevelSpec
	BGRes         graphic.Resource
	BGColor       sdl.Color
	Decorations   []Object
	TileObjects   [][]Object
	NumTiles      vector.Vec2D // NOTE: X, Y is TID
	Enemies       []Enemy
	VolatileObjs  *list.List // a list of volatileObject objects
	ObstMngr      *ObstacleManager
	EnemyObstMngr *ObstacleManager // obstacle manager for enemies
	TheHero       *Hero
	InitHeroRect  sdl.Rect
	Coins         int

	// Private

	effects *list.List

	// if not empty, it means we should switch to next level
	nextLevelName string
}

func (l *Level) Init() {
	l.fadeIn()
	l.TheHero.Reborn(l.InitHeroRect)
}

func (l *Level) HandleEvents(events *intsets.Sparse) {
	if events.Has(int(event.EVENT_KEYDOWN_F4)) {
		l.fadeIn()
	}
}

func (l *Level) Update(events *intsets.Sparse, ticks uint32) {
	// defensive prevention
	if nextLevel, shouldSwitch := l.GetNextLevel(); shouldSwitch {
		log.Fatalf("level should switch to %s, cannot update", nextLevel)
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

	if !l.TheHero.IsDead() {
		// update hero with events
		l.TheHero.HandleEvents(events, l)
		l.TheHero.Update(ticks, l)

		// if hero is out of level, kills it
		if l.isOutOfLevel(l.TheHero.GetRect()) {
			l.TheHero.Kill(l)
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

	// update effects, run on-finished hook and remove finished effects
	var finishedEffs []*list.Element
	for e := l.effects.Front(); e != nil; e = e.Next() {
		eff, ok := e.Value.(Effect)
		if !ok {
			log.Fatalf("eff is not an effect object: %T", e.Value)
		}
		eff.Update(ticks)

		if eff.Finished() {
			eff.OnFinished()
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

	// put hero into render queue
	if l.TheHero.GetZIndex() == ZINDEX_0 {
		log.Fatal("hero's z-index cannot be lowest")
	}
	zIndexObjs[l.TheHero.GetZIndex()] = append(zIndexObjs[l.TheHero.GetZIndex()], l.TheHero)

	// put live enemies into render queue
	for _, e := range l.Enemies {
		if e.IsDead() {
			continue
		}

		z := e.GetZIndex()
		zIndexObjs[z] = append(zIndexObjs[z], e)
	}

	// render z-index layers one-by-one
	for _, objs := range zIndexObjs {
		for _, o := range objs {
			o.Draw(camPos)
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
	// reset things needs to be reset with new level
	newLevel := BuildLevel(l.Spec)
	l.TileObjects = newLevel.TileObjects
	l.Enemies = newLevel.Enemies
	l.ObstMngr = newLevel.ObstMngr
	l.EnemyObstMngr = newLevel.EnemyObstMngr

	l.Init()
}

func (l *Level) ShouldSwitchLevel(nextLevelName string) {
	l.nextLevelName = nextLevelName
}

// a indicator to upper game that we should switch level
// return (next level name, should switch level)
func (l *Level) GetNextLevel() (string, bool) {
	if len(l.nextLevelName) == 0 {
		return "", false
	}
	return l.nextLevelName, true
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

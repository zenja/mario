package level

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Basic utils functions
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func enemySimpleMove(
	ticks uint32,
	lastTicks uint32,
	vel *vector.Vec2D,
	levelRect *sdl.Rect,
	level *Level) {

	enemySimpleMoveEx(ticks, lastTicks, vel, levelRect, level, nil, nil)
}

func enemySimpleMoveEx(
	ticks uint32,
	lastTicks uint32,
	vel *vector.Vec2D,
	levelRect *sdl.Rect,
	level *Level,
	onHitLeft func(),
	onHitRight func()) {

	gravity := vector.Vec2D{0, 50}
	vel.Add(gravity)

	maxVel := vector.Vec2D{int32(graphic.TILE_SIZE * 30 / 100), int32(graphic.TILE_SIZE * 30 / 100)}
	velocityStep := CalcVelocityStep(*vel, ticks, lastTicks, &maxVel)
	levelRect.X += velocityStep.X
	levelRect.Y += velocityStep.Y

	_, hitRight, hitBottom, hitLeft, _ := level.ObstMngr.SolveCollision(levelRect, SOLVE_COLLISION_ENEMY)

	if hitRight {
		vel.X = -vel.X
		if onHitRight != nil {
			onHitRight()
		}
	}
	if hitLeft {
		vel.X = -vel.X
		if onHitLeft != nil {
			onHitLeft()
		}
	}

	// prevent too big down velocity
	if velocityStep.Y > 0 && hitBottom {
		vel.Y = 0
	}
}

func hitEnemiesOnTop(selfRect *sdl.Rect, level *Level, ticks uint32) {
	for _, e := range level.Enemies {
		emyRectLower := sdl.Rect{
			X: e.GetRect().X,
			Y: e.GetRect().Y + 1,
			W: e.GetRect().W,
			H: e.GetRect().H,
		}
		if emyRectLower.HasIntersection(selfRect) {
			if !e.IsDead() {
				e.hitByBottomTile(level, ticks)
			}
		}
	}
}

func hurtHeroIfIntersectEnough(hero *Hero, emy Enemy, level *Level) {
	heroRect := hero.GetRect()
	emyRect := emy.GetRect()
	interRect, hasIntersection := heroRect.Intersect(&emyRect)
	if !hasIntersection {
		return
	}

	if interRect.W > int32(float64(emyRect.W)*0.3) && interRect.H > int32(float64(emyRect.H)*0.3) {
		hero.Hurt(level)
	}
}

func hurtHeroIfIntersectEnoughEx(hero *Hero, emy Enemy, level *Level, ratio float64) {
	heroRect := hero.GetRect()
	emyRect := emy.GetRect()
	interRect, hasIntersection := heroRect.Intersect(&emyRect)
	if !hasIntersection {
		return
	}

	if interRect.W > int32(float64(emyRect.W)*ratio) && interRect.H > int32(float64(emyRect.H)*ratio) {
		hero.Hurt(level)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Advanced helper functions
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func changeDirectionRandomly(randomness int, isFacingRight *bool, velX *int32) {
	if rand.Intn(randomness) == 7 {
		*isFacingRight = !(*isFacingRight)
		*velX = -(*velX)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// canSay
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type canSay struct {
	sentences    []string
	lastSayTicks uint32
}

func newCanSay(sentences []string) canSay {
	return canSay{sentences: sentences}
}

func (cs *canSay) say(ticks uint32, level *Level, minRGB, maxRGB int, getSentencePosFunc func() vector.Pos) {
	if len(cs.sentences) == 0 {
		return
	}

	randColor := sdl.Color{
		uint8(rand.Intn(maxRGB-minRGB) + minRGB),
		uint8(rand.Intn(maxRGB-minRGB) + minRGB),
		uint8(rand.Intn(maxRGB-minRGB) + minRGB),
		255,
	}
	randSentence := cs.sentences[rand.Intn(len(cs.sentences))]
	if ticks-cs.lastSayTicks > 3000 {
		level.AddEffect(NewShowTextEffect(randSentence, randColor, getSentencePosFunc, ticks, 2000))
		cs.lastSayTicks = ticks
	}
}

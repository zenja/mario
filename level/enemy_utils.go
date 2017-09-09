package level

import (
	"math"
	"math/rand"

	"log"

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

type bulletEnemyType int

const (
	BULLET_ENEMY_FIREBALL bulletEnemyType = iota
	BULLET_ENEMY_SWORD
	BULLET_ENEMY_APPLE
	BULLET_ENEMY_CHERRY
	BULLET_ENEMY_MOON
	BULLET_ENEMY_AXE
	BULLET_ENEMY_SKULL
)

func fireToHeroRandomly(randomness int, level *Level, self Enemy, t bulletEnemyType, finalVel float64) {
	if rand.Intn(randomness) != 7 {
		return
	}

	pos := vector.Pos{
		X: self.GetRect().X + self.GetRect().W/2,
		Y: self.GetRect().Y + self.GetRect().H/2,
	}

	vel := calcBulletEnemyVel(level.TheHero.levelRect, pos, finalVel)
	var e *bulletEnemy

	switch t {
	case BULLET_ENEMY_FIREBALL:
		e = NewFireBallEnemy(pos, vel)
	case BULLET_ENEMY_SWORD:
		e = NewSwordEnemy(pos, vel)
	case BULLET_ENEMY_APPLE:
		e = NewAppleEnemy(pos, vel)
	case BULLET_ENEMY_CHERRY:
		e = NewCherryEnemy(pos, vel)
	case BULLET_ENEMY_MOON:
		e = NewMoonEnemy(pos, vel)
	case BULLET_ENEMY_AXE:
		e = NewAxeEnemy(pos, vel)
	case BULLET_ENEMY_SKULL:
		e = NewSkullEnemy(pos, vel)
	default:
		log.Fatalf("not supported bullet enemy type: %d", t)
	}

	level.AddEnemy(e)
}

func fireAroundRandomly(randomness int, level *Level, self Enemy, t bulletEnemyType) {
	if rand.Intn(randomness) != 7 {
		return
	}

	pos := vector.Pos{
		X: self.GetRect().X + self.GetRect().W/2,
		Y: self.GetRect().Y + self.GetRect().H/2,
	}

	type bulletEnemyMaker func(startPos vector.Pos, initVel vector.Vec2D) *bulletEnemy
	var maker bulletEnemyMaker

	switch t {
	case BULLET_ENEMY_FIREBALL:
		maker = NewFireBallEnemy
	case BULLET_ENEMY_SWORD:
		maker = NewSwordEnemy
	case BULLET_ENEMY_APPLE:
		maker = NewAppleEnemy
	case BULLET_ENEMY_CHERRY:
		maker = NewCherryEnemy
	case BULLET_ENEMY_MOON:
		maker = NewMoonEnemy
	case BULLET_ENEMY_AXE:
		maker = NewAxeEnemy
	case BULLET_ENEMY_SKULL:
		maker = NewSkullEnemy
	default:
		log.Fatalf("not supported bullet enemy type: %d", t)
	}

	vels := []vector.Vec2D{
		{-200, 0},
		{200, 0},
		{0, -200},
		{-130, -130},
		{130, -130},
	}
	for _, v := range vels {
		e := maker(pos, v)
		level.AddEnemy(e)
	}
}

func calcBulletEnemyVel(heroRect sdl.Rect, bulletPos vector.Pos, finalVel float64) vector.Vec2D {
	heroX := heroRect.X
	heroY := heroRect.Y
	deltaX := heroX - bulletPos.X
	deltaY := heroY - bulletPos.Y
	var scale float64 = finalVel / math.Sqrt(math.Pow(float64(deltaX), 2)+math.Pow(float64(deltaY), 2))
	return vector.Vec2D{
		X: int32(float64(deltaX) * scale),
		Y: int32(float64(deltaY) * scale),
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

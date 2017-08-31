package level

import (
	"log"

	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/audio"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type hittableByHero interface {
	hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32)
}

type hittableByBottomTile interface {
	hitByBottomTile(level *Level, ticks uint32)
}

type hittableByBullet interface {
	hitByBullet(b bullet, level *Level, ticks uint32)
}

type Enemy interface {
	// Enemy is an object
	Object

	// Enemy is hittable by hero
	hittableByHero

	// Enemy is hittable by breaking tile (from enemy's bottom tile)
	hittableByBottomTile

	// Enemy is hittable by bullet
	hittableByBullet

	// if the enemy is dead, if so, don't need to update/draw
	IsDead() bool

	// after Kill(), IsDead() should return true
	Kill()
}

type basicEnemy struct {
	isDead bool
}

func (be *basicEnemy) IsDead() bool {
	return be.isDead
}

func (be *basicEnemy) Kill() {
	be.isDead = true
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// mushroomEnemy
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type mushroomEnemy struct {
	basicEnemy

	res0      graphic.Resource
	res1      graphic.Resource
	resHit    graphic.Resource
	resDown   graphic.Resource
	currRes   graphic.Resource
	levelRect sdl.Rect
	lastTicks uint32
	velocity  vector.Vec2D
}

func NewMushroomEnemy(startPos vector.Pos) *mushroomEnemy {
	res0 := graphic.Res(graphic.RESOURCE_TYPE_MUSHROOM_ENEMY_0)
	return &mushroomEnemy{
		res0:      res0,
		res1:      graphic.Res(graphic.RESOURCE_TYPE_MUSHROOM_ENEMY_1),
		resHit:    graphic.Res(graphic.RESOURCE_TYPE_MUSHROOM_ENEMY_HIT),
		resDown:   graphic.Res(graphic.RESOURCE_TYPE_MUSHROOM_ENEMY_DOWN),
		currRes:   res0,
		levelRect: sdl.Rect{startPos.X, startPos.Y, res0.GetW(), res0.GetH()},
		velocity:  vector.Vec2D{100, 0},
	}
}

func (m *mushroomEnemy) GetRect() sdl.Rect {
	return m.levelRect
}

func (m *mushroomEnemy) GetZIndex() int {
	return ZINDEX_4
}

func (m *mushroomEnemy) Update(ticks uint32, level *Level) {
	if m.lastTicks == 0 {
		m.lastTicks = ticks
		return
	}

	enemySimpleMove(ticks, m.lastTicks, &m.velocity, &m.levelRect, level)

	m.updateResource(ticks)

	m.lastTicks = ticks
}

func (m *mushroomEnemy) Draw(camPos vector.Pos) {
	graphic.DrawResource(m.currRes, m.levelRect, camPos)
}

func (m *mushroomEnemy) updateResource(ticks uint32) {
	if ticks%1000 < 500 {
		m.currRes = m.res0
	} else {
		m.currRes = m.res1
	}
}

func (m *mushroomEnemy) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	if direction == HIT_FROM_TOP_W_INTENT {
		// dead immediately!
		m.isDead = true

		// bounce the hero up
		h.velocity.Y = -1200

		// add dead effect
		level.AddEffect(NewShowOnceEffect(m.resHit, GetRectStartPos(m.levelRect), ticks, 500))

		audio.PlaySound(audio.SOUND_STOMP)
	} else {
		// hero is hurt
		hurtHeroIfIntersectEnough(h, m, level)
	}
}

func (m *mushroomEnemy) hitByBottomTile(level *Level, ticks uint32) {
	m.dieDown(true, level, ticks)
}

func (m *mushroomEnemy) hitByBullet(b bullet, level *Level, ticks uint32) {
	if b.GetDamage() == 1 {
		return
	}

	var dieToRight bool
	if b.GetRect().X < m.levelRect.X {
		dieToRight = true
	}
	m.dieDown(dieToRight, level, ticks)
}

func (m *mushroomEnemy) dieDown(toRight bool, level *Level, ticks uint32) {
	m.isDead = true
	level.AddEffect(NewDeadDownEffect(m.resDown, toRight, m.levelRect, ticks))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// tortoiseEnemy
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const (
	tortoiseBumpingVelocityXRight = 800
)

type tortoiseEnemy struct {
	basicEnemy

	resLeft0      graphic.Resource
	resLeft1      graphic.Resource
	resRight0     graphic.Resource
	resRight1     graphic.Resource
	resSemiInside graphic.Resource
	resInside     graphic.Resource
	currRes       graphic.Resource

	isFacingRight bool
	levelRect     sdl.Rect
	velocity      vector.Vec2D
	lastTicks     uint32

	insideStartTicks uint32 // when tortoise go inside
	bumpStartTicks   uint32 // when tortoise start bumping
}

func NewTortoiseEnemy(startPos vector.Pos) Enemy {
	resLeft0 := graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RED_LEFT_0)
	return &tortoiseEnemy{
		resLeft0:      resLeft0,
		resLeft1:      graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RED_LEFT_1),
		resRight0:     graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RED_RIGHT_0),
		resRight1:     graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RED_RIGHT_1),
		resSemiInside: graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RED_SEMI_INSIDE),
		resInside:     graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RED_INSIDE),
		currRes:       resLeft0,
		levelRect:     sdl.Rect{startPos.X, startPos.Y, resLeft0.GetW(), resLeft0.GetH()},
		velocity:      vector.Vec2D{-100, 0},
	}
}

func NewTortoiseEnemyRandomSpeed(startPos vector.Pos, faceRight bool, maxSpeedUp int) Enemy {
	resLeft0 := graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RED_LEFT_0)
	velX := 100 + rand.Intn(maxSpeedUp)
	if !faceRight {
		velX = -velX
	}
	return &tortoiseEnemy{
		resLeft0:      resLeft0,
		resLeft1:      graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RED_LEFT_1),
		resRight0:     graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RED_RIGHT_0),
		resRight1:     graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RED_RIGHT_1),
		resSemiInside: graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RED_SEMI_INSIDE),
		resInside:     graphic.Res(graphic.RESOURCE_TYPE_TORTOISE_RED_INSIDE),
		currRes:       resLeft0,
		levelRect:     sdl.Rect{startPos.X, startPos.Y, resLeft0.GetW(), resLeft0.GetH()},
		velocity:      vector.Vec2D{int32(velX), 0},
		isFacingRight: faceRight,
	}
}

func (t *tortoiseEnemy) GetRect() sdl.Rect {
	return t.levelRect
}

func (t *tortoiseEnemy) GetZIndex() int {
	return ZINDEX_4
}

func (t *tortoiseEnemy) Update(ticks uint32, level *Level) {
	if t.lastTicks == 0 {
		t.lastTicks = ticks
		return
	}

	onHitLeft := func() {
		t.isFacingRight = true
		if t.bumpStartTicks > 0 {
			level.AddEffect(t.newBangEffect(true))
		}
	}
	onHitRight := func() {
		t.isFacingRight = false
		if t.bumpStartTicks > 0 {
			level.AddEffect(t.newBangEffect(false))
		}
	}
	enemySimpleMoveEx(ticks, t.lastTicks, &t.velocity, &t.levelRect, level, onHitLeft, onHitRight)

	t.updateResource(ticks)

	t.lastTicks = ticks
}

func (t *tortoiseEnemy) Draw(camPos vector.Pos) {
	graphic.DrawResource(t.currRes, t.levelRect, camPos)
}

func (t *tortoiseEnemy) updateResource(ticks uint32) {
	if t.insideStartTicks > 0 || t.bumpStartTicks > 0 {
		t.switchResourceAndAdjustRect(t.resInside)
		return
	}

	if ticks%1000 < 500 {
		if t.isFacingRight {
			t.switchResourceAndAdjustRect(t.resRight0)
		} else {
			t.switchResourceAndAdjustRect(t.resLeft0)
		}
	} else {
		if t.isFacingRight {
			t.switchResourceAndAdjustRect(t.resRight1)
		} else {
			t.switchResourceAndAdjustRect(t.resLeft1)
		}
	}
}

func (t *tortoiseEnemy) switchResourceAndAdjustRect(newRes graphic.Resource) {
	oldRes := t.currRes
	t.currRes = newRes
	// make sure new res stands on same bottom position
	incH := newRes.GetH() - oldRes.GetH()
	t.levelRect.Y -= incH
	t.levelRect.W = t.currRes.GetW()
	t.levelRect.H = t.currRes.GetH()
}

func (t *tortoiseEnemy) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	if t.insideStartTicks > 0 && t.bumpStartTicks > 0 {
		log.Fatal("bug! insideStartTicks and bumpStartTicks cannot be positive at the same time!")
	}

	switch direction {
	case HIT_FROM_TOP_W_INTENT:
		// bounce the hero up
		h.velocity.Y = -1200

		switch {
		// case 1: normal state => go inside, don't move in X
		case t.insideStartTicks == 0 && t.bumpStartTicks == 0:
			t.toInsideState(ticks)

		// case 2: inside state => start bumping
		case t.bumpStartTicks == 0:
			// decide move right or left
			heroMidX := h.levelRect.X + h.levelRect.W/2
			tortoiseMidX := t.levelRect.X + t.levelRect.W/2
			if heroMidX < tortoiseMidX {
				t.toBumpingState(ticks, true)
			} else {
				t.toBumpingState(ticks, false)
			}

		// case 3: bumping state => stop bumping, turn to inside state, don't move in X
		default:
			t.toInsideState(ticks)
		}

		audio.PlaySound(audio.SOUND_STOMP)

	case HIT_FROM_LEFT_W_INTENT:
		if t.insideStartTicks > 0 {
			t.toBumpingState(ticks, true)
			audio.PlaySound(audio.SOUND_KICK)
		} else {
			hurtHeroIfIntersectEnough(h, t, level)
		}

	case HIT_FROM_RIGHT_W_INTENT:
		if t.insideStartTicks > 0 {
			t.toBumpingState(ticks, false)
			audio.PlaySound(audio.SOUND_KICK)
		} else {
			hurtHeroIfIntersectEnough(h, t, level)
		}

	default:
		// hero is hurt
		hurtHeroIfIntersectEnough(h, t, level)
	}
}

func (t *tortoiseEnemy) hitByBottomTile(level *Level, ticks uint32) {
	t.dieDown(true, level, ticks)
}

func (t *tortoiseEnemy) hitByBullet(b bullet, level *Level, ticks uint32) {
	if b.GetDamage() == 1 {
		return
	}

	var dieToRight bool
	if b.GetRect().X < t.levelRect.X {
		dieToRight = true
	}
	t.dieDown(dieToRight, level, ticks)
}

func (t *tortoiseEnemy) dieDown(toRight bool, level *Level, ticks uint32) {
	t.isDead = true
	level.AddEffect(NewDeadDownEffect(t.currRes, toRight, t.levelRect, ticks))
}

func (t *tortoiseEnemy) toInsideState(ticks uint32) {
	// change state
	t.insideStartTicks = ticks
	t.bumpStartTicks = 0

	t.velocity.X = 0
}

func (t *tortoiseEnemy) toBumpingState(ticks uint32, toRight bool) {
	// change state
	t.insideStartTicks = 0
	t.bumpStartTicks = ticks

	if toRight {
		// move right
		t.velocity.X = tortoiseBumpingVelocityXRight
	} else {
		// move left
		t.velocity.X = -tortoiseBumpingVelocityXRight
	}
}

func (t *tortoiseEnemy) newBangEffect(hitLeft bool) *showOnceEffect {
	var xDelta int32
	if hitLeft {
		xDelta = -20
	} else {
		xDelta = 20
	}
	bangStartPos := vector.Vec2D{
		t.levelRect.X + xDelta,
		t.levelRect.Y,
	}
	return NewShowOnceEffect(graphic.Res(graphic.RESOURCE_TYPE_BANG), bangStartPos, sdl.GetTicks(), 50)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// eaterFlower
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var _ Enemy = &eaterFlower{}

type eaterFlower struct {
	basicEnemy
	*animationTileObj

	maxY      int32
	minY      int32
	goingUp   bool
	lastTicks uint32
}

func NewEaterFlower(tid vector.TileID) *eaterFlower {
	res := graphic.Res(graphic.RESOURCE_TYPE_EATER_FLOWER_0)
	reses := []graphic.ResourceID{
		graphic.RESOURCE_TYPE_EATER_FLOWER_0,
		graphic.RESOURCE_TYPE_EATER_FLOWER_1,
	}
	tidRect := GetTileRect(tid)
	startX := tidRect.X + (graphic.TILE_SIZE*2-res.GetW())/2
	startY := tidRect.Y - graphic.TILE_SIZE
	return &eaterFlower{
		animationTileObj: NewAnimationObject(vector.Pos{startX, startY}, reses, 200, ZINDEX_3),
		maxY:             startY,
		minY:             startY - graphic.TILE_SIZE - res.GetH(),
		goingUp:          true,
	}
}

func (ef *eaterFlower) GetRect() sdl.Rect {
	return ef.levelRect
}

func (ef *eaterFlower) GetZIndex() int {
	return ZINDEX_1
}

func (ef *eaterFlower) Update(ticks uint32, level *Level) {
	if ef.lastTicks == 0 {
		ef.lastTicks = ticks
		return
	}

	ef.animationTileObj.Update(ticks, level)

	if ef.levelRect.Y >= ef.maxY {
		ef.levelRect.Y = ef.maxY
		ef.goingUp = true
	} else if ef.levelRect.Y < ef.minY {
		ef.levelRect.Y = ef.minY
		ef.goingUp = false
	}

	var velocity vector.Vec2D
	if ef.goingUp {
		velocity.Y = -100
	} else {
		velocity.Y = 100
	}
	step := CalcVelocityStep(velocity, ticks, ef.lastTicks, nil)
	ef.levelRect.Y += step.Y

	ef.lastTicks = ticks
}

func (ef *eaterFlower) Draw(camPos vector.Pos) {
	ef.animationTileObj.Draw(camPos)
}

func (ef *eaterFlower) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	hurtHeroIfIntersectEnough(h, ef, level)
}

func (ef *eaterFlower) hitByBottomTile(level *Level, ticks uint32) {
	// Do Nothing
}

func (ef *eaterFlower) hitByBullet(b bullet, level *Level, ticks uint32) {
	if b.GetDamage() == 1 {
		return
	}

	ef.isDead = true
	bangRes := graphic.Res(graphic.RESOURCE_TYPE_BANG)
	bangStartPos := vector.Vec2D{
		ef.levelRect.X,
		ef.levelRect.Y,
	}
	level.AddEffect(NewShowOnceEffect(bangRes, bangStartPos, sdl.GetTicks(), 50))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// goodMushroom
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type goodMushroom struct {
	basicEnemy

	res       graphic.Resource
	levelRect sdl.Rect
	lastTicks uint32
	velocity  vector.Vec2D
}

func NewGoodMushroom(startPos vector.Pos) *goodMushroom {
	res := graphic.Res(graphic.RESOURCE_TYPE_GOOD_MUSHROOM)
	return &goodMushroom{
		res:       res,
		levelRect: sdl.Rect{startPos.X, startPos.Y, res.GetW(), res.GetH()},
		velocity:  vector.Vec2D{-100, -500},
	}
}

func (gm *goodMushroom) GetRect() sdl.Rect {
	return gm.levelRect
}

func (gm *goodMushroom) GetZIndex() int {
	return ZINDEX_1
}

func (gm *goodMushroom) Update(ticks uint32, level *Level) {
	if gm.lastTicks == 0 {
		gm.lastTicks = ticks
		return
	}

	enemySimpleMove(ticks, gm.lastTicks, &gm.velocity, &gm.levelRect, level)

	gm.lastTicks = ticks
}

func (gm *goodMushroom) Draw(camPos vector.Pos) {
	graphic.DrawResource(gm.res, gm.levelRect, camPos)
}

func (gm *goodMushroom) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	gm.isDead = true
	// upgrade hero
	h.upgrade(level)
}

func (gm *goodMushroom) hitByBottomTile(level *Level, ticks uint32) {
	// bounce up
	gm.velocity.Y -= 500
}

func (gm *goodMushroom) hitByBullet(b bullet, level *Level, ticks uint32) {
	// No interaction with bullet; Do nothing
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// upgradeFlower
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var _ Enemy = &upgradeFlower{}

type upgradeFlower struct {
	basicEnemy

	res       graphic.Resource
	levelRect sdl.Rect
	lastTicks uint32
	velocity  vector.Vec2D
}

func NewUpgradeFlower(startPos vector.Pos) *upgradeFlower {
	res := graphic.Res(graphic.RESOURCE_TYPE_UPGRADE_FLOWER)
	return &upgradeFlower{
		res:       res,
		levelRect: sdl.Rect{startPos.X, startPos.Y, res.GetW(), res.GetH()},
		velocity:  vector.Vec2D{0, -500},
	}
}

func (uf *upgradeFlower) GetRect() sdl.Rect {
	return uf.levelRect
}

func (uf *upgradeFlower) GetZIndex() int {
	return ZINDEX_1
}

func (uf *upgradeFlower) Update(ticks uint32, level *Level) {
	if uf.lastTicks == 0 {
		uf.lastTicks = ticks
		return
	}

	enemySimpleMove(ticks, uf.lastTicks, &uf.velocity, &uf.levelRect, level)

	uf.lastTicks = ticks
}

func (uf *upgradeFlower) Draw(camPos vector.Pos) {
	graphic.DrawResource(uf.res, uf.levelRect, camPos)
}

func (uf *upgradeFlower) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	uf.isDead = true
	// upgrade hero to highest
	h.upgradeToHighest(level)
}

func (uf *upgradeFlower) hitByBottomTile(level *Level, ticks uint32) {
	// bounce up
	uf.velocity.Y -= 500
}

func (uf *upgradeFlower) hitByBullet(b bullet, level *Level, ticks uint32) {
	// No interaction with bullet; Do nothing
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// levelJumper
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var _ Enemy = &levelJumper{}

// levelJumper is a two tile wide invisible "enemy" used to go to another level
// it is usually placed on a pipe's top
type levelJumper struct {
	basicEnemy

	nextLevelName string
	levelRect     sdl.Rect
}

// NewLevelJumper
// leftTID: the tile ID of left pipe top
func NewLevelJumper(leftTID vector.TileID, nextLevelName string) *levelJumper {
	leftRect := GetTileRect(leftTID)
	return &levelJumper{
		nextLevelName: nextLevelName,
		levelRect: sdl.Rect{
			leftRect.X,
			leftRect.Y - 1, // Y - 1 so hero is able to hit it
			2 * graphic.TILE_SIZE,
			graphic.TILE_SIZE,
		},
	}
}

func (lj *levelJumper) GetRect() sdl.Rect {
	return lj.levelRect
}

func (lj *levelJumper) GetZIndex() int {
	// whatever z-index
	return ZINDEX_1
}

func (lj *levelJumper) Update(ticks uint32, level *Level) {
	// Do nothing
}

func (lj *levelJumper) Draw(camPos vector.Pos) {
	// Do nothing due to invisible
}

func (lj *levelJumper) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	if direction != HIT_FROM_TOP_W_INTENT {
		return
	}

	h.Disable()
	afterEffect := func() {
		level.ShouldSwitchLevel(lj.nextLevelName)
		h.Enable()
	}
	level.AddEffect(NewHeroIntoPipeEffect(h, ticks, afterEffect))
	audio.PlaySound(audio.SOUND_PIPE)
}

func (lj *levelJumper) hitByBottomTile(level *Level, ticks uint32) {
	// Do nothing
}

func (lj *levelJumper) hitByBullet(b bullet, level *Level, ticks uint32) {
	// Do nothing
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// coinEnemy
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var _ Enemy = &coinEnemy{}

type coinEnemy struct {
	*animationTileObj
	basicEnemy
}

func NewCoinEnemy(tid vector.TileID) *coinEnemy {
	reses := []graphic.ResourceID{
		graphic.RESOURCE_TYPE_COIN_0,
		graphic.RESOURCE_TYPE_COIN_1,
		graphic.RESOURCE_TYPE_COIN_2,
		graphic.RESOURCE_TYPE_COIN_3,
	}
	return &coinEnemy{
		animationTileObj: NewAnimationObjectTID(tid, reses, 200, ZINDEX_3),
	}
}

func (ce *coinEnemy) GetRect() sdl.Rect {
	return ce.levelRect
}

func (ce *coinEnemy) GetZIndex() int {
	return ZINDEX_1
}

func (ce *coinEnemy) Update(ticks uint32, level *Level) {
	ce.animationTileObj.Update(ticks, level)
}

func (ce *coinEnemy) Draw(camPos vector.Pos) {
	ce.animationTileObj.Draw(camPos)
}

func (ce *coinEnemy) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	ce.isDead = true
	level.Coins++
	audio.PlaySound(audio.SOUND_COIN)
}

func (ce *coinEnemy) hitByBottomTile(level *Level, ticks uint32) {
	// Do Nothing
}

func (ce *coinEnemy) hitByBullet(b bullet, level *Level, ticks uint32) {
	// Do nothing
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Boss A
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var _ Enemy = &coinEnemy{}

const bossAInitHP = 200

type bossA struct {
	basicEnemy

	resLeft0      graphic.Resource
	resLeft1      graphic.Resource
	resRight0     graphic.Resource
	resRight1     graphic.Resource
	currRes       graphic.Resource
	isFacingRight bool
	levelRect     sdl.Rect
	lastTicks     uint32
	velocity      vector.Vec2D
	hp            int
}

func NewBossA(startPos vector.Pos) *bossA {
	resLeft0 := graphic.Res(graphic.RESOURCE_TYPE_BOSS_A_LEFT_0)
	return &bossA{
		resLeft0:      resLeft0,
		resLeft1:      graphic.Res(graphic.RESOURCE_TYPE_BOSS_A_LEFT_1),
		resRight0:     graphic.Res(graphic.RESOURCE_TYPE_BOSS_A_RIGHT_0),
		resRight1:     graphic.Res(graphic.RESOURCE_TYPE_BOSS_A_RIGHT_1),
		currRes:       resLeft0,
		isFacingRight: true,
		levelRect:     sdl.Rect{startPos.X, startPos.Y, resLeft0.GetW(), resLeft0.GetH()},
		velocity:      vector.Vec2D{80, 0},
		hp:            bossAInitHP,
	}
}

func (b *bossA) GetRect() sdl.Rect {
	return b.levelRect
}

func (b *bossA) GetZIndex() int {
	return ZINDEX_1
}

func (b *bossA) Update(ticks uint32, level *Level) {
	if b.lastTicks == 0 {
		b.lastTicks = ticks
		return
	}

	onHitLeft := func() {
		b.isFacingRight = true
	}
	onHitRight := func() {
		b.isFacingRight = false
	}
	enemySimpleMoveEx(ticks, b.lastTicks, &b.velocity, &b.levelRect, level, onHitLeft, onHitRight)

	// Generate enemies randomly
	if rand.Intn(500) == 7 {
		level.AddEnemy(NewTortoiseEnemyRandomSpeed(
			vector.Pos{b.levelRect.X, b.levelRect.Y}, b.isFacingRight, 100))
	}

	b.updateResource(ticks)

	b.lastTicks = ticks
}

func (b *bossA) Draw(camPos vector.Pos) {
	graphic.DrawResource(b.currRes, b.levelRect, camPos)

	// Draw HP
	outerBox := sdl.Rect{
		b.levelRect.X,
		b.levelRect.Y - 20,
		b.levelRect.W,
		10,
	}
	innerBox := sdl.Rect{
		b.levelRect.X + 1,
		b.levelRect.Y - 19,
		(b.levelRect.W - 2) * int32(b.hp) / bossAInitHP,
		8,
	}
	graphic.DrawRect(outerBox, camPos, 0, 0, 0, 255)
	graphic.FillRect(innerBox, camPos, 255, 0, 0, 255)
}

func (b *bossA) updateResource(ticks uint32) {
	if ticks%1000 < 500 {
		if b.isFacingRight {
			b.currRes = b.resRight0
		} else {
			b.currRes = b.resLeft0
		}
	} else {
		if b.isFacingRight {
			b.currRes = b.resRight1
		} else {
			b.currRes = b.resLeft1
		}
	}
}

func (b *bossA) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	switch direction {
	case HIT_FROM_TOP_W_INTENT:
		// bounce the hero up
		h.velocity.Y = -1200
		audio.PlaySound(audio.SOUND_STOMP)

	default:
		// hero is hurt
		hurtHeroIfIntersectEnoughEx(h, b, level, 0.1)
	}
}

func (b *bossA) hitByBottomTile(level *Level, ticks uint32) {
	// Do Nothing
}

func (boss *bossA) hitByBullet(blt bullet, level *Level, ticks uint32) {
	boss.hp -= blt.GetDamage()
	if boss.hp <= 0 {
		boss.isDead = true
		var dieToRight bool
		if blt.GetRect().X < boss.levelRect.X {
			dieToRight = true
		}

		// show effects
		boomRes := graphic.Res(graphic.RESOURCE_TYPE_BOSS_BOOM)
		level.AddEffect(NewShowOnceEffect(
			boomRes, vector.Vec2D{boss.levelRect.X, boss.levelRect.Y}, sdl.GetTicks(), 1000))
		level.AddEffect(NewDeadDownEffect(boss.currRes, dieToRight, boss.levelRect, ticks))
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helper functions
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

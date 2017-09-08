package level

import (
	"math/rand"

	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/audio"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Boss A
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var _ Enemy = &bossA{}

const bossAInitHP = 1000

type bossA struct {
	basicEnemy
	canSay

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
		canSay:        newCanSay(bossASentences),
		resLeft0:      resLeft0,
		resLeft1:      graphic.Res(graphic.RESOURCE_TYPE_BOSS_A_LEFT_1),
		resRight0:     graphic.Res(graphic.RESOURCE_TYPE_BOSS_A_RIGHT_0),
		resRight1:     graphic.Res(graphic.RESOURCE_TYPE_BOSS_A_RIGHT_1),
		currRes:       resLeft0,
		isFacingRight: false,
		levelRect:     sdl.Rect{startPos.X, startPos.Y, resLeft0.GetW(), resLeft0.GetH()},
		velocity:      vector.Vec2D{-80, 0},
		hp:            bossAInitHP,
	}
}

func (b *bossA) GetRect() sdl.Rect {
	return b.levelRect
}

func (b *bossA) GetZIndex() int {
	return ZINDEX_1
}

func (b *bossA) getSentencePos() vector.Pos {
	return vector.Pos{
		X: b.levelRect.X - 30,
		Y: b.levelRect.Y - 70,
	}
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

	// Generate new self randomly
	if rand.Intn(1900) == 7 {
		pos := vector.Pos{
			X: b.levelRect.X + 2*b.velocity.X,
			Y: b.levelRect.Y,
		}
		level.AddEnemy(NewBossA(pos))
		audio.PlaySound(audio.SOUND_BOSS_LAUGH)
	}

	// Generate new boss B randomly
	if rand.Intn(500) == 7 {
		pos := vector.Pos{
			X: b.levelRect.X + 2*b.velocity.X,
			Y: b.levelRect.Y,
		}
		level.AddEnemy(NewRandomBossB(pos))
	}

	// Keep showing random sentences
	b.say(ticks, level, 100, 256, b.getSentencePos)

	// Randomly change direction
	changeDirectionRandomly(100, &b.isFacingRight, &b.velocity.X)

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
	graphic.FillRect(innerBox, camPos, 156, 54, 181, 255)
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

		// if die, show effects & play sound
		boomRes := graphic.Res(graphic.RESOURCE_TYPE_BOSS_BOOM)
		level.AddEffect(NewShowOnceEffect(
			boomRes, vector.Vec2D{boss.levelRect.X, boss.levelRect.Y}, sdl.GetTicks(), 1000))
		level.AddEffect(NewDeadDownEffect(boss.currRes, dieToRight, boss.levelRect, ticks))
		// play multiple times for better effect
		audio.PlaySound(audio.SOUND_KO)
		audio.PlaySound(audio.SOUND_KO)
		audio.PlaySound(audio.SOUND_KO)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Boss B: Richard's direct reports
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var _ Enemy = &bossB{}

const bossBInitHP = 200

type bossB struct {
	basicEnemy
	canSay

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

func NewRandomBossB(startPos vector.Pos) *bossB {
	userID := richardLeadershipUserIDs[rand.Intn(len(richardLeadershipUserIDs))]
	return NewBossB(startPos, userID)
}

func NewBossB(startPos vector.Pos, userID string) *bossB {
	sentences, ok := bossBSentenceMap[userID]
	if !ok {
		log.Fatalf("Boss B does not support user %s", userID)
	}

	resPack := graphic.GetBossBResPack(userID)
	resLeft0 := resPack.ResLeft0
	return &bossB{
		canSay:        newCanSay(sentences),
		resLeft0:      resLeft0,
		resLeft1:      resPack.ResLeft1,
		resRight0:     resPack.ResRight0,
		resRight1:     resPack.ResRight1,
		currRes:       resLeft0,
		isFacingRight: false,
		levelRect:     sdl.Rect{startPos.X, startPos.Y, resLeft0.GetW(), resLeft0.GetH()},
		velocity:      vector.Vec2D{-80, 0},
		hp:            bossBInitHP,
	}
}

func (b *bossB) GetRect() sdl.Rect {
	return b.levelRect
}

func (b *bossB) GetZIndex() int {
	return ZINDEX_1
}

func (b *bossB) getSentencePos() vector.Pos {
	return vector.Pos{
		X: b.levelRect.X - 30,
		Y: b.levelRect.Y - 70,
	}
}

func (b *bossB) Update(ticks uint32, level *Level) {
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
	if rand.Intn(150) == 7 {
		level.AddEnemy(NewRandomJupiterTortoiseEnemyEx(
			vector.Pos{b.levelRect.X, b.levelRect.Y}, b.isFacingRight, 150))
	}

	// Keep showing random sentences
	b.say(ticks, level, 100, 256, b.getSentencePos)

	// Randomly change direction
	changeDirectionRandomly(100, &b.isFacingRight, &b.velocity.X)

	b.updateResource(ticks)

	b.lastTicks = ticks
}

func (b *bossB) Draw(camPos vector.Pos) {
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
		(b.levelRect.W - 2) * int32(b.hp) / bossBInitHP,
		8,
	}
	graphic.DrawRect(outerBox, camPos, 0, 0, 0, 255)
	graphic.FillRect(innerBox, camPos, 255, 0, 0, 255)
}

func (b *bossB) updateResource(ticks uint32) {
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

func (b *bossB) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
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

func (b *bossB) hitByBottomTile(level *Level, ticks uint32) {
	// Do Nothing
}

func (boss *bossB) hitByBullet(blt bullet, level *Level, ticks uint32) {
	boss.hp -= blt.GetDamage()
	if boss.hp <= 0 {
		boss.isDead = true
		var dieToRight bool
		if blt.GetRect().X < boss.levelRect.X {
			dieToRight = true
		}

		// if die, show effects & play sound
		boomRes := graphic.Res(graphic.RESOURCE_TYPE_BOSS_BOOM)
		level.AddEffect(NewShowOnceEffect(
			boomRes, vector.Vec2D{boss.levelRect.X, boss.levelRect.Y}, sdl.GetTicks(), 1000))
		level.AddEffect(NewDeadDownEffect(boss.currRes, dieToRight, boss.levelRect, ticks))
		// play multiple times for better effect
		audio.PlaySound(audio.SOUND_KO)
		audio.PlaySound(audio.SOUND_KO)
		audio.PlaySound(audio.SOUND_KO)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Boss Sentences
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var bossASentences []string = []string{
	"I have a dream!",
	"Red lobster!!",
	"Let's extend the meeting...",
	"I will ask Peter to fire you...",
	"I love my work!!",
}

var chranSentences []string = []string{
	"I am enjoying my sabbatical",
	"Don't bother me...",
	"I have 25 meetings today :D",
	"My next meeting is in 3 minutes",
	"Have you found the bug??",
	"I really need to fire you...",
	"Don't be shy",
	"Let's schedule a meeting",
	"Your PPT sucks...",
	"Your code works like a shit...",
	"Let me help Richard!",
}

var fchen5Sentences []string = []string{
	"Oh no, there is a variable shift...",
	"No worries I will handle it",
	"Have you tried our variable catalogue?",
	"Let me help Richard!",
}

var xhaoSentences []string = []string{
	"Let me help Richard!",
	"Talk is cheap, show me the code",
}

var qingyliSentences []string = []string{
	"Hope you will enjoy today's pizza!",
	"I really enjoy this innovation day!",
	"Let me help Richard!",
}

// user id -> sentences slide
var bossBSentenceMap map[string][]string = map[string][]string{
	"chran":   chranSentences,
	"fchen5":  fchen5Sentences,
	"xhao":    xhaoSentences,
	"qingyli": qingyliSentences,
}

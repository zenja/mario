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
// Basic Boss
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var _ Enemy = &basicBoss{}

var (
	hpColorRed    = sdl.Color{255, 0, 0, 255}
	hpColorPurple = sdl.Color{156, 54, 181, 255}
)

type basicBoss struct {
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
	maxHP         int
	hp            int
	hpColor       sdl.Color

	extraUpdateActions func(b *basicBoss, level *Level, ticks uint32)
}

func NewBasicBoss(
	startPos vector.Pos,
	resLeft0, resLeft1, resRight0, resRight1 graphic.Resource,
	initHP int,
	hpColor sdl.Color,
	sentences []string,
	extraUpdateActions func(*basicBoss, *Level, uint32)) *basicBoss {

	return &basicBoss{
		canSay:             newCanSay(sentences),
		resLeft0:           resLeft0,
		resLeft1:           resLeft1,
		resRight0:          resRight0,
		resRight1:          resRight1,
		currRes:            resLeft0,
		isFacingRight:      false,
		levelRect:          sdl.Rect{startPos.X, startPos.Y, resLeft0.GetW(), resLeft0.GetH()},
		velocity:           vector.Vec2D{-80, 0},
		maxHP:              initHP,
		hp:                 initHP,
		hpColor:            hpColor,
		extraUpdateActions: extraUpdateActions,
	}
}

func (b *basicBoss) GetRect() sdl.Rect {
	return b.levelRect
}

func (b *basicBoss) GetZIndex() int {
	return ZINDEX_1
}

func (b *basicBoss) getSentencePos() vector.Pos {
	return vector.Pos{
		X: b.levelRect.X - 30,
		Y: b.levelRect.Y - 70,
	}
}

func (b *basicBoss) Update(ticks uint32, level *Level) {
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

	// execute extra actions
	if b.extraUpdateActions != nil {
		b.extraUpdateActions(b, level, ticks)
	}

	// Randomly change direction
	changeDirectionRandomly(100, &b.isFacingRight, &b.velocity.X)

	b.updateResource(ticks)

	b.lastTicks = ticks
}

func (b *basicBoss) Draw(camPos vector.Pos) {
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
		(b.levelRect.W - 2) * int32(b.hp) / int32(b.maxHP),
		8,
	}
	graphic.DrawRect(outerBox, camPos, 0, 0, 0, 255)
	graphic.FillRect(innerBox, camPos, b.hpColor.R, b.hpColor.G, b.hpColor.B, b.hpColor.A)
}

func (b *basicBoss) updateResource(ticks uint32) {
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

func (b *basicBoss) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	switch direction {
	case HIT_FROM_TOP_W_INTENT:
		// bounce the hero up
		h.velocity.Y = -1200
		audio.PlaySound(audio.SOUND_STOMP)
		// reduce some HP from boss
		b.hp -= 20
		if b.hp <= 0 {
			b.die(true, level, ticks)
		}

	default:
		// hero is hurt
		hurtHeroIfIntersectEnoughEx(h, b, level, 0.1)
	}
}

func (b *basicBoss) hitByBottomTile(level *Level, ticks uint32) {
	// Do Nothing
}

func (boss *basicBoss) hitByBullet(blt bullet, level *Level, ticks uint32) {
	boss.hp -= blt.GetDamage()
	var dieToRight bool
	if blt.GetRect().X < boss.levelRect.X {
		dieToRight = true
	}
	if boss.hp <= 0 {
		boss.die(dieToRight, level, ticks)
	}
}

func (boss *basicBoss) die(dieToRight bool, level *Level, ticks uint32) {
	boss.Kill()

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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Boss A
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const bossAInitHP = 700

type bossA struct {
	*basicBoss
}

func NewBossA(startPos vector.Pos) *bossA {
	resLeft0 := graphic.Res(graphic.RESOURCE_TYPE_BOSS_A_LEFT_0)
	basicBoss := NewBasicBoss(
		startPos,
		resLeft0,
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_A_LEFT_1),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_A_RIGHT_0),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_A_RIGHT_1),
		bossAInitHP,
		hpColorPurple,
		bossASentences,
		bossAExtraUpdateActions,
	)
	return &bossA{
		basicBoss: basicBoss,
	}
}

func bossAExtraUpdateActions(b *basicBoss, level *Level, ticks uint32) {
	// Generate new self randomly
	if rand.Intn(2200) == 7 {
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
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Boss B: Richard's direct reports
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const bossBInitHP = 150

type bossB struct {
	*basicBoss
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
	basicBoss := NewBasicBoss(
		startPos,
		resPack.ResLeft0,
		resPack.ResLeft1,
		resPack.ResRight0,
		resPack.ResRight1,
		bossBInitHP,
		hpColorRed,
		sentences,
		bossBExtraUpdateActions,
	)
	return &bossB{
		basicBoss: basicBoss,
	}
}

func bossBExtraUpdateActions(b *basicBoss, level *Level, ticks uint32) {
	// Generate enemies randomly
	if rand.Intn(250) == 7 {
		level.AddEnemy(NewRandomICTortoiseEnemyEx(
			vector.Pos{b.levelRect.X, b.levelRect.Y}, b.isFacingRight, 150))
	}

	// Keep showing random sentences
	b.say(ticks, level, 100, 256, b.getSentencePos)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// BossC
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const bossCInitHP = 400

type bossC struct {
	*basicBoss
}

func NewBossC(startPos vector.Pos) *bossC {
	resLeft0 := graphic.Res(graphic.RESOURCE_TYPE_BOSS_C_LEFT_0)
	basicBoss := NewBasicBoss(
		startPos,
		resLeft0,
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_C_LEFT_1),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_C_RIGHT_0),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_C_RIGHT_1),
		bossCInitHP,
		hpColorPurple,
		bossCSentences,
		bossCExtraUpdateActions,
	)
	return &bossC{
		basicBoss: basicBoss,
	}
}

func bossCExtraUpdateActions(b *basicBoss, level *Level, ticks uint32) {
	fireAroundRandomly(200, level, b, BULLET_ENEMY_SWORD)

	// Keep showing random sentences
	b.say(ticks, level, 100, 256, b.getSentencePos)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// BossD
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const bossDInitHP = 400

type bossD struct {
	*basicBoss
}

func NewBossD(startPos vector.Pos) *bossD {
	resLeft0 := graphic.Res(graphic.RESOURCE_TYPE_BOSS_D_LEFT_0)
	basicBoss := NewBasicBoss(
		startPos,
		resLeft0,
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_D_LEFT_1),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_D_RIGHT_0),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_D_RIGHT_1),
		bossDInitHP,
		hpColorPurple,
		bossDSentences,
		bossDExtraUpdateActions,
	)
	return &bossD{
		basicBoss: basicBoss,
	}
}

func bossDExtraUpdateActions(b *basicBoss, level *Level, ticks uint32) {
	fireAroundRandomly(200, level, b, BULLET_ENEMY_MOON)

	// Keep showing random sentences
	b.say(ticks, level, 100, 256, b.getSentencePos)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// BossE
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const bossEInitHP = 400

type bossE struct {
	*basicBoss
}

func NewBossE(startPos vector.Pos) *bossE {
	resLeft0 := graphic.Res(graphic.RESOURCE_TYPE_BOSS_E_LEFT_0)
	basicBoss := NewBasicBoss(
		startPos,
		resLeft0,
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_E_LEFT_1),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_E_RIGHT_0),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_E_RIGHT_1),
		bossEInitHP,
		hpColorPurple,
		bossESentences,
		bossEExtraUpdateActions,
	)
	return &bossE{
		basicBoss: basicBoss,
	}
}

func bossEExtraUpdateActions(b *basicBoss, level *Level, ticks uint32) {
	fireToHeroRandomly(70, level, b, BULLET_ENEMY_AXE, 350)

	// Keep showing random sentences
	b.say(ticks, level, 100, 256, b.getSentencePos)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// BossF
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const bossFInitHP = 400

type bossF struct {
	*basicBoss
}

func NewBossF(startPos vector.Pos) *bossF {
	resLeft0 := graphic.Res(graphic.RESOURCE_TYPE_BOSS_F_LEFT_0)
	basicBoss := NewBasicBoss(
		startPos,
		resLeft0,
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_F_LEFT_1),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_F_RIGHT_0),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_F_RIGHT_1),
		bossFInitHP,
		hpColorPurple,
		bossFSentences,
		bossFExtraUpdateActions,
	)
	return &bossF{
		basicBoss: basicBoss,
	}
}

func bossFExtraUpdateActions(b *basicBoss, level *Level, ticks uint32) {
	fireToHeroRandomly(70, level, b, BULLET_ENEMY_CHERRY, 350)

	// Keep showing random sentences
	b.say(ticks, level, 100, 256, b.getSentencePos)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// BossG
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const bossGInitHP = 400

type bossG struct {
	*basicBoss
}

func NewBossG(startPos vector.Pos) *bossG {
	resLeft0 := graphic.Res(graphic.RESOURCE_TYPE_BOSS_G_LEFT_0)
	basicBoss := NewBasicBoss(
		startPos,
		resLeft0,
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_G_LEFT_1),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_G_RIGHT_0),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_G_RIGHT_1),
		bossGInitHP,
		hpColorPurple,
		bossGSentences,
		bossGExtraUpdateActions,
	)
	return &bossG{
		basicBoss: basicBoss,
	}
}

func bossGExtraUpdateActions(b *basicBoss, level *Level, ticks uint32) {
	fireToHeroRandomly(70, level, b, BULLET_ENEMY_SKULL, 350)

	// Keep showing random sentences
	b.say(ticks, level, 100, 256, b.getSentencePos)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// BossH
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const bossHInitHP = 400

type bossH struct {
	*basicBoss
}

func NewBossH(startPos vector.Pos) *bossH {
	resLeft0 := graphic.Res(graphic.RESOURCE_TYPE_BOSS_H_LEFT_0)
	basicBoss := NewBasicBoss(
		startPos,
		resLeft0,
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_H_LEFT_1),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_H_RIGHT_0),
		graphic.Res(graphic.RESOURCE_TYPE_BOSS_H_RIGHT_1),
		bossHInitHP,
		hpColorPurple,
		bossHSentences,
		bossHExtraUpdateActions,
	)
	return &bossH{
		basicBoss: basicBoss,
	}
}

func bossHExtraUpdateActions(b *basicBoss, level *Level, ticks uint32) {
	fireAroundRandomly(200, level, b, BULLET_ENEMY_APPLE)

	// Keep showing random sentences
	b.say(ticks, level, 100, 256, b.getSentencePos)
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

var bossCSentences []string = []string{}
var bossDSentences []string = []string{}
var bossESentences []string = []string{}
var bossFSentences []string = []string{}
var bossGSentences []string = []string{}
var bossHSentences []string = []string{}

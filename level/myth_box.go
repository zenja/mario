package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

// assert that mythBox is hit-able by hero
var _ hittableByHero = &mythBox{}

type mythBox struct {
	// resources
	resNormal      graphic.Resource
	resNormalLight graphic.Resource
	resEmpty       graphic.Resource // empty, no coins
	currRes        graphic.Resource

	// myth box has both a tile rect and current level rect,
	// because we allow myth box to move a little bit after being hit
	tileRect  sdl.Rect
	levelRect sdl.Rect

	actor mythBoxActor

	isBounding bool
	isEmpty    bool
	velocity   vector.Vec2D
	lastTicks  uint32
}

type mythBoxActor interface {
	onEffectiveBottomHit(mb *mythBox, level *Level, ticks uint32)
	onBoundingFinished(mb *mythBox, level *Level, ticks uint32)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Coin actor
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// coinActor is an mythBoxActor
var _ mythBoxActor = &coinActor{}

type coinActor struct {
	numCoinsLeft int
}

func (ca *coinActor) onEffectiveBottomHit(mb *mythBox, level *Level, ticks uint32) {
	if ca.numCoinsLeft > 0 {
		// add a coin effect
		mbTID := GetTileID(vector.Pos{mb.tileRect.X, mb.tileRect.Y}, false, false)
		level.AddEffect(NewCoinEffect(vector.TileID{mbTID.X, mbTID.Y - 1}, level.ResourceRegistry, ticks))
	}
}

func (ca *coinActor) onBoundingFinished(mb *mythBox, level *Level, ticks uint32) {
	ca.numCoinsLeft--
	if ca.numCoinsLeft <= 0 {
		mb.Empty()
	}
}

func NewCoinMythBox(startPos vector.Pos, numCoins int, resourceRegistry map[graphic.ResourceID]graphic.Resource) *mythBox {
	actor := coinActor{numCoinsLeft: numCoins}
	return newMythBox(startPos, &actor, resourceRegistry)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Mushroom actor
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var _ mythBoxActor = &mushroomActor{}

type mushroomActor struct {
	enemy *goodMushroom
}

func (ma *mushroomActor) onEffectiveBottomHit(mb *mythBox, level *Level, ticks uint32) {
	level.AddEnemy(ma.enemy)
}

func (ma *mushroomActor) onBoundingFinished(mb *mythBox, level *Level, ticks uint32) {
	mb.Empty()
}

func NewMushroomMythBox(startPos vector.Pos, resourceRegistry map[graphic.ResourceID]graphic.Resource) *mythBox {
	enemyStartPos := vector.Pos{startPos.X, startPos.Y - graphic.TILE_SIZE}
	enemy := NewGoodMushroom(enemyStartPos, resourceRegistry)
	actor := mushroomActor{enemy}
	return newMythBox(startPos, &actor, resourceRegistry)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Myth box methods
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func newMythBox(startPos vector.Pos, actor mythBoxActor, resourceRegistry map[graphic.ResourceID]graphic.Resource) *mythBox {
	resNormal, _ := resourceRegistry[graphic.RESOURCE_TYPE_MYTH_BOX_NORMAL]
	resNormalLight, _ := resourceRegistry[graphic.RESOURCE_TYPE_MYTH_BOX_NORMAL_LIGHT]
	resEmpty, _ := resourceRegistry[graphic.RESOURCE_TYPE_MYTH_BOX_EMPTY]
	tileRect := sdl.Rect{startPos.X, startPos.Y, graphic.TILE_SIZE, graphic.TILE_SIZE}
	return &mythBox{
		resNormal:      resNormal,
		resNormalLight: resNormalLight,
		resEmpty:       resEmpty,
		currRes:        resNormal,
		tileRect:       tileRect,
		levelRect:      tileRect,
		actor:          actor,
	}
}

func (mb *mythBox) Draw(g *graphic.Graphic, camPos vector.Pos) {
	g.DrawResource(mb.currRes, mb.levelRect, camPos)
}

func (mb *mythBox) Update(events *intsets.Sparse, ticks uint32, level *Level) {
	if mb.lastTicks == 0 {
		mb.lastTicks = ticks
		return
	}

	// if has coin, show blink animation
	if !mb.isEmpty {
		// update res for blink animation
		if ticks%400 < 200 {
			mb.currRes = mb.resNormal
		} else {
			mb.currRes = mb.resNormalLight
		}
	}

	if mb.isBounding {
		gravity := vector.Vec2D{0, 10}
		mb.velocity.Add(gravity)

		velocityStep := CalcVelocityStep(mb.velocity, ticks, mb.lastTicks, nil)

		// apply velocity step
		mb.levelRect.X += velocityStep.X
		mb.levelRect.Y += velocityStep.Y

		// if reach origin (Y) position, the bounding is stopped
		if mb.levelRect.Y >= mb.tileRect.Y {
			mb.levelRect.Y = mb.tileRect.Y
			mb.isBounding = false
			mb.actor.onBoundingFinished(mb, level, ticks)
		}
	} else {
		mb.levelRect.X = mb.tileRect.X
		mb.levelRect.Y = mb.tileRect.Y
	}

	mb.lastTicks = ticks
}

func (mb *mythBox) GetRect() sdl.Rect {
	return mb.levelRect
}

func (mb *mythBox) GetZIndex() int {
	return ZINDEX_1
}

func (mb *mythBox) Empty() {
	mb.isEmpty = true
	mb.currRes = mb.resEmpty
}

func (mb *mythBox) IsEmpty() bool {
	return mb.isEmpty
}

func (mb *mythBox) StartBounding() {
	mb.isBounding = true
	mb.velocity.Y = -100
}

func (mb *mythBox) StopBoundingIfNeeded() {
	// if reach origin (Y) position, the bounding is stopped
	if mb.levelRect.Y >= mb.tileRect.Y {
		mb.levelRect.Y = mb.tileRect.Y
		mb.isBounding = false
	}
}

func (mb *mythBox) IsBounding() bool {
	return mb.isBounding
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Private major methods
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (mb *mythBox) hitByHero(h *Hero, direction hitDirection, level *Level, ticks uint32) {
	// can only be hit from bottom
	if direction != HIT_FROM_BOTTOM_W_INTENT {
		return
	}

	// empty myth box won't react to hit
	if mb.isEmpty {
		return
	}

	// only react if box is not bounding, to avoid bounding on bounding
	if !mb.isBounding {
		mb.StartBounding()
		mb.actor.onEffectiveBottomHit(mb, level, ticks)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Private helper methods
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

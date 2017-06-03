package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

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

	numCoinsLeft int

	isBounding bool
	velocity   vector.Vec2D
	lastTicks  uint32
}

func NewMythBox(startPos vector.Pos, numCoins int, resourceRegistry map[graphic.ResourceID]graphic.Resource) Object {
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
		numCoinsLeft:   numCoins,
	}
}

func (mb *mythBox) Draw(g *graphic.Graphic, camPos vector.Pos) {
	drawResource(g, mb.currRes, mb.levelRect, camPos)
}

func (mb *mythBox) Update(events *intsets.Sparse, ticks uint32, level *Level) {
	if mb.lastTicks == 0 {
		mb.lastTicks = ticks
		return
	}

	// if has coin, show blink animation
	if mb.numCoinsLeft > 0 {
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

		velocityStep := mb.velocity
		velocityStep.Multiply(int32(ticks - mb.lastTicks))
		velocityStep.Divide(1000)

		// apply velocity step
		mb.levelRect.X += velocityStep.X
		mb.levelRect.Y += velocityStep.Y

		// update coin, set empty res if no coin
		mb.numCoinsLeft--
		if mb.numCoinsLeft <= 0 {
			mb.currRes = mb.resEmpty
		}

		// if reach origin (Y) position, the bounding is stopped
		if mb.levelRect.Y >= mb.tileRect.Y {
			mb.levelRect.Y = mb.tileRect.Y
			mb.isBounding = false
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Private major methods
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (mb *mythBox) hitByHero() {
	if !mb.isBounding && mb.numCoinsLeft > 0 {
		mb.isBounding = true
		mb.velocity.Y = -100
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Private helper methods
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

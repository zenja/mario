package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

var _ Effect = &heroIntoPipeEffect{}

type heroIntoPipeEffect struct {
	res            graphic.Resource
	levelRect      sdl.Rect
	startTicks     uint32
	lastTicks      uint32
	finished       bool
	onFinishedHook func()
}

func NewHeroIntoPipeEffect(h *Hero, ticks uint32, onFinished func()) *heroIntoPipeEffect {
	return &heroIntoPipeEffect{
		res:            h.currRes,
		levelRect:      h.getRenderRect(),
		startTicks:     ticks,
		lastTicks:      ticks,
		onFinishedHook: onFinished,
	}
}

func (hipe *heroIntoPipeEffect) Update(ticks uint32) {
	if ticks-hipe.startTicks > 1500 {
		hipe.finished = true
		return
	}

	velocity := vector.Vec2D{0, 100}
	velStep := CalcVelocityStep(velocity, ticks, hipe.lastTicks, nil)
	hipe.levelRect.X += velStep.X
	hipe.levelRect.Y += velStep.Y
	hipe.levelRect.W -= velStep.X
	hipe.levelRect.H -= velStep.Y

	hipe.lastTicks = ticks
}

func (hipe *heroIntoPipeEffect) Draw(camPos vector.Pos, ticks uint32) {
	if !hipe.Finished() {
		graphic.DrawResource(hipe.res, hipe.levelRect, camPos)
	}
}

func (hipe *heroIntoPipeEffect) Finished() bool {
	return hipe.finished
}

func (hipe *heroIntoPipeEffect) OnFinished() {
	if hipe.onFinishedHook != nil {
		hipe.onFinishedHook()
	}
}

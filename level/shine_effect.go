package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

var _ Effect = &shineEffect{}

const (
	shine_effect_duraiton_ms = 1000
)

type shineEffect struct {
	res0       graphic.Resource
	res1       graphic.Resource
	res2       graphic.Resource
	currRes0   graphic.Resource
	currRes1   graphic.Resource
	currRes2   graphic.Resource
	hero       *Hero
	startTicks uint32
	finished   bool
}

func NewShineEffect(h *Hero, ticks uint32) *shineEffect {
	return &shineEffect{
		res0:       graphic.Res(graphic.RESOURCE_TYPE_SHINE_0),
		res1:       graphic.Res(graphic.RESOURCE_TYPE_SHINE_1),
		res2:       graphic.Res(graphic.RESOURCE_TYPE_SHINE_2),
		hero:       h,
		startTicks: ticks,
	}
}

func (se *shineEffect) Update(ticks uint32) {
	if ticks-se.startTicks > shine_effect_duraiton_ms {
		se.finished = true
		return
	}

	se.currRes0 = se.getRes(ticks, se.res0, se.res1, se.res2)
	se.currRes1 = se.getRes(ticks, se.res1, se.res2, se.res0)
	se.currRes2 = se.getRes(ticks, se.res2, se.res0, se.res1)
}

func (se *shineEffect) Draw(camPos vector.Pos, ticks uint32) {
	graphic.DrawResource(se.currRes0, se.getRect0(), camPos)
	graphic.DrawResource(se.currRes1, se.getRect1(), camPos)
	graphic.DrawResource(se.currRes2, se.getRect2(), camPos)
}

func (se *shineEffect) Finished() bool {
	return se.finished
}

func (se *shineEffect) getRect0() sdl.Rect {
	heroRect := se.hero.levelRect
	return sdl.Rect{
		heroRect.X - 15,
		heroRect.Y - 25,
		se.res0.GetW(),
		se.res0.GetH(),
	}
}

func (se *shineEffect) getRect1() sdl.Rect {
	heroRect := se.hero.levelRect
	return sdl.Rect{
		heroRect.X + 30,
		heroRect.Y + 10,
		se.res0.GetW(),
		se.res0.GetH(),
	}
}

func (se *shineEffect) getRect2() sdl.Rect {
	heroRect := se.hero.levelRect
	return sdl.Rect{
		heroRect.X + 5,
		heroRect.Y + 55,
		se.res0.GetW(),
		se.res0.GetH(),
	}
}

func (se *shineEffect) getRes(ticks uint32, resA, resB, resC graphic.Resource) graphic.Resource {
	residual := (ticks - se.startTicks) % 600
	switch {
	case residual < 200:
		return resA
	case residual < 400:
		return resB
	default:
		return resC
	}
}

func (se *shineEffect) OnFinished() {
	// Do nothing
}

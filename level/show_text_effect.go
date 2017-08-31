package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

// showTextEffect is an Effect
var _ Effect = &showTextEffect{}

// showTextEffect shows a resource for a while and then disappear
type showTextEffect struct {
	text       string
	color      sdl.Color
	posGetter  func() vector.Pos
	startTicks uint32
	durationMs uint32
	finished   bool
}

func NewShowTextEffect(text string, color sdl.Color, posGetter func() vector.Pos, ticks uint32, durationMs uint32) *showTextEffect {
	return &showTextEffect{
		text:       text,
		color:      color,
		posGetter:  posGetter,
		startTicks: ticks,
		durationMs: durationMs,
		finished:   false,
	}
}

func (ste *showTextEffect) Update(ticks uint32) {
	if ticks-ste.startTicks > ste.durationMs {
		ste.finished = true
	}
}

func (ste *showTextEffect) Draw(camPos vector.Pos, ticks uint32) {
	if !ste.Finished() {
		graphic.DrawTextRelative(ste.text, ste.posGetter(), camPos, ste.color)
	}
}

func (ste *showTextEffect) Finished() bool {
	return ste.finished
}

func (ste *showTextEffect) OnFinished() {
	// Do nothing
}

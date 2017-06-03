package level

import (
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type Effect interface {
	UpdateAndDraw(g *graphic.Graphic, camPos vector.Pos, ticks uint32)
	Finished() bool
}

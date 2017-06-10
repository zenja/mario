package level

import (
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type Effect interface {
	Update(ticks uint32)
	Draw(g *graphic.Graphic, camPos vector.Pos, ticks uint32)
	Finished() bool
}

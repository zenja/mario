package level

import "github.com/zenja/mario/vector"

type Effect interface {
	Update(ticks uint32)
	Draw(camPos vector.Pos, ticks uint32)
	Finished() bool
	OnFinished()
}

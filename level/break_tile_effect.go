package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

// breakTileEffect is an Effect
var _ Effect = &breakTileEffect{}

type breakTileEffect struct {
	pieceRes graphic.Resource

	rectLT sdl.Rect
	rectRT sdl.Rect
	rectLB sdl.Rect
	rectRB sdl.Rect

	// velocities
	velLT vector.Vec2D
	velRT vector.Vec2D
	velLB vector.Vec2D
	velRB vector.Vec2D

	startTicks uint32
	lastTicks  uint32
	finished   bool
}

func NewBreakTileEffect(pieceRes graphic.Resource, tid vector.TileID, ticks uint32) *breakTileEffect {
	tileRect := GetTileRect(tid)
	var size int32 = graphic.TILE_SIZE / 2
	return &breakTileEffect{
		pieceRes:   pieceRes,
		rectLT:     sdl.Rect{tileRect.X, tileRect.Y, size, size},
		rectRT:     sdl.Rect{tileRect.X + size, tileRect.Y, size, size},
		rectLB:     sdl.Rect{tileRect.X, tileRect.Y + size, size, size},
		rectRB:     sdl.Rect{tileRect.X + size, tileRect.Y + size, size, size},
		velLT:      vector.Vec2D{-500, -1000},
		velRT:      vector.Vec2D{500, -1000},
		velLB:      vector.Vec2D{-250, -1000},
		velRB:      vector.Vec2D{250, -1000},
		startTicks: ticks,
		lastTicks:  ticks,
		finished:   false,
	}
}

func (bte *breakTileEffect) Update(ticks uint32) {
	vels := []*vector.Vec2D{
		&bte.velLT,
		&bte.velRT,
		&bte.velLB,
		&bte.velRB,
	}
	rects := []*sdl.Rect{
		&bte.rectLT,
		&bte.rectRT,
		&bte.rectLB,
		&bte.rectRB,
	}
	for i := range vels {
		vels[i].Y += 50

		velocityStep := CalcVelocityStep(*vels[i], ticks, bte.lastTicks, nil)
		rects[i].X += velocityStep.X
		rects[i].Y += velocityStep.Y
	}

	if ticks-bte.startTicks > 1000 {
		bte.finished = true
	}

	bte.lastTicks = ticks
}

func (bte *breakTileEffect) Draw(g *graphic.Graphic, camPos vector.Pos, ticks uint32) {

	if !bte.Finished() {
		g.DrawResource(bte.pieceRes, bte.rectLT, camPos)
		g.DrawResource(bte.pieceRes, bte.rectRT, camPos)
		g.DrawResource(bte.pieceRes, bte.rectLB, camPos)
		g.DrawResource(bte.pieceRes, bte.rectRB, camPos)
	}
}

func (bte *breakTileEffect) Finished() bool {
	return bte.finished
}

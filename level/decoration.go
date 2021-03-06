package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

var _ Object = &decoration{}

type decoration struct {
	reses      []graphic.Resource
	currResIdx int
	startPos   vector.Pos
	frameMs    uint32
}

func NewDecoration(tid vector.TileID, resIDs []graphic.ResourceID, frameMs uint32) *decoration {

	var reses []graphic.Resource
	for _, id := range resIDs {
		reses = append(reses, graphic.Res(id))
	}

	tidRect := GetTileRect(tid)
	startPos := vector.Pos{
		tidRect.X,
		// make sure the bottom is on a tile
		tidRect.Y + (graphic.TILE_SIZE - reses[0].GetH()),
	}

	return &decoration{
		reses:      reses,
		currResIdx: 0,
		startPos:   startPos,
		frameMs:    frameMs,
	}
}

func (d *decoration) Draw(camPos vector.Pos) {
	graphic.DrawResource(d.reses[d.currResIdx], d.GetRect(), camPos)
}

func (d *decoration) Update(ticks uint32, _ *Level) {
	d.currResIdx = int((ticks / d.frameMs) % uint32(len(d.reses)))
}

func (d *decoration) GetRect() sdl.Rect {
	return sdl.Rect{d.startPos.X, d.startPos.Y, d.reses[d.currResIdx].GetW(), d.reses[d.currResIdx].GetH()}
}

func (d *decoration) GetZIndex() int {
	return ZINDEX_0
}

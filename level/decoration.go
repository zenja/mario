package level

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

var _ Object = &decoration{}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// decoration
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// textDecoration
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type textDecoration struct {
	sentences    []string
	currIdx      int
	startPos     vector.Pos
	lastSayTicks uint32
}

func NewPrincessTextDecoration(tid vector.TileID) *textDecoration {
	sentences := []string{
		"    Congratulations, you are fired!!",
		"         Welcome to Shanghai!",
		"           We love PayPal!",
		"      Hope you enjoy the journey!",
		"              Surprise!",
	}
	return NewTextDecoration(tid, sentences)
}

func NewTextDecoration(tid vector.TileID, sentences []string) *textDecoration {
	tidRect := GetTileRect(tid)
	startPos := vector.Pos{
		tidRect.X,
		tidRect.Y,
	}

	return &textDecoration{
		sentences: sentences,
		startPos:  startPos,
	}
}

func (td *textDecoration) Draw(camPos vector.Pos) {
}

func (td *textDecoration) Update(ticks uint32, level *Level) {
	if len(td.sentences) == 0 {
		return
	}
	randColor := sdl.Color{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		255,
	}
	s := td.sentences[td.currIdx]
	if ticks-td.lastSayTicks > 3000 {
		level.AddEffect(NewShowTextEffect(s, randColor, td.getPos, ticks, 2000))
		td.lastSayTicks = ticks
		td.currIdx++
		if td.currIdx >= len(td.sentences) {
			td.currIdx = td.currIdx % len(td.sentences)
		}
	}
}

func (td *textDecoration) getPos() vector.Pos {
	return td.startPos
}

func (td *textDecoration) GetRect() sdl.Rect {
	return sdl.Rect{td.startPos.X, td.startPos.Y, 500, 50}
}

func (td *textDecoration) GetZIndex() int {
	return ZINDEX_0
}

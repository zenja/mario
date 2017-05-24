package object

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
)

type Object interface {
	Draw(g *graphic.Graphic, xStart, yStart int32)
}

type singleTileObject struct {
	tileID graphic.TileID
	// the position on the level
	levelPos *sdl.Rect
}

func NewSingleTileObject(tileID graphic.TileID, xStart, yStart int32) Object {
	return &singleTileObject{
		tileID:   tileID,
		levelPos: &sdl.Rect{xStart, yStart, graphic.TILE_SIZE, graphic.TILE_SIZE},
	}
}

// DrawObject draws an object to a given camera screen (xStart, yStart, graphic.SCREEN_WIDTH, graphic.SCREEN_HEIGHT)
func (sto *singleTileObject) Draw(g *graphic.Graphic, xStart, yStart int32) {
	rectInTile, rectInCamera := visibleRectInCamera(sto.levelPos, xStart, yStart)
	if rectInTile != nil {
		g.RenderTile(sto.tileID, rectInTile, rectInCamera)
	}
}

// visibleRectInCamera returns a rect relative to camera which is (partly) visible
// return nil if the rect is not visible in camera at all
func visibleRectInCamera(rect *sdl.Rect, xCameraStart, yCameraStart int32) (rectInTile *sdl.Rect, rectInCamera *sdl.Rect) {
	if rect.X+rect.W < xCameraStart || rect.X > xCameraStart+graphic.SCREEN_WIDTH ||
		rect.Y+rect.H < yCameraStart || rect.Y > yCameraStart+graphic.SCREEN_HEIGHT {
		return nil, nil
	}

	xStartInLevel := min(rect.X, xCameraStart)
	xEndInLevel := min(rect.X+rect.W, xCameraStart+graphic.SCREEN_WIDTH)
	yStartInLevel := min(rect.Y, yCameraStart)
	yEndInLevel := min(rect.Y+rect.H, yCameraStart+graphic.SCREEN_HEIGHT)

	rectInTile = &sdl.Rect{
		xStartInLevel - rect.X,
		yStartInLevel - rect.Y,
		xEndInLevel - xStartInLevel,
		yEndInLevel - yStartInLevel,
	}
	rectInCamera = &sdl.Rect{
		xStartInLevel - xCameraStart,
		yStartInLevel - yCameraStart,
		xEndInLevel - xStartInLevel,
		yEndInLevel - yStartInLevel,
	}
	return
}

func min(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

package object

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
)

type Object interface {
	Draw(g *graphic.Graphic, xCamStart, yCamStart int32)
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
func (sto *singleTileObject) Draw(g *graphic.Graphic, xCamStart, yCamStart int32) {
	drawTile(g, sto.tileID, sto.levelPos, xCamStart, yCamStart)
}

// drawTile is a helper function to draw a tile on level to camera
func drawTile(g *graphic.Graphic, tileID graphic.TileID, levelPos *sdl.Rect, xCamStart, yCamStart int32) {
	rectInTile, rectInCamera := visibleRectInCamera(levelPos, xCamStart, yCamStart)
	if rectInTile != nil {
		g.RenderTile(tileID, rectInTile, rectInCamera)
	}
}

// visibleRectInCamera returns a rect relative to camera which is (partly) visible
// return nil if the rect is not visible in camera at all
func visibleRectInCamera(rect *sdl.Rect, xCamStart, yCamStart int32) (rectInTile *sdl.Rect, rectInCamera *sdl.Rect) {
	if rect.X+rect.W < xCamStart || rect.X > xCamStart+graphic.SCREEN_WIDTH ||
		rect.Y+rect.H < yCamStart || rect.Y > yCamStart+graphic.SCREEN_HEIGHT {
		return nil, nil
	}

	xStartInLevel := max(rect.X, xCamStart)
	xEndInLevel := min(rect.X+rect.W, xCamStart+graphic.SCREEN_WIDTH)
	yStartInLevel := max(rect.Y, yCamStart)
	yEndInLevel := min(rect.Y+rect.H, yCamStart+graphic.SCREEN_HEIGHT)

	rectInTile = &sdl.Rect{
		xStartInLevel - rect.X,
		yStartInLevel - rect.Y,
		xEndInLevel - xStartInLevel,
		yEndInLevel - yStartInLevel,
	}
	rectInCamera = &sdl.Rect{
		xStartInLevel - xCamStart,
		yStartInLevel - yCamStart,
		xEndInLevel - xStartInLevel,
		yEndInLevel - yStartInLevel,
	}
	//fmt.Printf("Camera: %d, %d\n", xCamStart, yCamStart)
	//fmt.Printf("Object rect: %d, %d, %d, %d\n", rect.X, rect.Y, rect.W, rect.H)
	//fmt.Printf("Rect in level: %d, %d, %d, %d\n", xStartInLevel, yStartInLevel, xEndInLevel-xStartInLevel, yEndInLevel-yStartInLevel)
	//fmt.Printf("Rect in tile: %d, %d, %d, %d\n", rectInTile.X, rectInTile.Y, rectInTile.W, rectInTile.H)
	//fmt.Printf("Rect in Camera: %d, %d, %d, %d\n", rectInCamera.X, rectInCamera.Y, rectInCamera.W, rectInCamera.H)
	//fmt.Println()
	return
}

func min(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func max(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

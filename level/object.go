package level

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/math_utils"
	"github.com/zenja/mario/vector"
	"golang.org/x/tools/container/intsets"
)

// Note that Z-Index must start from 0 and increase one-by-one
const ZINDEX_NUM = 5
const (
	ZINDEX_0 = iota
	ZINDEX_1
	ZINDEX_2
	ZINDEX_3
	ZINDEX_4
)

type Object interface {
	Draw(g *graphic.Graphic, camPos vector.Pos)
	Update(events *intsets.Sparse, ticks uint32, level *Level)
	GetRect() sdl.Rect
	GetZIndex() int
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Single tile object
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type singleTileObject struct {
	resource graphic.Resource
	// the position rect on the level
	levelRect sdl.Rect
	zIndex    int
}

func NewSingleTileObject(resource graphic.Resource, startPos vector.Pos, zIndex int) Object {
	return &singleTileObject{
		resource:  resource,
		levelRect: sdl.Rect{startPos.X, startPos.Y, resource.GetW(), resource.GetH()},
		zIndex:    zIndex,
	}
}

func (sto *singleTileObject) Draw(g *graphic.Graphic, camPos vector.Pos) {
	drawResource(g, sto.resource, sto.levelRect, camPos)
}

func (sto *singleTileObject) Update(events *intsets.Sparse, ticks uint32, level *Level) {
	// Do nothing
}

func (sto *singleTileObject) GetRect() sdl.Rect {
	return sto.levelRect
}

func (sto *singleTileObject) GetZIndex() int {
	return sto.zIndex
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helper functions
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// drawResource is a helper function to draw a resource on level to camera
func drawResource(g *graphic.Graphic, resource graphic.Resource, levelPos sdl.Rect, camPos vector.Pos) {
	rectInResource, rectInCamera := visibleRectInCamera(levelPos, camPos.X, camPos.Y)
	if rectInResource != nil {
		g.RenderResource(resource, rectInResource, rectInCamera)
	}
}

// visibleRectInCamera returns a rect relative to camera which is (partly) visible
// return nil if the rect is not visible in camera at all
func visibleRectInCamera(rect sdl.Rect, xCamStart, yCamStart int32) (rectInTile *sdl.Rect, rectInCamera *sdl.Rect) {
	if rect.X+rect.W < xCamStart || rect.X > xCamStart+graphic.SCREEN_WIDTH ||
		rect.Y+rect.H < yCamStart || rect.Y > yCamStart+graphic.SCREEN_HEIGHT {
		return nil, nil
	}

	xStartInLevel := math_utils.Max(rect.X, xCamStart)
	xEndInLevel := math_utils.Min(rect.X+rect.W, xCamStart+graphic.SCREEN_WIDTH)
	yStartInLevel := math_utils.Max(rect.Y, yCamStart)
	yEndInLevel := math_utils.Min(rect.Y+rect.H, yCamStart+graphic.SCREEN_HEIGHT)

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

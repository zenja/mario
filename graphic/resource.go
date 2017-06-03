package graphic

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/zenja/mario/math_utils"
	"github.com/zenja/mario/vector"
)

type ResourceID int

type Resource interface {
	GetTexture() *sdl.Texture
	GetW() int32
	GetH() int32
}

// Resource IDs
const (
	RESOURCE_TYPE_BRICK = iota

	RESOURCE_TYPE_GROUD_GRASS_LEFT
	RESOURCE_TYPE_GROUD_GRASS_MID
	RESOURCE_TYPE_GROUD_GRASS_RIGHT

	RESOURCE_TYPE_GROUD_INNER_MID

	RESOURCE_TYPE_MYTH_BOX_NORMAL
	RESOURCE_TYPE_MYTH_BOX_NORMAL_LIGHT
	RESOURCE_TYPE_MYTH_BOX_EMPTY

	RESOURCE_TYPE_COIN

	RESOURCE_TYPE_HERO_STAND_LEFT
	RESOURCE_TYPE_HERO_WALKING_LEFT
	RESOURCE_TYPE_HERO_STAND_RIGHT
	RESOURCE_TYPE_HERO_WALKING_RIGHT
)

const TILE_SIZE = 50

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TileResource
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type TileResource struct {
	texture *sdl.Texture
}

func (tr *TileResource) GetTexture() *sdl.Texture {
	return tr.texture
}

func (tr *TileResource) GetW() int32 {
	return TILE_SIZE
}

func (tr *TileResource) GetH() int32 {
	return TILE_SIZE
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// NonTileResource
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type NonTileResource struct {
	texture *sdl.Texture
	w, h    int32
}

func (ntr *NonTileResource) GetTexture() *sdl.Texture {
	return ntr.texture
}

func (ntr *NonTileResource) GetW() int32 {
	return ntr.w
}

func (ntr *NonTileResource) GetH() int32 {
	return ntr.h
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Public helper methods
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// visibleRectInCamera returns a rect relative to camera which is (partly) visible
// return nil if the rect is not visible in camera at all
func VisibleRectInCamera(rect sdl.Rect, xCamStart, yCamStart int32) (rectInTile *sdl.Rect, rectInCamera *sdl.Rect) {
	if rect.X+rect.W < xCamStart || rect.X > xCamStart+SCREEN_WIDTH ||
		rect.Y+rect.H < yCamStart || rect.Y > yCamStart+SCREEN_HEIGHT {
		return nil, nil
	}

	xStartInLevel := math_utils.Max(rect.X, xCamStart)
	xEndInLevel := math_utils.Min(rect.X+rect.W, xCamStart+SCREEN_WIDTH)
	yStartInLevel := math_utils.Max(rect.Y, yCamStart)
	yEndInLevel := math_utils.Min(rect.Y+rect.H, yCamStart+SCREEN_HEIGHT)

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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Graphic methods relative to resource
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (g *Graphic) GetResource(resourceID ResourceID) Resource {
	r, ok := g.ResourceRegistry[resourceID]
	if !ok {
		log.Fatalf("ResourceID %d not found in resource registry", resourceID)
	}
	return r
}

// drawResource is a helper function to draw a resource on level to camera
func (g *Graphic) DrawResource(resource Resource, levelPos sdl.Rect, camPos vector.Pos) {
	rectInResource, rectInCamera := VisibleRectInCamera(levelPos, camPos.X, camPos.Y)
	if rectInResource != nil {
		g.RenderResource(resource, rectInResource, rectInCamera)
	}
}

// registerTileResource loads a sprite into a TileResource from a file
func (g *Graphic) registerTileResource(filename string, id ResourceID) {
	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	g.registerTileFromSurface(surface, id)
}

// registerNonTailResource loads a sprite into a NonTileResource from a file
func (g *Graphic) registerNonTailResource(filename string, id ResourceID) {
	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	g.registerNonTileFromSurface(surface, id)
}

func (g *Graphic) registerFlippedNonTailResource(filename string, id ResourceID, flipHorizontal bool) {
	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	g.registerFlippedNonTileFromSurface(surface, id, flipHorizontal)
}

// registerTileFromSurface loads a sprite from a surface into a TileResource object
// User need to free the surface himself
func (g *Graphic) registerTileFromSurface(surface *sdl.Surface, id ResourceID) {
	g.registerResourceFromSurface(surface, id, TILE_SIZE, TILE_SIZE, true)
}

// registerNonTileFromSurface loads a sprite from a surface into a NonTileResource object
// User need to free the surface himself
func (g *Graphic) registerNonTileFromSurface(surface *sdl.Surface, id ResourceID) {
	g.registerResourceFromSurface(surface, id, surface.W, surface.H, false)
}

func (g *Graphic) registerFlippedNonTileFromSurface(surface *sdl.Surface, id ResourceID, flipHorizontal bool) {
	g.registerFlippedResourceFromSurface(surface, id, surface.W, surface.H, flipHorizontal, false)
}

// registerResourceFromSurface loads a sprite from a surface into a Resource object
// User need to free the surface himself
func (g *Graphic) registerResourceFromSurface(surface *sdl.Surface, id ResourceID, width, height int32, isTile bool) {
	if isTile && (width != TILE_SIZE || height != TILE_SIZE) {
		log.Fatalf("declared to be tile but has wrong width (%d) or height (%d)", width, height)
	}

	texture, err := g.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatal(err)
	}

	// make sure the tile is in good shape
	if surface.W != width || surface.H != height {
		oldTexture := texture
		texture, err = g.clipTexture(oldTexture, &sdl.Rect{0, 0, width, height})
		if err != nil {
			log.Fatal(err)
		}
		// release original texture
		oldTexture.Destroy()
	}

	if isTile {
		g.ResourceRegistry[id] = &TileResource{texture: texture}
	} else {
		g.ResourceRegistry[id] = &NonTileResource{texture: texture, w: width, h: height}
	}
}

func (g *Graphic) registerFlippedResourceFromSurface(
	surface *sdl.Surface,
	id ResourceID,
	width,
	height int32,
	flipHorizontal bool,
	isTile bool) {

	if isTile && (width != TILE_SIZE || height != TILE_SIZE) {
		log.Fatalf("declared to be tile but has wrong width (%d) or height (%d)", width, height)
	}

	texture, err := g.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatal(err)
	}

	// make sure the tile is in good shape
	if surface.W != width || surface.H != height {
		oldTexture := texture
		texture, err = g.clipTexture(oldTexture, &sdl.Rect{0, 0, width, height})
		if err != nil {
			log.Fatal(err)
		}
		// release original texture
		oldTexture.Destroy()
	}

	// flip texture
	oldTexture := texture
	texture, err = g.flipTexture(texture, width, height, flipHorizontal)
	oldTexture.Destroy()

	if isTile {
		g.ResourceRegistry[id] = &TileResource{texture: texture}
	} else {
		g.ResourceRegistry[id] = &NonTileResource{texture: texture, w: width, h: height}
	}
}

// RenderResource renders a tile (or a part of tile specified by srcRect) to a given position in screen
func (g *Graphic) RenderResource(resource Resource, srcRect *sdl.Rect, dstRect *sdl.Rect) {
	g.renderer.Copy(resource.GetTexture(), srcRect, dstRect)
}

func (g *Graphic) loadAllResources() {
	// -------------------------------
	// load tile resources
	// -------------------------------

	g.registerTileResource("assets/brick.png", RESOURCE_TYPE_BRICK)

	g.registerTileResource("assets/ground-grass-left.png", RESOURCE_TYPE_GROUD_GRASS_LEFT)
	g.registerTileResource("assets/ground-grass-mid.png", RESOURCE_TYPE_GROUD_GRASS_MID)
	g.registerTileResource("assets/ground-grass-right.png", RESOURCE_TYPE_GROUD_GRASS_RIGHT)

	g.registerTileResource("assets/ground-inner-mid.png", RESOURCE_TYPE_GROUD_INNER_MID)

	g.registerTileResource("assets/myth-box-normal.png", RESOURCE_TYPE_MYTH_BOX_NORMAL)
	g.registerTileResource("assets/myth-box-normal-light.png", RESOURCE_TYPE_MYTH_BOX_NORMAL_LIGHT)
	g.registerTileResource("assets/myth-box-empty.png", RESOURCE_TYPE_MYTH_BOX_EMPTY)

	g.registerTileResource("assets/coin.png", RESOURCE_TYPE_COIN)

	// -------------------------------
	// Load non-tile resources
	// -------------------------------

	// hero
	g.registerNonTailResource("assets/hero-stand.png", RESOURCE_TYPE_HERO_STAND_RIGHT)
	g.registerNonTailResource("assets/hero-walking.png", RESOURCE_TYPE_HERO_WALKING_RIGHT)
	g.registerFlippedNonTailResource("assets/hero-stand.png", RESOURCE_TYPE_HERO_STAND_LEFT, true)
	g.registerFlippedNonTailResource("assets/hero-walking.png", RESOURCE_TYPE_HERO_WALKING_LEFT, true)
}

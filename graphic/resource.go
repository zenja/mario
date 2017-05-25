package graphic

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

type ResourceID int

type Resource interface {
	GetTexture() *sdl.Texture
	GetW() int32
	GetH() int32
}

// Resource IDs
const (
	RESOURCE_TYPE_GROUD = iota
	RESOURCE_TYPE_HERO
)

const TILE_SIZE = 64

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
// Graphic methods relative to resource
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (g *Graphic) GetResource(resourceID ResourceID) Resource {
	r, ok := g.ResourceRegistry[resourceID]
	if !ok {
		log.Fatalf("ResourceID %d not found in resource registry", resourceID)
	}
	return r
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

// RenderResource renders a tile (or a part of tile specified by srcRect) to a given position in screen
func (g *Graphic) RenderResource(resource Resource, srcRect *sdl.Rect, dstRect *sdl.Rect) {
	g.renderer.Copy(resource.GetTexture(), srcRect, dstRect)
}

func (g *Graphic) loadAllResources() {
	// load tile resources
	g.registerTileResource("assets/texture.png", RESOURCE_TYPE_GROUD)
	// load non-tile resources
	g.registerNonTailResource("assets/hero.png", RESOURCE_TYPE_HERO)
}

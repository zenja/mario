package graphic

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

type TileID int

const (
	// Tile IDs
	TILE_TYPE_GROUD = iota
)

const TILE_SIZE = 64

var TILE_RECT = sdl.Rect{X: 0, Y: 0, W: TILE_SIZE, H: TILE_SIZE}

type Tile struct {
	texture *sdl.Texture
}

// registerTile loads a sprite into a Tile from a file
func (g *Graphic) registerTile(filename string, id TileID) {
	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	g.registerTileFromSurface(surface, id)
}

// registerTileFromSurface loads a sprite into a Tile from a surface
// User need to free the surface himself
func (g *Graphic) registerTileFromSurface(surface *sdl.Surface, id TileID) {
	texture, err := g.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatal(err)
	}

	// make sure the tile is in good shape
	if surface.W != TILE_SIZE || surface.H != TILE_SIZE {
		oldTexture := texture
		texture, err = g.clipTexture(oldTexture, &TILE_RECT)
		if err != nil {
			log.Fatal(err)
		}
		// release original texture
		oldTexture.Destroy()
	}

	g.ResourceRegistry[id] = &Tile{texture: texture}
}

// RenderTile renders a tile (or a part of tile specified by srcRect) to a given position in screen
func (g *Graphic) RenderTile(id TileID, srcRect *sdl.Rect, dstRect *sdl.Rect) {
	tile, ok := g.ResourceRegistry[id]
	if !ok {
		log.Fatal(fmt.Errorf("resource not found: %d", id))
	}
	g.renderer.Copy(tile.texture, srcRect, dstRect)
}

func (g *Graphic) loadAllTiles() {
	g.registerTile("assets/texture.png", TILE_TYPE_GROUD)
}

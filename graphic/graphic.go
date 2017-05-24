package graphic

import (
	"log"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

const (
	SCREEN_WIDTH  = 1024
	SCREEN_HEIGHT = 768
)

type Graphic struct {
	window   *sdl.Window
	renderer *sdl.Renderer

	ResourceRegistry map[TileID]*Tile
}

func New() *Graphic {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatal(err)
	}

	// Create window
	window, err := sdl.CreateWindow("Mario", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal(err)
	}

	// Create renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		log.Fatal(err)
	}

	// Init font system
	err = ttf.Init()
	if err != nil {
		log.Fatal(err)
	}


	// Load tiles
	g := &Graphic{window: window, renderer: renderer, ResourceRegistry: make(map[TileID]*Tile)}
	g.loadAllTiles()

	return g
}

func (g *Graphic) DestroyAndQuit() {
	g.renderer.Destroy()
	g.window.Destroy()

	ttf.Quit()
	img.Quit()
	sdl.Quit()
}

// Show the screen
func (g *Graphic) ClearScreen() {
	g.renderer.Clear()
}

// Show the screen
func (g *Graphic) ShowScreen() {
	g.renderer.Present()
}

func (g *Graphic) Delay(ms uint32) {
	sdl.Delay(ms)
}

// clipTexture is a helper function to create a new texture from a region of a texture
// User needs to free the input texture himself if needed
func (g *Graphic) clipTexture(texture *sdl.Texture, rect *sdl.Rect) (*sdl.Texture, error) {
	newTexture, err := g.renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_TARGET, int(rect.W), int(rect.H))
	if err != nil {
		return nil, errors.Wrap(err, "failed to clip texture")
	}
	if err = g.renderer.SetRenderTarget(newTexture); err != nil {
		return nil, errors.Wrap(err, "failed to set render target")
	}
	if err = g.renderer.Copy(texture, nil, rect); err != nil {
		return nil, errors.Wrap(err, "failed to render texture")
	}
	// reset render target
	if err = g.renderer.SetRenderTarget(nil); err != nil {
		return nil, errors.Wrap(err, "failed to reset render target")
	}
	return newTexture, nil
}

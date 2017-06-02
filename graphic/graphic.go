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

const (
	FPS           = 60
	DELAY_TIME_MS = 1000 / FPS
)

type Graphic struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	font     *ttf.Font

	ResourceRegistry map[ResourceID]Resource
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

	// Load font
	font, err := ttf.OpenFont("assets/fonts/Menlo-Regular.ttf", 18)
	if err != nil {
		log.Fatal(err)
	}

	// Load resources
	g := &Graphic{
		window:           window,
		renderer:         renderer,
		font:             font,
		ResourceRegistry: make(map[ResourceID]Resource),
	}
	g.loadAllResources()

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
func (g *Graphic) ClearScreenWithColor(color sdl.Color) {
	var err error
	err = g.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	if err != nil {
		log.Fatal("failed to set renderer draw color")
	}
	err = g.renderer.Clear()
	if err != nil {
		log.Fatal("failed to clear renderer")
	}
	// reset draw color
	err = g.renderer.SetDrawColor(0, 0, 0, 0)
	if err != nil {
		log.Fatal("failed to reset renderer draw color")
	}
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
	newTexture, err := g.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, int(rect.W), int(rect.H))
	if err != nil {
		return nil, errors.Wrap(err, "failed to clip texture")
	}

	// will make pixels with alpha 0 fully transparent
	if err = newTexture.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		return nil, errors.Wrap(err, "failed to set blend mode")
	}

	if err = g.renderer.SetRenderTarget(newTexture); err != nil {
		return nil, errors.Wrap(err, "failed to set render target")
	}

	// this together with blend mode will make transparent area
	if err = g.renderer.SetDrawColor(0, 0, 0, 0); err != nil {
		return nil, errors.Wrap(err, "failed to reset draw color")
	}

	if err = g.renderer.Clear(); err != nil {
		return nil, errors.Wrap(err, "failed to clear renderer")
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

// flipTexture is a helper function to create a flipped texture from a region of a texture
// User needs to free the input texture himself if needed
func (g *Graphic) flipTexture(texture *sdl.Texture, width int32, height int32, flipHorizontal bool) (*sdl.Texture, error) {
	newTexture, err := g.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, int(width), int(height))
	if err != nil {
		return nil, errors.Wrap(err, "failed to clip texture")
	}

	// will make pixels with alpha 0 fully transparent
	if err = newTexture.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		return nil, errors.Wrap(err, "failed to set blend mode")
	}

	if err = g.renderer.SetRenderTarget(newTexture); err != nil {
		return nil, errors.Wrap(err, "failed to set render target")
	}

	// this together with blend mode will make transparent area
	if err = g.renderer.SetDrawColor(0, 0, 0, 0); err != nil {
		return nil, errors.Wrap(err, "failed to reset draw color")
	}

	if err := g.renderer.Clear(); err != nil {
		return nil, errors.Wrap(err, "failed to clear renderer")
	}

	var flipFlag sdl.RendererFlip
	if flipHorizontal {
		flipFlag = sdl.FLIP_HORIZONTAL
	} else {
		flipFlag = sdl.FLIP_VERTICAL
	}
	if err := g.renderer.CopyEx(texture, nil, nil, 0, nil, flipFlag); err != nil {
		return nil, errors.Wrap(err, "failed to render texture")
	}

	// reset render target
	if err = g.renderer.SetRenderTarget(nil); err != nil {
		return nil, errors.Wrap(err, "failed to reset render target")
	}
	return newTexture, nil
}

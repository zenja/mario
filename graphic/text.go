package graphic

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/vector"
)

func (g *Graphic) DrawText(text string, pos vector.Pos, color sdl.Color) {
	surface, err := g.font.RenderUTF8_Solid(text, color)
	if err != nil {
		log.Fatal(err)
	}

	texture, err := g.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatal(err)
	}

	width := surface.W
	height := surface.H

	// Free loaded surface
	surface.Free()

	g.renderer.Copy(texture, nil, &sdl.Rect{pos.X, pos.Y, width, height})
}

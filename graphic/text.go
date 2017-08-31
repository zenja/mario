package graphic

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/vector"
)

func DrawText(text string, pos vector.Pos, color sdl.Color) {
	surface, err := font.RenderUTF8_Solid(text, color)
	if err != nil {
		log.Fatal(err)
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatal(err)
	}
	defer texture.Destroy()

	width := surface.W
	height := surface.H

	// Free loaded surface
	surface.Free()

	renderer.Copy(texture, nil, &sdl.Rect{pos.X, pos.Y, width, height})
}

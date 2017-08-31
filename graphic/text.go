package graphic

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/vector"
)

func DrawTextAbsolute(text string, pos vector.Pos, color sdl.Color) {
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

func DrawTextRelative(text string, pos vector.Pos, camPos vector.Pos, color sdl.Color) {
	levelRect := sdl.Rect{
		pos.X,
		pos.Y,
		1,
		1,
	}
	_, rectInCamera := VisibleRectInCamera(levelRect, camPos.X, camPos.Y)

	if rectInCamera == nil {
		return
	}

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

	renderer.Copy(texture, nil, &sdl.Rect{rectInCamera.X, rectInCamera.Y, width, height})
}

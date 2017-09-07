package graphic

import (
	"log"

	"fmt"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

func getFacedResource(
	filename string,
	faceID string,
	mainWidth, mainHeight int32,
	faceWidth, faceHeight int32,
	faceX, faceY int32,
	faceAngle float64,
	flipHorizontal bool,
	flipVertical bool) Resource {

	mainSurface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer mainSurface.Free()

	mainTexture, err := loadTextureFromSurface(mainSurface, mainWidth, mainHeight, 0, flipHorizontal, flipVertical)
	if err != nil {
		log.Fatal(err)
	}
	defer mainTexture.Destroy()

	faceSurface, err := img.Load(fmt.Sprintf("assets/faces/%s.png", faceID))
	if err != nil {
		log.Fatal(err)
	}
	defer faceSurface.Free()

	faceTexture, err := loadTextureFromSurface(faceSurface, faceWidth, faceHeight, 0, flipHorizontal, flipVertical)
	if err != nil {
		log.Fatal(err)
	}
	defer faceTexture.Destroy()

	combined, err := combineTexture(
		mainTexture, faceTexture, mainWidth, mainHeight, faceWidth, faceHeight, faceX, faceY, faceAngle)
	if err != nil {
		log.Fatal(err)
	}

	return &BasicResource{texture: combined, w: mainWidth, h: mainHeight}
}

func registerFacedResource(
	filename string,
	faceID string,
	id ResourceID,
	mainWidth, mainHeight int32,
	faceWidth, faceHeight int32,
	faceX, faceY int32,
	faceAngle float64,
	flipHorizontal bool,
	flipVertical bool) {

	res := getFacedResource(
		filename, faceID, mainWidth, mainHeight,
		faceWidth, faceHeight, faceX, faceY, faceAngle, flipHorizontal, flipVertical)
	resourceRegistry[id] = res
}

func combineTexture(
	backTexture *sdl.Texture,
	frontTexture *sdl.Texture,
	backWidth, backHeight int32,
	frontWidth, frontHeight int32,
	frontX, frontY int32,
	frontAngle float64) (*sdl.Texture, error) {

	newTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, int(backWidth), int(backHeight))
	if err != nil {
		return nil, errors.Wrap(err, "failed to clip texture")
	}

	// will make pixels with alpha 0 fully transparent
	if err = newTexture.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		return nil, errors.Wrap(err, "failed to set blend mode")
	}

	if err = renderer.SetRenderTarget(newTexture); err != nil {
		return nil, errors.Wrap(err, "failed to set render target")
	}

	// this together with blend mode will make transparent area
	if err = renderer.SetDrawColor(0, 0, 0, 0); err != nil {
		return nil, errors.Wrap(err, "failed to reset draw color")
	}

	if err = renderer.Clear(); err != nil {
		return nil, errors.Wrap(err, "failed to clear renderer")
	}

	if err = renderer.Copy(backTexture, nil, nil); err != nil {
		return nil, errors.Wrap(err, "failed to render texture")
	}

	if err = renderer.Copy(frontTexture, nil, &sdl.Rect{frontX, frontY, frontWidth, frontHeight}); err != nil {
		return nil, errors.Wrap(err, "failed to render texture")
	}

	renderer.CopyEx(frontTexture, nil, &sdl.Rect{frontX, frontY, frontWidth, frontHeight}, frontAngle, nil, sdl.FLIP_NONE)

	// reset render target
	if err = renderer.SetRenderTarget(nil); err != nil {
		return nil, errors.Wrap(err, "failed to reset render target")
	}

	return newTexture, nil
}

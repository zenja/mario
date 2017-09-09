package graphic

import (
	"log"

	"github.com/pkg/errors"
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
	SetResourceAlpha(alpha uint8)
}

var resourceRegistry map[ResourceID]Resource = make(map[ResourceID]Resource)

// Resource IDs
const (
	RESOURCE_TYPE_CURR_BG ResourceID = iota

	RESOURCE_TYPE_BRICK_RED
	RESOURCE_TYPE_BRICK_PIECE_RED
	RESOURCE_TYPE_BRICK_YELLOW
	RESOURCE_TYPE_BRICK_PIECE_YELLOW

	RESOURCE_TYPE_GRASS_GROUD_LEFT
	RESOURCE_TYPE_GRASS_GROUD_MID
	RESOURCE_TYPE_GRASS_GROUD_RIGHT

	RESOURCE_TYPE_GROUD_LEFT
	RESOURCE_TYPE_GROUD_MID
	RESOURCE_TYPE_GROUD_RIGHT

	RESOURCE_TYPE_WATER_0
	RESOURCE_TYPE_WATER_1
	RESOURCE_TYPE_WATER_2
	RESOURCE_TYPE_WATER_3
	RESOURCE_TYPE_WATER_4
	RESOURCE_TYPE_WATER_5
	RESOURCE_TYPE_WATER_6
	RESOURCE_TYPE_WATER_FULL

	RESOURCE_TYPE_DEC_GRASS_0
	RESOURCE_TYPE_DEC_GRASS_1

	RESOURCE_TYPE_DEC_TREE_0

	RESOURCE_TYPE_DEC_FAT_TREE_GREEN
	RESOURCE_TYPE_DEC_FAT_TREE_RED
	RESOURCE_TYPE_DEC_FAT_TREE_PINK
	RESOURCE_TYPE_DEC_FAT_TREE_WHITE

	RESOURCE_TYPE_DEC_PRINCESS_0
	RESOURCE_TYPE_DEC_PRINCESS_1

	RESOURCE_TYPE_DEC_PRINCESS_IS_WAITING

	RESOURCE_TYPE_DEC_SUPER_MARIO_PAYPAL

	RESOURCE_TYPE_PAYPAL_IS_NEW_MONEY

	RESOURCE_TYPE_DEC_HIGH_ENERGY_AHEAD

	RESOURCE_TYPE_MYTH_BOX_NORMAL
	RESOURCE_TYPE_MYTH_BOX_NORMAL_LIGHT
	RESOURCE_TYPE_MYTH_BOX_EMPTY

	RESOURCE_TYPE_PIPE_LEFT_TOP
	RESOURCE_TYPE_PIPE_RIGHT_TOP
	RESOURCE_TYPE_PIPE_LEFT_MID
	RESOURCE_TYPE_PIPE_RIGHT_MID
	RESOURCE_TYPE_PIPE_LEFT_BOTTOM
	RESOURCE_TYPE_PIPE_RIGHT_BOTTOM

	RESOURCE_TYPE_COIN_0
	RESOURCE_TYPE_COIN_1
	RESOURCE_TYPE_COIN_2
	RESOURCE_TYPE_COIN_3

	RESOURCE_TYPE_GOOD_MUSHROOM

	RESOURCE_TYPE_MUSHROOM_ENEMY_0
	RESOURCE_TYPE_MUSHROOM_ENEMY_1
	RESOURCE_TYPE_MUSHROOM_ENEMY_HIT
	RESOURCE_TYPE_MUSHROOM_ENEMY_DOWN

	RESOURCE_TYPE_TORTOISE_RED_INSIDE
	RESOURCE_TYPE_TORTOISE_RED_SEMI_INSIDE

	RESOURCE_TYPE_BOSS_A_LEFT_0
	RESOURCE_TYPE_BOSS_A_LEFT_1
	RESOURCE_TYPE_BOSS_A_RIGHT_0
	RESOURCE_TYPE_BOSS_A_RIGHT_1

	RESOURCE_TYPE_BOSS_C_LEFT_0
	RESOURCE_TYPE_BOSS_C_LEFT_1
	RESOURCE_TYPE_BOSS_C_RIGHT_0
	RESOURCE_TYPE_BOSS_C_RIGHT_1

	RESOURCE_TYPE_BOSS_D_LEFT_0
	RESOURCE_TYPE_BOSS_D_LEFT_1
	RESOURCE_TYPE_BOSS_D_RIGHT_0
	RESOURCE_TYPE_BOSS_D_RIGHT_1

	RESOURCE_TYPE_BOSS_E_LEFT_0
	RESOURCE_TYPE_BOSS_E_LEFT_1
	RESOURCE_TYPE_BOSS_E_RIGHT_0
	RESOURCE_TYPE_BOSS_E_RIGHT_1

	RESOURCE_TYPE_BOSS_F_LEFT_0
	RESOURCE_TYPE_BOSS_F_LEFT_1
	RESOURCE_TYPE_BOSS_F_RIGHT_0
	RESOURCE_TYPE_BOSS_F_RIGHT_1

	RESOURCE_TYPE_BOSS_G_LEFT_0
	RESOURCE_TYPE_BOSS_G_LEFT_1
	RESOURCE_TYPE_BOSS_G_RIGHT_0
	RESOURCE_TYPE_BOSS_G_RIGHT_1

	RESOURCE_TYPE_BOSS_H_LEFT_0
	RESOURCE_TYPE_BOSS_H_LEFT_1
	RESOURCE_TYPE_BOSS_H_RIGHT_0
	RESOURCE_TYPE_BOSS_H_RIGHT_1

	RESOURCE_TYPE_BOSS_BOOM

	RESOURCE_TYPE_BANG

	RESOURCE_TYPE_FIREBALL_0
	RESOURCE_TYPE_FIREBALL_1
	RESOURCE_TYPE_FIREBALL_2
	RESOURCE_TYPE_FIREBALL_3
	RESOURCE_TYPE_FIREBALL_BOOM

	RESOURCE_TYPE_SHIT_0
	RESOURCE_TYPE_SHIT_1
	RESOURCE_TYPE_SHIT_2
	RESOURCE_TYPE_SHIT_3
	RESOURCE_TYPE_SHIT_BOOM

	RESOURCE_TYPE_BUG_0
	RESOURCE_TYPE_BUG_1
	RESOURCE_TYPE_BUG_2
	RESOURCE_TYPE_BUG_3
	RESOURCE_TYPE_BUG_BOOM

	RESOURCE_TYPE_SHINE_0
	RESOURCE_TYPE_SHINE_1
	RESOURCE_TYPE_SHINE_2

	RESOURCE_TYPE_SWORD_0
	RESOURCE_TYPE_SWORD_1
	RESOURCE_TYPE_SWORD_2
	RESOURCE_TYPE_SWORD_3

	RESOURCE_TYPE_ICEBALL_0
	RESOURCE_TYPE_ICEBALL_1
	RESOURCE_TYPE_ICEBALL_2
	RESOURCE_TYPE_ICEBALL_3

	RESOURCE_TYPE_APPLE_0
	RESOURCE_TYPE_APPLE_1
	RESOURCE_TYPE_APPLE_2
	RESOURCE_TYPE_APPLE_3

	RESOURCE_TYPE_CHERRY_0
	RESOURCE_TYPE_CHERRY_1
	RESOURCE_TYPE_CHERRY_2
	RESOURCE_TYPE_CHERRY_3

	RESOURCE_TYPE_MOON_0
	RESOURCE_TYPE_MOON_1
	RESOURCE_TYPE_MOON_2
	RESOURCE_TYPE_MOON_3

	RESOURCE_TYPE_AXE_0
	RESOURCE_TYPE_AXE_1
	RESOURCE_TYPE_AXE_2
	RESOURCE_TYPE_AXE_3

	RESOURCE_TYPE_SKULL_0
	RESOURCE_TYPE_SKULL_1
	RESOURCE_TYPE_SKULL_2
	RESOURCE_TYPE_SKULL_3

	RESOURCE_TYPE_UPGRADE_FLOWER

	RESOURCE_TYPE_EATER_FLOWER_0
	RESOURCE_TYPE_EATER_FLOWER_1

	RESOURCE_TYPE_BLACK_SCREEN

	RESOURCE_TYPE_HERO_0_STAND_LEFT
	RESOURCE_TYPE_HERO_0_WALKING_LEFT
	RESOURCE_TYPE_HERO_0_JUMP_LEFT

	RESOURCE_TYPE_HERO_0_STAND_RIGHT
	RESOURCE_TYPE_HERO_0_WALKING_RIGHT
	RESOURCE_TYPE_HERO_0_JUMP_RIGHT

	RESOURCE_TYPE_HERO_1_STAND_LEFT
	RESOURCE_TYPE_HERO_1_WALKING_LEFT
	RESOURCE_TYPE_HERO_1_JUMP_LEFT

	RESOURCE_TYPE_HERO_1_STAND_RIGHT
	RESOURCE_TYPE_HERO_1_WALKING_RIGHT
	RESOURCE_TYPE_HERO_1_JUMP_RIGHT

	RESOURCE_TYPE_HERO_2_STAND_LEFT
	RESOURCE_TYPE_HERO_2_WALKING_LEFT
	RESOURCE_TYPE_HERO_2_JUMP_LEFT

	RESOURCE_TYPE_HERO_2_STAND_RIGHT
	RESOURCE_TYPE_HERO_2_WALKING_RIGHT
	RESOURCE_TYPE_HERO_2_JUMP_RIGHT
)

const TILE_SIZE = 50

const (
	hero_0_width  = 50
	hero_0_height = 75
	hero_1_width  = 55
	hero_1_height = 93
	hero_2_width  = 55
	hero_2_height = 93

	tortoise_walking_width  = 50
	tortoise_walking_height = 65
	tortoise_inside_width   = 50
	tortoise_inside_height  = 43
)

func Res(id ResourceID) Resource {
	return resourceRegistry[id]
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// BasicResource
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type BasicResource struct {
	texture *sdl.Texture
	w, h    int32
}

func (br *BasicResource) GetTexture() *sdl.Texture {
	return br.texture
}

func (br *BasicResource) GetW() int32 {
	return br.w
}

func (br *BasicResource) GetH() int32 {
	return br.h
}

func (br *BasicResource) SetResourceAlpha(alpha uint8) {
	br.GetTexture().SetAlphaMod(alpha)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Public helper functions
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
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Other public utils
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func DrawRect(rect sdl.Rect, camPos vector.Pos, r, g, b, a uint8) {
	originR, originG, originB, originA, err := renderer.GetDrawColor()
	if err != nil {
		log.Fatalf("failed to get draw color: %s", err)
	}
	_, rectInCam := VisibleRectInCamera(rect, camPos.X, camPos.Y)
	if rectInCam != nil {
		renderer.SetDrawColor(r, g, b, a)
		renderer.DrawRect(rectInCam)
		renderer.SetDrawColor(originR, originG, originB, originA)
	}
}

func FillRect(rect sdl.Rect, camPos vector.Pos, r, g, b, a uint8) {
	originR, originG, originB, originA, err := renderer.GetDrawColor()
	if err != nil {
		log.Fatalf("failed to get draw color: %s", err)
	}
	_, rectInCam := VisibleRectInCamera(rect, camPos.X, camPos.Y)
	if rectInCam != nil {
		renderer.SetDrawColor(r, g, b, a)
		renderer.FillRect(rectInCam)
		renderer.SetDrawColor(originR, originG, originB, originA)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Basic resource related functions
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func GetResource(resourceID ResourceID) Resource {
	r, ok := resourceRegistry[resourceID]
	if !ok {
		log.Fatalf("ResourceID %d not found in resource registry", resourceID)
	}
	return r
}

// drawResource is a helper function to draw a resource on level to camera
func DrawResource(resource Resource, levelRect sdl.Rect, camPos vector.Pos) {
	rectInResource, rectInCamera := VisibleRectInCamera(levelRect, camPos.X, camPos.Y)
	if rectInResource != nil {
		RenderResource(resource, rectInResource, rectInCamera)
	}
}

func registerResource(texture *sdl.Texture, width, height int32, id ResourceID) {
	resourceRegistry[id] = &BasicResource{texture: texture, w: width, h: height}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Register resource from file
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// registerResourceFromFile loads a sprite into a BasicResource from a file
func registerResourceFromFile(filename string, id ResourceID) {
	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	registerResourceFromSurface(surface, id, surface.W, surface.H)
}

func registerTileResourceFromFile(filename string, id ResourceID) {
	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	registerResourceFromSurface(surface, id, TILE_SIZE, TILE_SIZE)
}

func registerScaledResourceFromFile(filename string, id ResourceID, dstWidth int32, dstHeight int32) {
	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	registerScaledResourceFromSurface(surface, id, dstWidth, dstHeight)
}

// RegisterBackgroundResource register a level background resource, scale it to have level's height
// This function has to be public because it is used when parsing a level
func RegisterBackgroundResource(filename string, id ResourceID, tilesInY int) {
	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	dstHeight := int32(tilesInY * TILE_SIZE)
	dstWidth := surface.W * (dstHeight / surface.H)

	registerScaledResourceFromSurface(surface, id, dstWidth, dstHeight)
}

func registerResourceFromFileEx(
	filename string,
	id ResourceID,
	width,
	height int32,
	angle float64,
	flipHorizontal bool,
	flipVertical bool) {

	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	registerResourceFromSurfaceEx(surface, id, width, height, angle, flipHorizontal, flipVertical)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Register resource from surface
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// registerResourceFromSurface loads a sprite from a surface into a Resource object
// User need to free the surface himself
func registerResourceFromSurface(surface *sdl.Surface, id ResourceID, width, height int32) {
	registerResourceFromSurfaceEx(surface, id, width, height, 0, false, false)
}

func registerScaledResourceFromSurface(surface *sdl.Surface, id ResourceID, dstWidth, dstHeight int32) {
	registerResourceFromSurface(surface, id, dstWidth, dstHeight)
}

func registerFlippedResourceFromSurface(
	surface *sdl.Surface,
	id ResourceID,
	width,
	height int32,
	flipHorizontal bool) {

	registerResourceFromSurfaceEx(surface, id, width, height, 0, flipHorizontal, !flipHorizontal)
}

func registerResourceFromSurfaceEx(
	surface *sdl.Surface,
	id ResourceID,
	width,
	height int32,
	angle float64,
	flipHorizontal bool,
	flipVertical bool) {

	texture, err := loadTextureFromSurface(surface, width, height, angle, flipHorizontal, flipVertical)
	if err != nil {
		log.Fatal(err)
	}
	registerResource(texture, width, height, id)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Graphic functions relative to texture
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func loadSimpleTextureFromFile(filename string) (*sdl.Texture, error) {
	surface, err := img.Load(filename)
	if err != nil {
		return nil, err
	}
	defer surface.Free()

	return loadTextureFromSurface(surface, surface.W, surface.H, 0, false, false)
}

func loadTextureFromFile(
	filename string,
	width, height int32,
	angle float64,
	flipHorizontal bool,
	flipVertical bool) (*sdl.Texture, error) {

	surface, err := img.Load(filename)
	if err != nil {
		return nil, err
	}
	defer surface.Free()

	return loadTextureFromSurface(surface, width, height, angle, flipHorizontal, flipVertical)
}

func loadTextureFromSurface(
	surface *sdl.Surface,
	width, height int32,
	angle float64,
	flipHorizontal bool,
	flipVertical bool) (*sdl.Texture, error) {

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}
	defer texture.Destroy()

	newTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, int(width), int(height))
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

	var flipMode sdl.RendererFlip
	if flipHorizontal {
		flipMode = sdl.FLIP_HORIZONTAL
	} else if flipVertical {
		flipMode = sdl.FLIP_VERTICAL
	} else {
		flipMode = sdl.FLIP_NONE
	}
	renderer.CopyEx(texture, nil, &sdl.Rect{0, 0, width, height}, angle, nil, flipMode)

	// reset render target
	if err = renderer.SetRenderTarget(nil); err != nil {
		return nil, errors.Wrap(err, "failed to reset render target")
	}

	return newTexture, nil
}

// RenderResource renders a tile (or a part of tile specified by srcRect) to a given position in screen
func RenderResource(resource Resource, srcRect *sdl.Rect, dstRect *sdl.Rect) {
	renderer.Copy(resource.GetTexture(), srcRect, dstRect)
}

func LoadAllResources(heroUserID string) {
	// -------------------------------
	// load tile resources
	// -------------------------------

	// brick
	registerTileResourceFromFile("assets/brick-red.png", RESOURCE_TYPE_BRICK_RED)
	registerTileResourceFromFile("assets/brick-yellow.png", RESOURCE_TYPE_BRICK_YELLOW)

	// grass
	registerTileResourceFromFile("assets/grass-ground-left.png", RESOURCE_TYPE_GRASS_GROUD_LEFT)
	registerTileResourceFromFile("assets/grass-ground-mid.png", RESOURCE_TYPE_GRASS_GROUD_MID)
	registerTileResourceFromFile("assets/grass-ground-right.png", RESOURCE_TYPE_GRASS_GROUD_RIGHT)

	// ground
	registerTileResourceFromFile("assets/ground-left.png", RESOURCE_TYPE_GROUD_LEFT)
	registerTileResourceFromFile("assets/ground-mid.png", RESOURCE_TYPE_GROUD_MID)
	registerTileResourceFromFile("assets/ground-right.png", RESOURCE_TYPE_GROUD_RIGHT)

	// myth box
	registerTileResourceFromFile("assets/myth-box-normal.png", RESOURCE_TYPE_MYTH_BOX_NORMAL)
	registerTileResourceFromFile("assets/myth-box-normal-light.png", RESOURCE_TYPE_MYTH_BOX_NORMAL_LIGHT)
	registerTileResourceFromFile("assets/myth-box-empty.png", RESOURCE_TYPE_MYTH_BOX_EMPTY)

	// pipe
	registerTileResourceFromFile("assets/pipe-left-top.png", RESOURCE_TYPE_PIPE_LEFT_TOP)
	registerTileResourceFromFile("assets/pipe-right-top.png", RESOURCE_TYPE_PIPE_RIGHT_TOP)
	registerTileResourceFromFile("assets/pipe-left-mid.png", RESOURCE_TYPE_PIPE_LEFT_MID)
	registerTileResourceFromFile("assets/pipe-right-mid.png", RESOURCE_TYPE_PIPE_RIGHT_MID)
	registerTileResourceFromFile("assets/pipe-left-bottom.png", RESOURCE_TYPE_PIPE_LEFT_BOTTOM)
	registerTileResourceFromFile("assets/pipe-right-bottom.png", RESOURCE_TYPE_PIPE_RIGHT_BOTTOM)

	// coin
	registerTileResourceFromFile("assets/coin-0.png", RESOURCE_TYPE_COIN_0)
	registerTileResourceFromFile("assets/coin-1.png", RESOURCE_TYPE_COIN_1)
	registerTileResourceFromFile("assets/coin-2.png", RESOURCE_TYPE_COIN_2)
	registerTileResourceFromFile("assets/coin-3.png", RESOURCE_TYPE_COIN_3)

	// good mushroom
	registerTileResourceFromFile("assets/mushroom.png", RESOURCE_TYPE_GOOD_MUSHROOM)

	// upgrade flower
	registerTileResourceFromFile("assets/upgrade-flower.png", RESOURCE_TYPE_UPGRADE_FLOWER)

	// water
	registerTileResourceFromFile("assets/water-0.png", RESOURCE_TYPE_WATER_0)
	registerTileResourceFromFile("assets/water-1.png", RESOURCE_TYPE_WATER_1)
	registerTileResourceFromFile("assets/water-2.png", RESOURCE_TYPE_WATER_2)
	registerTileResourceFromFile("assets/water-3.png", RESOURCE_TYPE_WATER_3)
	registerTileResourceFromFile("assets/water-4.png", RESOURCE_TYPE_WATER_4)
	registerTileResourceFromFile("assets/water-5.png", RESOURCE_TYPE_WATER_5)
	registerTileResourceFromFile("assets/water-6.png", RESOURCE_TYPE_WATER_6)
	registerTileResourceFromFile("assets/water-pixel.png", RESOURCE_TYPE_WATER_FULL)

	// -------------------------------
	// Load non-tile resources
	// -------------------------------

	// hero 0
	registerFacedResource("assets/hero-0-stand.png", heroUserID, RESOURCE_TYPE_HERO_0_STAND_RIGHT,
		hero_0_width, hero_0_height, 40, 50, 5, 0, 5, false, false)
	registerFacedResource("assets/hero-0-walking.png", heroUserID, RESOURCE_TYPE_HERO_0_WALKING_RIGHT,
		hero_0_width, hero_0_height, 40, 50, 5, 0, -5, false, false)
	registerFacedResource("assets/hero-0-jump.png", heroUserID, RESOURCE_TYPE_HERO_0_JUMP_RIGHT,
		hero_0_width, hero_0_height, 40, 50, 5, 0, 0, false, false)
	registerFacedResource("assets/hero-0-stand.png", heroUserID, RESOURCE_TYPE_HERO_0_STAND_LEFT,
		hero_0_width, hero_0_height, 40, 50, 5, 0, 5, true, false)
	registerFacedResource("assets/hero-0-walking.png", heroUserID, RESOURCE_TYPE_HERO_0_WALKING_LEFT,
		hero_0_width, hero_0_height, 40, 50, 5, 0, -5, true, false)
	registerFacedResource("assets/hero-0-jump.png", heroUserID, RESOURCE_TYPE_HERO_0_JUMP_LEFT,
		hero_0_width, hero_0_height, 40, 50, 5, 0, 0, true, false)

	// hero 1
	registerFacedResource("assets/hero-1-stand.png", heroUserID, RESOURCE_TYPE_HERO_1_STAND_RIGHT,
		hero_1_width, hero_1_height, 48, 60, 5, 0, 5, false, false)
	registerFacedResource("assets/hero-1-walking.png", heroUserID, RESOURCE_TYPE_HERO_1_WALKING_RIGHT,
		hero_1_width, hero_1_height, 48, 60, 5, 0, -5, false, false)
	registerFacedResource("assets/hero-1-jump.png", heroUserID, RESOURCE_TYPE_HERO_1_JUMP_RIGHT,
		hero_1_width, hero_1_height, 48, 60, 5, 0, 0, false, false)
	registerFacedResource("assets/hero-1-stand.png", heroUserID, RESOURCE_TYPE_HERO_1_STAND_LEFT,
		hero_1_width, hero_1_height, 48, 60, 5, 0, 5, true, false)
	registerFacedResource("assets/hero-1-walking.png", heroUserID, RESOURCE_TYPE_HERO_1_WALKING_LEFT,
		hero_1_width, hero_1_height, 48, 60, 5, 0, -5, true, false)
	registerFacedResource("assets/hero-1-jump.png", heroUserID, RESOURCE_TYPE_HERO_1_JUMP_LEFT,
		hero_1_width, hero_1_height, 48, 60, 5, 0, 0, true, false)

	// hero 2
	registerFacedResource("assets/hero-2-stand.png", heroUserID, RESOURCE_TYPE_HERO_2_STAND_RIGHT,
		hero_2_width, hero_2_height, 48, 60, 5, 0, 5, false, false)
	registerFacedResource("assets/hero-2-walking.png", heroUserID, RESOURCE_TYPE_HERO_2_WALKING_RIGHT,
		hero_2_width, hero_2_height, 48, 60, 5, 0, -5, false, false)
	registerFacedResource("assets/hero-2-jump.png", heroUserID, RESOURCE_TYPE_HERO_2_JUMP_RIGHT,
		hero_2_width, hero_2_height, 48, 60, 5, 0, 0, false, false)
	registerFacedResource("assets/hero-2-stand.png", heroUserID, RESOURCE_TYPE_HERO_2_STAND_LEFT,
		hero_2_width, hero_2_height, 48, 60, 5, 0, 5, true, false)
	registerFacedResource("assets/hero-2-walking.png", heroUserID, RESOURCE_TYPE_HERO_2_WALKING_LEFT,
		hero_2_width, hero_2_height, 48, 60, 5, 0, -5, true, false)
	registerFacedResource("assets/hero-2-jump.png", heroUserID, RESOURCE_TYPE_HERO_2_JUMP_LEFT,
		hero_2_width, hero_2_height, 48, 60, 5, 0, 0, true, false)

	// decoration: grass
	registerResourceFromFile("assets/dec-grass-0.png", RESOURCE_TYPE_DEC_GRASS_0)
	registerResourceFromFile("assets/dec-grass-1.png", RESOURCE_TYPE_DEC_GRASS_1)

	// decoration: trees
	registerResourceFromFile("assets/dec-tree-0.png", RESOURCE_TYPE_DEC_TREE_0)
	registerScaledResourceFromFile("assets/dec-fat-tree-green.png", RESOURCE_TYPE_DEC_FAT_TREE_GREEN, 200, 220)
	registerScaledResourceFromFile("assets/dec-fat-tree-red.png", RESOURCE_TYPE_DEC_FAT_TREE_RED, 200, 220)
	registerScaledResourceFromFile("assets/dec-fat-tree-pink.png", RESOURCE_TYPE_DEC_FAT_TREE_PINK, 200, 220)
	registerScaledResourceFromFile("assets/dec-fat-tree-white.png", RESOURCE_TYPE_DEC_FAT_TREE_WHITE, 200, 220)

	// decoration: princess
	registerResourceFromFile("assets/princess-0.png", RESOURCE_TYPE_DEC_PRINCESS_0)
	registerResourceFromFile("assets/princess-1.png", RESOURCE_TYPE_DEC_PRINCESS_1)

	// decoration: princess is waiting
	registerScaledResourceFromFile("assets/princess-is-waiting.png", RESOURCE_TYPE_DEC_PRINCESS_IS_WAITING, 150, 200)

	// decoration: super mario paypal
	registerResourceFromFile("assets/super-mario-paypal.png", RESOURCE_TYPE_DEC_SUPER_MARIO_PAYPAL)

	// decoration: paypal is new money
	registerResourceFromFile("assets/paypal-is-new-money.png", RESOURCE_TYPE_PAYPAL_IS_NEW_MONEY)

	// decoration: high energy ahead
	registerScaledResourceFromFile("assets/high-energy-ahead.png", RESOURCE_TYPE_DEC_HIGH_ENERGY_AHEAD, 150, 200)

	// broken pieces
	registerScaledResourceFromFile("assets/brick-piece-red.png", RESOURCE_TYPE_BRICK_PIECE_RED, TILE_SIZE/2, TILE_SIZE/2)
	registerScaledResourceFromFile("assets/brick-piece-yellow.png", RESOURCE_TYPE_BRICK_PIECE_YELLOW, TILE_SIZE/2, TILE_SIZE/2)

	// mushroom enemy
	registerScaledResourceFromFile("assets/mushroom-enemy-0.png", RESOURCE_TYPE_MUSHROOM_ENEMY_0, TILE_SIZE, TILE_SIZE)
	registerScaledResourceFromFile("assets/mushroom-enemy-1.png", RESOURCE_TYPE_MUSHROOM_ENEMY_1, TILE_SIZE, TILE_SIZE)
	registerScaledResourceFromFile("assets/mushroom-enemy-hit.png", RESOURCE_TYPE_MUSHROOM_ENEMY_HIT, TILE_SIZE, TILE_SIZE)
	registerResourceFromFileEx("assets/mushroom-enemy-0.png", RESOURCE_TYPE_MUSHROOM_ENEMY_DOWN, TILE_SIZE, TILE_SIZE, 0, false, true)

	// tortoise enemy res without face
	registerScaledResourceFromFile("assets/tortoise-red-inside.png", RESOURCE_TYPE_TORTOISE_RED_INSIDE, tortoise_inside_width, tortoise_inside_height)
	registerScaledResourceFromFile("assets/tortoise-red-semi-inside.png", RESOURCE_TYPE_TORTOISE_RED_SEMI_INSIDE, tortoise_inside_width, tortoise_inside_height)

	// boss A
	registerFacedResource("assets/boss-a-right-0.png", "yunfeli", RESOURCE_TYPE_BOSS_A_RIGHT_0,
		149, 140, 65, 90, 80, 0, 5, false, false)
	registerFacedResource("assets/boss-a-right-1.png", "yunfeli", RESOURCE_TYPE_BOSS_A_RIGHT_1,
		149, 140, 65, 90, 80, 0, -5, false, false)
	registerFacedResource("assets/boss-a-right-0.png", "yunfeli", RESOURCE_TYPE_BOSS_A_LEFT_0,
		149, 140, 65, 90, 0, 0, 5, true, false)
	registerFacedResource("assets/boss-a-right-1.png", "yunfeli", RESOURCE_TYPE_BOSS_A_LEFT_1,
		149, 140, 65, 90, 0, 0, -5, true, false)

	// boss C
	registerFacedResource("assets/knight-a-left-0.png", "minwu", RESOURCE_TYPE_BOSS_C_RIGHT_0,
		149, 140, 50, 65, 45, 0, 5, true, false)
	registerFacedResource("assets/knight-a-left-1.png", "minwu", RESOURCE_TYPE_BOSS_C_RIGHT_1,
		149, 140, 50, 65, 45, 0, -5, true, false)
	registerFacedResource("assets/knight-a-left-0.png", "minwu", RESOURCE_TYPE_BOSS_C_LEFT_0,
		149, 140, 50, 65, 55, 0, 5, false, false)
	registerFacedResource("assets/knight-a-left-1.png", "minwu", RESOURCE_TYPE_BOSS_C_LEFT_1,
		149, 140, 50, 65, 55, 0, -5, false, false)

	// boss D
	registerFacedResource("assets/soldier-0.png", "huiwang", RESOURCE_TYPE_BOSS_D_RIGHT_0,
		149, 140, 50, 65, 50, 0, 2, false, false)
	registerFacedResource("assets/soldier-1.png", "huiwang", RESOURCE_TYPE_BOSS_D_RIGHT_1,
		149, 140, 50, 65, 50, 0, -2, false, false)
	registerFacedResource("assets/soldier-0.png", "huiwang", RESOURCE_TYPE_BOSS_D_LEFT_0,
		149, 140, 50, 65, 50, 0, 2, true, false)
	registerFacedResource("assets/soldier-1.png", "huiwang", RESOURCE_TYPE_BOSS_D_LEFT_1,
		149, 140, 50, 65, 50, 0, -2, true, false)

	// boss E
	registerFacedResource("assets/ox-right-0.png", "pregev", RESOURCE_TYPE_BOSS_E_RIGHT_0,
		100, 140, 50, 65, 30, 0, 5, false, false)
	registerFacedResource("assets/ox-right-1.png", "pregev", RESOURCE_TYPE_BOSS_E_RIGHT_1,
		100, 140, 50, 65, 30, 0, -5, false, false)
	registerFacedResource("assets/ox-right-0.png", "pregev", RESOURCE_TYPE_BOSS_E_LEFT_0,
		100, 140, 50, 65, 15, 0, 5, true, false)
	registerFacedResource("assets/ox-right-1.png", "pregev", RESOURCE_TYPE_BOSS_E_LEFT_1,
		100, 140, 50, 65, 15, 0, -5, true, false)

	// boss F
	registerFacedResource("assets/fire-sonic-right-0.png", "uarad", RESOURCE_TYPE_BOSS_F_RIGHT_0,
		100, 140, 50, 65, 30, 0, 5, false, false)
	registerFacedResource("assets/fire-sonic-right-1.png", "uarad", RESOURCE_TYPE_BOSS_F_RIGHT_1,
		100, 140, 50, 65, 30, 0, -5, false, false)
	registerFacedResource("assets/fire-sonic-right-0.png", "uarad", RESOURCE_TYPE_BOSS_F_LEFT_0,
		100, 140, 50, 65, 10, 0, 5, true, false)
	registerFacedResource("assets/fire-sonic-right-1.png", "uarad", RESOURCE_TYPE_BOSS_F_LEFT_1,
		100, 140, 50, 65, 10, 0, -5, true, false)

	// boss G
	registerFacedResource("assets/ghost-right-0.png", "gronen", RESOURCE_TYPE_BOSS_G_RIGHT_0,
		120, 140, 50, 65, 53, 0, 5, false, false)
	registerFacedResource("assets/ghost-right-1.png", "gronen", RESOURCE_TYPE_BOSS_G_RIGHT_1,
		120, 140, 50, 65, 53, 0, -5, false, false)
	registerFacedResource("assets/ghost-right-0.png", "gronen", RESOURCE_TYPE_BOSS_G_LEFT_0,
		120, 140, 50, 65, 15, 0, 5, true, false)
	registerFacedResource("assets/ghost-right-1.png", "gronen", RESOURCE_TYPE_BOSS_G_LEFT_1,
		120, 140, 50, 65, 15, 0, -5, true, false)

	// boss H
	registerFacedResource("assets/ghost-right-0.png", "mparnes", RESOURCE_TYPE_BOSS_H_RIGHT_0,
		149, 140, 65, 90, 80, 0, 5, false, false)
	registerFacedResource("assets/ghost-right-1.png", "mparnes", RESOURCE_TYPE_BOSS_H_RIGHT_1,
		149, 140, 65, 90, 80, 0, -5, false, false)
	registerFacedResource("assets/ghost-right-0.png", "mparnes", RESOURCE_TYPE_BOSS_H_LEFT_0,
		149, 140, 65, 90, 0, 0, 5, true, false)
	registerFacedResource("assets/ghost-right-1.png", "mparnes", RESOURCE_TYPE_BOSS_H_LEFT_1,
		149, 140, 65, 90, 0, 0, -5, true, false)

	// boss boom
	registerScaledResourceFromFile("assets/boss-boom.png", RESOURCE_TYPE_BOSS_BOOM, 300, 200)

	// fireball
	registerResourceFromFileEx("assets/fireball.png", RESOURCE_TYPE_FIREBALL_0, 30, 30, 0, false, false)
	registerResourceFromFileEx("assets/fireball.png", RESOURCE_TYPE_FIREBALL_1, 30, 30, 90, false, false)
	registerResourceFromFileEx("assets/fireball.png", RESOURCE_TYPE_FIREBALL_2, 30, 30, 180, false, false)
	registerResourceFromFileEx("assets/fireball.png", RESOURCE_TYPE_FIREBALL_3, 30, 30, 270, false, false)
	registerScaledResourceFromFile("assets/fireball-boom.png", RESOURCE_TYPE_FIREBALL_BOOM, 40, 40)

	// shit
	registerResourceFromFileEx("assets/shit.png", RESOURCE_TYPE_SHIT_0, 30, 30, 0, false, false)
	registerResourceFromFileEx("assets/shit.png", RESOURCE_TYPE_SHIT_1, 30, 30, -90, false, false)
	registerResourceFromFileEx("assets/shit.png", RESOURCE_TYPE_SHIT_2, 30, 30, -180, false, false)
	registerResourceFromFileEx("assets/shit.png", RESOURCE_TYPE_SHIT_3, 30, 30, -270, false, false)
	registerScaledResourceFromFile("assets/shit-boom.png", RESOURCE_TYPE_SHIT_BOOM, 50, 50)

	// bug
	registerResourceFromFileEx("assets/bug.png", RESOURCE_TYPE_BUG_0, 30, 30, 0, false, false)
	registerResourceFromFileEx("assets/bug.png", RESOURCE_TYPE_BUG_1, 30, 30, -90, false, false)
	registerResourceFromFileEx("assets/bug.png", RESOURCE_TYPE_BUG_2, 30, 30, -180, false, false)
	registerResourceFromFileEx("assets/bug.png", RESOURCE_TYPE_BUG_3, 30, 30, -270, false, false)
	registerScaledResourceFromFile("assets/bug-boom.png", RESOURCE_TYPE_BUG_BOOM, 50, 50)

	// eater flower
	registerScaledResourceFromFile("assets/eater-flower-0.png", RESOURCE_TYPE_EATER_FLOWER_0, 52, 75)
	registerScaledResourceFromFile("assets/eater-flower-1.png", RESOURCE_TYPE_EATER_FLOWER_1, 52, 75)

	// shine effect
	registerResourceFromFile("assets/shine-0.png", RESOURCE_TYPE_SHINE_0)
	registerResourceFromFile("assets/shine-1.png", RESOURCE_TYPE_SHINE_1)
	registerResourceFromFile("assets/shine-2.png", RESOURCE_TYPE_SHINE_2)

	// bang
	registerScaledResourceFromFile("assets/bang.png", RESOURCE_TYPE_BANG, 50, 50)

	// black screen
	registerScaledResourceFromFile("assets/black-pixel.png", RESOURCE_TYPE_BLACK_SCREEN, SCREEN_WIDTH, SCREEN_HEIGHT)

	// sword
	registerResourceFromFileEx("assets/sword.png", RESOURCE_TYPE_SWORD_0, 45, 45, 0, false, false)
	registerResourceFromFileEx("assets/sword.png", RESOURCE_TYPE_SWORD_1, 45, 45, 90, false, false)
	registerResourceFromFileEx("assets/sword.png", RESOURCE_TYPE_SWORD_2, 45, 45, 180, false, false)
	registerResourceFromFileEx("assets/sword.png", RESOURCE_TYPE_SWORD_3, 45, 45, 270, false, false)

	// iceball
	registerResourceFromFileEx("assets/iceball.png", RESOURCE_TYPE_ICEBALL_0, 30, 30, 0, false, false)
	registerResourceFromFileEx("assets/iceball.png", RESOURCE_TYPE_ICEBALL_1, 30, 30, 90, false, false)
	registerResourceFromFileEx("assets/iceball.png", RESOURCE_TYPE_ICEBALL_2, 30, 30, 180, false, false)
	registerResourceFromFileEx("assets/iceball.png", RESOURCE_TYPE_ICEBALL_3, 30, 30, 270, false, false)

	// apple
	registerResourceFromFileEx("assets/apple.png", RESOURCE_TYPE_APPLE_0, 30, 30, 0, false, false)
	registerResourceFromFileEx("assets/apple.png", RESOURCE_TYPE_APPLE_1, 30, 30, 90, false, false)
	registerResourceFromFileEx("assets/apple.png", RESOURCE_TYPE_APPLE_2, 30, 30, 180, false, false)
	registerResourceFromFileEx("assets/apple.png", RESOURCE_TYPE_APPLE_3, 30, 30, 270, false, false)

	// cherry
	registerResourceFromFileEx("assets/cherry.png", RESOURCE_TYPE_CHERRY_0, 30, 30, 0, false, false)
	registerResourceFromFileEx("assets/cherry.png", RESOURCE_TYPE_CHERRY_1, 30, 30, 90, false, false)
	registerResourceFromFileEx("assets/cherry.png", RESOURCE_TYPE_CHERRY_2, 30, 30, 180, false, false)
	registerResourceFromFileEx("assets/cherry.png", RESOURCE_TYPE_CHERRY_3, 30, 30, 270, false, false)

	// moon
	registerResourceFromFileEx("assets/moon.png", RESOURCE_TYPE_MOON_0, 45, 45, 0, false, false)
	registerResourceFromFileEx("assets/moon.png", RESOURCE_TYPE_MOON_1, 45, 45, 90, false, false)
	registerResourceFromFileEx("assets/moon.png", RESOURCE_TYPE_MOON_2, 45, 45, 180, false, false)
	registerResourceFromFileEx("assets/moon.png", RESOURCE_TYPE_MOON_3, 45, 45, 270, false, false)

	// axe
	registerResourceFromFileEx("assets/axe.png", RESOURCE_TYPE_AXE_0, 30, 30, 0, false, false)
	registerResourceFromFileEx("assets/axe.png", RESOURCE_TYPE_AXE_1, 30, 30, 90, false, false)
	registerResourceFromFileEx("assets/axe.png", RESOURCE_TYPE_AXE_2, 30, 30, 180, false, false)
	registerResourceFromFileEx("assets/axe.png", RESOURCE_TYPE_AXE_3, 30, 30, 270, false, false)

	// skull
	registerResourceFromFileEx("assets/skull.png", RESOURCE_TYPE_SKULL_0, 45, 45, 0, false, false)
	registerResourceFromFileEx("assets/skull.png", RESOURCE_TYPE_SKULL_1, 45, 45, 90, false, false)
	registerResourceFromFileEx("assets/skull.png", RESOURCE_TYPE_SKULL_2, 45, 45, 180, false, false)
	registerResourceFromFileEx("assets/skull.png", RESOURCE_TYPE_SKULL_3, 45, 45, 270, false, false)

	// register resource packs
	registerAllTortoiseResPack()
	registerAllBossBResPack()
}

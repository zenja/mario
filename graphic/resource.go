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
	SetResourceAlpha(alpha uint8)
}

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

	RESOURCE_TYPE_PAYPAL_IS_NEW_MONEY

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

	RESOURCE_TYPE_TORTOISE_RED_LEFT_0
	RESOURCE_TYPE_TORTOISE_RED_LEFT_1
	RESOURCE_TYPE_TORTOISE_RED_RIGHT_0
	RESOURCE_TYPE_TORTOISE_RED_RIGHT_1
	RESOURCE_TYPE_TORTOISE_RED_INSIDE
	RESOURCE_TYPE_TORTOISE_RED_SEMI_INSIDE

	RESOURCE_TYPE_BOSS_A_LEFT_0
	RESOURCE_TYPE_BOSS_A_LEFT_1
	RESOURCE_TYPE_BOSS_A_RIGHT_0
	RESOURCE_TYPE_BOSS_A_RIGHT_1

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

func (tr *TileResource) SetResourceAlpha(alpha uint8) {
	tr.GetTexture().SetAlphaMod(alpha)
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

func (ntr *NonTileResource) SetResourceAlpha(alpha uint8) {
	ntr.GetTexture().SetAlphaMod(alpha)
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
	//fmt.Printf("Camera: %d, %d\n", xCamStart, yCamStart)
	//fmt.Printf("Object rect: %d, %d, %d, %d\n", rect.X, rect.Y, rect.W, rect.H)
	//fmt.Printf("Rect in level: %d, %d, %d, %d\n", xStartInLevel, yStartInLevel, xEndInLevel-xStartInLevel, yEndInLevel-yStartInLevel)
	//fmt.Printf("Rect in tile: %d, %d, %d, %d\n", rectInTile.X, rectInTile.Y, rectInTile.W, rectInTile.H)
	//fmt.Printf("Rect in Camera: %d, %d, %d, %d\n", rectInCamera.X, rectInCamera.Y, rectInCamera.W, rectInCamera.H)
	//fmt.Println()
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
	renderer.SetDrawColor(r, g, b, a)
	_, rectInCam := VisibleRectInCamera(rect, camPos.X, camPos.Y)
	renderer.DrawRect(rectInCam)
	renderer.SetDrawColor(originR, originG, originB, originA)
}

func FillRect(rect sdl.Rect, camPos vector.Pos, r, g, b, a uint8) {
	originR, originG, originB, originA, err := renderer.GetDrawColor()
	if err != nil {
		log.Fatalf("failed to get draw color: %s", err)
	}
	renderer.SetDrawColor(r, g, b, a)
	_, rectInCam := VisibleRectInCamera(rect, camPos.X, camPos.Y)
	renderer.FillRect(rectInCam)
	renderer.SetDrawColor(originR, originG, originB, originA)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Graphic functions relative to resource
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

// registerTileResource loads a sprite into a TileResource from a file
func registerTileResource(filename string, id ResourceID) {
	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	registerTileFromSurface(surface, id)
}

// registerNonTileResource loads a sprite into a NonTileResource from a file
func registerNonTileResource(filename string, id ResourceID) {
	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	registerNonTileFromSurface(surface, id)
}

func registerScaledNonTileResource(filename string, id ResourceID, dstWidth int32, dstHeight int32) {
	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	registerScaledNonTileFromSurface(surface, id, dstWidth, dstHeight)
}

func registerFlippedNonTileResource(filename string, id ResourceID, flipHorizontal bool) {
	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	registerFlippedNonTileFromSurface(surface, id, flipHorizontal)
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

	registerScaledNonTileFromSurface(surface, id, dstWidth, dstHeight)
}

func registerResourceEx(
	filename string,
	id ResourceID,
	width,
	height int32,
	isTile bool,
	flipHorizontal bool,
	flipVertical bool) {

	surface, err := img.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Free()

	registerResourceFromSurfaceEx(surface, id, width, height, isTile, flipHorizontal, flipVertical)
}

// registerTileFromSurface loads a sprite from a surface into a TileResource object
// User need to free the surface himself
func registerTileFromSurface(surface *sdl.Surface, id ResourceID) {
	registerResourceFromSurface(surface, id, TILE_SIZE, TILE_SIZE, true)
}

// registerNonTileFromSurface loads a sprite from a surface into a NonTileResource object
// User need to free the surface himself
func registerNonTileFromSurface(surface *sdl.Surface, id ResourceID) {
	registerResourceFromSurface(surface, id, surface.W, surface.H, false)
}

func registerScaledNonTileFromSurface(surface *sdl.Surface, id ResourceID, dstWidth, dstHeight int32) {
	registerResourceFromSurface(surface, id, dstWidth, dstHeight, false)
}

func registerFlippedNonTileFromSurface(surface *sdl.Surface, id ResourceID, flipHorizontal bool) {
	registerFlippedResourceFromSurface(surface, id, surface.W, surface.H, false, flipHorizontal)
}

// registerResourceFromSurface loads a sprite from a surface into a Resource object
// User need to free the surface himself
func registerResourceFromSurface(surface *sdl.Surface, id ResourceID, width, height int32, isTile bool) {
	registerResourceFromSurfaceEx(surface, id, width, height, isTile, false, false)
}

func registerFlippedResourceFromSurface(
	surface *sdl.Surface,
	id ResourceID,
	width,
	height int32,
	isTile bool,
	flipHorizontal bool) {

	registerResourceFromSurfaceEx(surface, id, width, height, isTile, flipHorizontal, !flipHorizontal)
}

func registerResourceFromSurfaceEx(
	surface *sdl.Surface,
	id ResourceID,
	width,
	height int32,
	isTile bool,
	flipHorizontal bool,
	flipVertical bool) {

	if isTile && (width != TILE_SIZE || height != TILE_SIZE) {
		log.Fatalf("declared to be tile but has wrong width (%d) or height (%d)", width, height)
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatal(err)
	}

	// make sure the tile is in good shape
	if surface.W != width || surface.H != height {
		oldTexture := texture
		texture, err = clipTexture(oldTexture, &sdl.Rect{0, 0, width, height})
		if err != nil {
			log.Fatal(err)
		}
		// release original texture
		oldTexture.Destroy()
	}

	// flip texture if needed
	if flipHorizontal {
		oldTexture := texture
		texture, err = flipTexture(texture, width, height, true)
		if err != nil {
			log.Fatal(err)
		}
		oldTexture.Destroy()
	}
	if flipVertical {
		oldTexture := texture
		texture, err = flipTexture(texture, width, height, false)
		if err != nil {
			log.Fatal(err)
		}
		oldTexture.Destroy()
	}

	if isTile {
		resourceRegistry[id] = &TileResource{texture: texture}
	} else {
		resourceRegistry[id] = &NonTileResource{texture: texture, w: width, h: height}
	}
}

// RenderResource renders a tile (or a part of tile specified by srcRect) to a given position in screen
func RenderResource(resource Resource, srcRect *sdl.Rect, dstRect *sdl.Rect) {
	renderer.Copy(resource.GetTexture(), srcRect, dstRect)
}

func loadAllResources() {
	// -------------------------------
	// load tile resources
	// -------------------------------

	// brick
	registerTileResource("assets/brick-red.png", RESOURCE_TYPE_BRICK_RED)
	registerTileResource("assets/brick-yellow.png", RESOURCE_TYPE_BRICK_YELLOW)

	// grass
	registerTileResource("assets/grass-ground-left.png", RESOURCE_TYPE_GRASS_GROUD_LEFT)
	registerTileResource("assets/grass-ground-mid.png", RESOURCE_TYPE_GRASS_GROUD_MID)
	registerTileResource("assets/grass-ground-right.png", RESOURCE_TYPE_GRASS_GROUD_RIGHT)

	// ground
	registerTileResource("assets/ground-left.png", RESOURCE_TYPE_GROUD_LEFT)
	registerTileResource("assets/ground-mid.png", RESOURCE_TYPE_GROUD_MID)
	registerTileResource("assets/ground-right.png", RESOURCE_TYPE_GROUD_RIGHT)

	// myth box
	registerTileResource("assets/myth-box-normal.png", RESOURCE_TYPE_MYTH_BOX_NORMAL)
	registerTileResource("assets/myth-box-normal-light.png", RESOURCE_TYPE_MYTH_BOX_NORMAL_LIGHT)
	registerTileResource("assets/myth-box-empty.png", RESOURCE_TYPE_MYTH_BOX_EMPTY)

	// pipe
	registerTileResource("assets/pipe-left-top.png", RESOURCE_TYPE_PIPE_LEFT_TOP)
	registerTileResource("assets/pipe-right-top.png", RESOURCE_TYPE_PIPE_RIGHT_TOP)
	registerTileResource("assets/pipe-left-mid.png", RESOURCE_TYPE_PIPE_LEFT_MID)
	registerTileResource("assets/pipe-right-mid.png", RESOURCE_TYPE_PIPE_RIGHT_MID)
	registerTileResource("assets/pipe-left-bottom.png", RESOURCE_TYPE_PIPE_LEFT_BOTTOM)
	registerTileResource("assets/pipe-right-bottom.png", RESOURCE_TYPE_PIPE_RIGHT_BOTTOM)

	// coin
	registerTileResource("assets/coin-0.png", RESOURCE_TYPE_COIN_0)
	registerTileResource("assets/coin-1.png", RESOURCE_TYPE_COIN_1)
	registerTileResource("assets/coin-2.png", RESOURCE_TYPE_COIN_2)
	registerTileResource("assets/coin-3.png", RESOURCE_TYPE_COIN_3)

	// good mushroom
	registerTileResource("assets/mushroom.png", RESOURCE_TYPE_GOOD_MUSHROOM)

	// upgrade flower
	registerTileResource("assets/upgrade-flower.png", RESOURCE_TYPE_UPGRADE_FLOWER)

	// water
	registerTileResource("assets/water-0.png", RESOURCE_TYPE_WATER_0)
	registerTileResource("assets/water-1.png", RESOURCE_TYPE_WATER_1)
	registerTileResource("assets/water-2.png", RESOURCE_TYPE_WATER_2)
	registerTileResource("assets/water-3.png", RESOURCE_TYPE_WATER_3)
	registerTileResource("assets/water-4.png", RESOURCE_TYPE_WATER_4)
	registerTileResource("assets/water-5.png", RESOURCE_TYPE_WATER_5)
	registerTileResource("assets/water-6.png", RESOURCE_TYPE_WATER_6)
	registerTileResource("assets/water-pixel.png", RESOURCE_TYPE_WATER_FULL)

	// -------------------------------
	// Load non-tile resources
	// -------------------------------

	// hero 0
	registerFacedResource("assets/hero-0-stand.png", "xwang16", RESOURCE_TYPE_HERO_0_STAND_RIGHT,
		hero_0_width, hero_0_height, 40, 50, 5, 0, 5, false, false)
	registerFacedResource("assets/hero-0-walking.png", "xwang16", RESOURCE_TYPE_HERO_0_WALKING_RIGHT,
		hero_0_width, hero_0_height, 40, 50, 5, 0, -5, false, false)
	registerFacedResource("assets/hero-0-jump.png", "xwang16", RESOURCE_TYPE_HERO_0_JUMP_RIGHT,
		hero_0_width, hero_0_height, 40, 50, 5, 0, 0, false, false)
	registerFacedResource("assets/hero-0-stand.png", "xwang16", RESOURCE_TYPE_HERO_0_STAND_LEFT,
		hero_0_width, hero_0_height, 40, 50, 5, 0, 5, true, false)
	registerFacedResource("assets/hero-0-walking.png", "xwang16", RESOURCE_TYPE_HERO_0_WALKING_LEFT,
		hero_0_width, hero_0_height, 40, 50, 5, 0, -5, true, false)
	registerFacedResource("assets/hero-0-jump.png", "xwang16", RESOURCE_TYPE_HERO_0_JUMP_LEFT,
		hero_0_width, hero_0_height, 40, 50, 5, 0, 0, true, false)

	// hero 1
	registerFacedResource("assets/hero-1-stand.png", "xwang16", RESOURCE_TYPE_HERO_1_STAND_RIGHT,
		hero_1_width, hero_1_height, 45, 55, 5, 0, 5, false, false)
	registerFacedResource("assets/hero-1-walking.png", "xwang16", RESOURCE_TYPE_HERO_1_WALKING_RIGHT,
		hero_1_width, hero_1_height, 45, 55, 5, 0, -5, false, false)
	registerFacedResource("assets/hero-1-jump.png", "xwang16", RESOURCE_TYPE_HERO_1_JUMP_RIGHT,
		hero_1_width, hero_1_height, 45, 55, 5, 0, 0, false, false)
	registerFacedResource("assets/hero-1-stand.png", "xwang16", RESOURCE_TYPE_HERO_1_STAND_LEFT,
		hero_1_width, hero_1_height, 45, 55, 5, 0, 5, true, false)
	registerFacedResource("assets/hero-1-walking.png", "xwang16", RESOURCE_TYPE_HERO_1_WALKING_LEFT,
		hero_1_width, hero_1_height, 45, 55, 5, 0, -5, true, false)
	registerFacedResource("assets/hero-1-jump.png", "xwang16", RESOURCE_TYPE_HERO_1_JUMP_LEFT,
		hero_1_width, hero_1_height, 45, 55, 5, 0, 0, true, false)

	// hero 2
	registerFacedResource("assets/hero-2-stand.png", "xwang16", RESOURCE_TYPE_HERO_2_STAND_RIGHT,
		hero_2_width, hero_2_height, 45, 55, 5, 0, 5, false, false)
	registerFacedResource("assets/hero-2-walking.png", "xwang16", RESOURCE_TYPE_HERO_2_WALKING_RIGHT,
		hero_2_width, hero_2_height, 45, 55, 5, 0, -5, false, false)
	registerFacedResource("assets/hero-2-jump.png", "xwang16", RESOURCE_TYPE_HERO_2_JUMP_RIGHT,
		hero_2_width, hero_2_height, 45, 55, 5, 0, 0, false, false)
	registerFacedResource("assets/hero-2-stand.png", "xwang16", RESOURCE_TYPE_HERO_2_STAND_LEFT,
		hero_2_width, hero_2_height, 45, 55, 5, 0, 5, true, false)
	registerFacedResource("assets/hero-2-walking.png", "xwang16", RESOURCE_TYPE_HERO_2_WALKING_LEFT,
		hero_2_width, hero_2_height, 45, 55, 5, 0, -5, true, false)
	registerFacedResource("assets/hero-2-jump.png", "xwang16", RESOURCE_TYPE_HERO_2_JUMP_LEFT,
		hero_2_width, hero_2_height, 45, 55, 5, 0, 0, true, false)

	// decoration: grass
	registerNonTileResource("assets/dec-grass-0.png", RESOURCE_TYPE_DEC_GRASS_0)
	registerNonTileResource("assets/dec-grass-1.png", RESOURCE_TYPE_DEC_GRASS_1)

	// decoration: trees
	registerNonTileResource("assets/dec-tree-0.png", RESOURCE_TYPE_DEC_TREE_0)

	// decoration: paypal is new money
	registerNonTileResource("assets/paypal-is-new-money.png", RESOURCE_TYPE_PAYPAL_IS_NEW_MONEY)

	// broken pieces
	registerScaledNonTileResource("assets/brick-piece-red.png", RESOURCE_TYPE_BRICK_PIECE_RED, TILE_SIZE/2, TILE_SIZE/2)
	registerScaledNonTileResource("assets/brick-piece-yellow.png", RESOURCE_TYPE_BRICK_PIECE_YELLOW, TILE_SIZE/2, TILE_SIZE/2)

	// mushroom enemy
	registerScaledNonTileResource("assets/mushroom-enemy-0.png", RESOURCE_TYPE_MUSHROOM_ENEMY_0, TILE_SIZE, TILE_SIZE)
	registerScaledNonTileResource("assets/mushroom-enemy-1.png", RESOURCE_TYPE_MUSHROOM_ENEMY_1, TILE_SIZE, TILE_SIZE)
	registerScaledNonTileResource("assets/mushroom-enemy-hit.png", RESOURCE_TYPE_MUSHROOM_ENEMY_HIT, TILE_SIZE, TILE_SIZE)
	registerResourceEx("assets/mushroom-enemy-0.png", RESOURCE_TYPE_MUSHROOM_ENEMY_DOWN, TILE_SIZE, TILE_SIZE, false, false, true)

	// tortoise enemy
	registerFacedResource("assets/tortoise-red-right-0.png", "chran", RESOURCE_TYPE_TORTOISE_RED_RIGHT_0,
		tortoise_walking_width, tortoise_walking_height, 30, 45, 15, 0, 5, false, false)
	registerFacedResource("assets/tortoise-red-right-1.png", "chran", RESOURCE_TYPE_TORTOISE_RED_RIGHT_1,
		tortoise_walking_width, tortoise_walking_height, 30, 45, 15, 0, -5, false, false)
	registerFacedResource("assets/tortoise-red-right-0.png", "chran", RESOURCE_TYPE_TORTOISE_RED_LEFT_0,
		tortoise_walking_width, tortoise_walking_height, 30, 45, 0, 0, 5, true, false)
	registerFacedResource("assets/tortoise-red-right-1.png", "chran", RESOURCE_TYPE_TORTOISE_RED_LEFT_1,
		tortoise_walking_width, tortoise_walking_height, 30, 45, 0, 0, -5, true, false)
	registerScaledNonTileResource("assets/tortoise-red-inside.png", RESOURCE_TYPE_TORTOISE_RED_INSIDE, tortoise_inside_width, tortoise_inside_height)
	registerScaledNonTileResource("assets/tortoise-red-semi-inside.png", RESOURCE_TYPE_TORTOISE_RED_SEMI_INSIDE, tortoise_inside_width, tortoise_inside_height)

	// boss A
	registerFacedResource("assets/boss-a-right-0.png", "yunfeli", RESOURCE_TYPE_BOSS_A_RIGHT_0,
		149, 140, 65, 90, 80, 0, 5, false, false)
	registerFacedResource("assets/boss-a-right-1.png", "yunfeli", RESOURCE_TYPE_BOSS_A_RIGHT_1,
		149, 140, 65, 90, 80, 0, -5, false, false)
	registerFacedResource("assets/boss-a-right-0.png", "yunfeli", RESOURCE_TYPE_BOSS_A_LEFT_0,
		149, 140, 65, 90, 0, 0, 5, true, false)
	registerFacedResource("assets/boss-a-right-1.png", "yunfeli", RESOURCE_TYPE_BOSS_A_LEFT_1,
		149, 140, 65, 90, 0, 0, -5, true, false)

	// boss boom
	registerScaledNonTileResource("assets/boss-boom.png", RESOURCE_TYPE_BOSS_BOOM, 300, 200)

	// fireball
	registerScaledNonTileResource("assets/fireball-0.png", RESOURCE_TYPE_FIREBALL_0, 30, 30)
	registerScaledNonTileResource("assets/fireball-1.png", RESOURCE_TYPE_FIREBALL_1, 30, 30)
	registerScaledNonTileResource("assets/fireball-2.png", RESOURCE_TYPE_FIREBALL_2, 30, 30)
	registerScaledNonTileResource("assets/fireball-3.png", RESOURCE_TYPE_FIREBALL_3, 30, 30)
	registerScaledNonTileResource("assets/fireball-boom.png", RESOURCE_TYPE_FIREBALL_BOOM, 40, 40)

	// shit
	registerScaledNonTileResource("assets/shit-0.png", RESOURCE_TYPE_SHIT_0, 30, 30)
	registerScaledNonTileResource("assets/shit-1.png", RESOURCE_TYPE_SHIT_1, 30, 30)
	registerScaledNonTileResource("assets/shit-2.png", RESOURCE_TYPE_SHIT_2, 30, 30)
	registerScaledNonTileResource("assets/shit-3.png", RESOURCE_TYPE_SHIT_3, 30, 30)
	registerScaledNonTileResource("assets/shit-boom.png", RESOURCE_TYPE_SHIT_BOOM, 50, 50)

	// bug
	registerScaledNonTileResource("assets/bug-0.png", RESOURCE_TYPE_BUG_0, 30, 30)
	registerScaledNonTileResource("assets/bug-1.png", RESOURCE_TYPE_BUG_1, 30, 30)
	registerScaledNonTileResource("assets/bug-2.png", RESOURCE_TYPE_BUG_2, 30, 30)
	registerScaledNonTileResource("assets/bug-3.png", RESOURCE_TYPE_BUG_3, 30, 30)
	registerScaledNonTileResource("assets/bug-boom.png", RESOURCE_TYPE_BUG_BOOM, 50, 50)

	// eater flower
	registerScaledNonTileResource("assets/eater-flower-0.png", RESOURCE_TYPE_EATER_FLOWER_0, 52, 75)
	registerScaledNonTileResource("assets/eater-flower-1.png", RESOURCE_TYPE_EATER_FLOWER_1, 52, 75)

	// shine effect
	registerNonTileResource("assets/shine-0.png", RESOURCE_TYPE_SHINE_0)
	registerNonTileResource("assets/shine-1.png", RESOURCE_TYPE_SHINE_1)
	registerNonTileResource("assets/shine-2.png", RESOURCE_TYPE_SHINE_2)

	// bang
	registerScaledNonTileResource("assets/bang.png", RESOURCE_TYPE_BANG, 50, 50)

	// black screen
	registerScaledNonTileResource("assets/black-pixel.png", RESOURCE_TYPE_BLACK_SCREEN, SCREEN_WIDTH, SCREEN_HEIGHT)
}

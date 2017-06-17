package level

import (
	"container/list"
	"log"

	"strings"

	"strconv"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

const (
	level_dir = "assets/levels"
)

type levelSpec struct {
	name          string
	nextLevelName string
	bgFilename    string // file name of background file
	bgColor       sdl.Color
	levelArr      [][]byte
	decArr        [][]byte // decoration array
}

func BuildLevel(spec *levelSpec) *Level {
	graphic.RegisterBackgroundResource(spec.bgFilename, graphic.RESOURCE_TYPE_BG_0, len(spec.levelArr))
	bgRes := graphic.Res(graphic.RESOURCE_TYPE_BG_0)

	// NOTE: index is tid.X, tid.Y
	var tileObjs [][]Object

	var enemies []Enemy

	numTiles := vector.Vec2D{int32(len(spec.levelArr[0])), int32(len(spec.levelArr))}
	obstMngr := NewObstacleManager(len(spec.levelArr[0]), len(spec.levelArr))
	enemyObstMngr := NewObstacleManager(len(spec.levelArr[0]), len(spec.levelArr))
	var hero *Hero

	// init tileObjs array
	for i := 0; i < int(numTiles.X); i++ {
		tileObjs = append(tileObjs, make([]Object, numTiles.Y))
	}

	addAsFullObstTile := func(tid vector.TileID, o Object) {
		tileObjs[tid.X][tid.Y] = o
		obstMngr.AddTileObst(tid)
		enemyObstMngr.AddTileObst(tid)
	}

	addAsEnemyOnlyObstTile := func(tid vector.TileID, o Object) {
		tileObjs[tid.X][tid.Y] = o
		enemyObstMngr.AddTileObst(tid)
	}

	addAsNoObstTile := func(tid vector.TileID, o Object) {
		tileObjs[tid.X][tid.Y] = o
	}

	var decorations []Object
	addDecoration := func(d *decoration) {
		decorations = append(decorations, d)
	}

	// parse level
	var currentPos vector.Pos
	for tidY := 0; tidY < int(numTiles.Y); tidY++ {
		currentPos.X = 0
		for tidX := 0; tidX < int(numTiles.X); tidX++ {
			tid := vector.TileID{int32(tidX), int32(tidY)}
			// note that levelArr's index is not TID, need reverse
			switch spec.levelArr[tidY][tidX] {
			// Invisible block
			case '#':
				addAsFullObstTile(tid, NewInvisibleTileObject(tid))

				// Invisible block only to enemies
			case '"':
				addAsEnemyOnlyObstTile(tid, NewInvisibleTileObject(tid))

				// Brick
			case 'B':
				mainRes := graphic.Res(graphic.RESOURCE_TYPE_BRICK)
				pieceRes := graphic.Res(graphic.RESOURCE_TYPE_BRICK_PIECE)
				o := NewBreakableTileObject(mainRes, pieceRes, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

				// Ground with left grass
			case 'L':
				res := graphic.Res(graphic.RESOURCE_TYPE_GROUD_GRASS_LEFT)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

				// Ground with mid grass
			case 'G':
				res := graphic.Res(graphic.RESOURCE_TYPE_GROUD_GRASS_MID)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

				// Ground with right grass
			case 'R':
				res := graphic.Res(graphic.RESOURCE_TYPE_GROUD_GRASS_RIGHT)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

				// Inner ground in middle
			case 'I':
				res := graphic.Res(graphic.RESOURCE_TYPE_GROUD_INNER_MID)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

				// Myth box for coins
			case 'C':
				addAsFullObstTile(tid, NewCoinMythBox(currentPos, 3))

				// Myth box for mushrooms
			case 'M':
				addAsFullObstTile(tid, NewMushroomMythBox(currentPos))

				// left middle of pipe
			case '[':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_MID)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

				// right middle of pipe
			case ']':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_MID)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

				// left top of pipe
			case '{':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_TOP)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

				// right top of pipe
			case '}':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_TOP)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

				// left bottom of pipe
			case '<':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_BOTTOM)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

				// right bottom of pipe
			case '>':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_BOTTOM)
				o := NewSingleTileObject(res, currentPos, ZINDEX_0)
				addAsFullObstTile(tid, o)

				// water surface
			case 'W':
				o := NewWaterSurfaceAnimationObject(tid)
				addAsNoObstTile(tid, o)

				// water inside
			case 'w':
				res := graphic.Res(graphic.RESOURCE_TYPE_WATER_FULL)
				o := NewSingleTileObject(res, currentPos, ZINDEX_1)
				addAsNoObstTile(tid, o)

				// Enemy 1: mushroom enemy
			case '1':
				enemies = append(enemies, NewMushroomEnemy(currentPos))

				// Enemy 2: tortoise enemy
			case '2':
				enemies = append(enemies, NewTortoiseEnemy(currentPos))

				// Hero
			case 'H':
				if hero != nil {
					log.Fatal("more than one hero found")
				}
				hero = NewHero(currentPos, 0.2, 0.1)
			}
			currentPos.X += graphic.TILE_SIZE
		}
		currentPos.Y += graphic.TILE_SIZE
	}

	// parse decorations
	currentPos = vector.Pos{}
	for tidY := 0; tidY < int(numTiles.Y); tidY++ {
		currentPos.X = 0
		for tidX := 0; tidX < int(numTiles.X); tidX++ {
			tid := vector.TileID{int32(tidX), int32(tidY)}
			// note that decArr's index is not TID, need reverse
			switch spec.decArr[tidY][tidX] {
			case '1':
				resIds := []graphic.ResourceID{
					graphic.RESOURCE_TYPE_DEC_GRASS_0,
					graphic.RESOURCE_TYPE_DEC_GRASS_1,
				}
				addDecoration(NewDecoration(tid, resIds, 800))
			}
		}
	}

	if hero == nil {
		log.Fatal("no hero found when parsing level")
	}

	return &Level{
		BGRes:         bgRes,
		Decorations:   decorations,
		TileObjects:   tileObjs,
		Enemies:       enemies,
		VolatileObjs:  list.New(),
		ObstMngr:      obstMngr,
		EnemyObstMngr: enemyObstMngr,
		TheHero:       hero,
		InitHeroRect:  hero.levelRect,
		BGColor:       spec.bgColor,
		NumTiles:      numTiles,
		effects:       list.New(),
	}
}

func ParseLevelSpec(levelFile string) *levelSpec {
	conf, err := toml.LoadFile(levelFile)
	if err != nil {
		log.Fatal(err)
	}

	name := conf.Get("basic.name").(string)
	nextLevelName := conf.Get("basic.next-level-name").(string)

	bgFilename := conf.Get("graphic.bg-file").(string)

	bgColor, err := parseRGB(conf.Get("graphic.bg-color-rgb").(string))
	if err != nil {
		log.Fatal(err)
	}

	levelDef, err := parseLevelArr(conf.Get("level.def").(string))
	if err != nil {
		log.Fatal(err)
	}

	levelDecDef, err := parseLevelArr(conf.Get("level.dec-def").(string))
	if err != nil {
		log.Fatal(err)
	}

	return &levelSpec{
		name:          name,
		nextLevelName: nextLevelName,
		bgFilename:    bgFilename,
		bgColor:       bgColor,
		levelArr:      levelDef,
		decArr:        levelDecDef,
	}
}

func parseRGB(str string) (sdl.Color, error) {
	splits := strings.Split(strings.Replace(str, " ", "", -1), ",")
	if len(splits) != 3 {
		return sdl.Color{}, errors.Errorf("rgb should have three parts: %s", str)
	}

	r, err := strconv.Atoi(splits[0])
	if err != nil {
		return sdl.Color{}, errors.Errorf("failed to parse %s as int in rgb str %s", splits[0], str)
	}

	g, err := strconv.Atoi(splits[1])
	if err != nil {
		return sdl.Color{}, errors.Errorf("failed to parse %s as int in rgb str %s", splits[1], str)
	}

	b, err := strconv.Atoi(splits[2])
	if err != nil {
		return sdl.Color{}, errors.Errorf("failed to parse %s as int in rgb str %s", splits[2], str)
	}

	return sdl.Color{uint8(r), uint8(g), uint8(b), 255}, nil
}

func parseLevelArr(str string) ([][]byte, error) {
	trimmed := strings.Trim(str, "\n")
	var result [][]byte
	lines := strings.Split(trimmed, "\n")

	if len(lines) == 0 {
		return nil, errors.New("failed to parse level arr: input is empty")
	}

	width := len(lines[0])

	for _, l := range lines {
		if len(l) != width {
			return nil, errors.Errorf("failed to parse level arr: first line length is %d, but %d found", width, len(l))
		}
		result = append(result, []byte(l))
	}
	return result, nil
}
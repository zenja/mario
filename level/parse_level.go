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

type LevelSpec struct {
	Name           string
	NextLevelNames []string
	BgFilename     string // file name of background file
	BgColor        sdl.Color
	LevelArr       [][]byte
	DecArr         [][]byte // decoration array
}

func BuildLevel(spec *LevelSpec) *Level {
	graphic.RegisterBackgroundResource(spec.BgFilename, graphic.RESOURCE_TYPE_BG_0, len(spec.LevelArr))
	bgRes := graphic.Res(graphic.RESOURCE_TYPE_BG_0)

	// NOTE: index is tid.X, tid.Y
	var tileObjs [][]Object

	var enemies []Enemy

	numTiles := vector.Vec2D{int32(len(spec.LevelArr[0])), int32(len(spec.LevelArr))}
	obstMngr := NewObstacleManager(len(spec.LevelArr[0]), len(spec.LevelArr))
	enemyObstMngr := NewObstacleManager(len(spec.LevelArr[0]), len(spec.LevelArr))
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
	var nextLevelJumperIdx int = 0
	for tidY := 0; tidY < int(numTiles.Y); tidY++ {
		currentPos.X = 0
		for tidX := 0; tidX < int(numTiles.X); tidX++ {
			tid := vector.TileID{int32(tidX), int32(tidY)}
			// note that levelArr's index is not TID, need reverse
			switch spec.LevelArr[tidY][tidX] {
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
				o := NewSingleTileObject(res, currentPos, ZINDEX_4)
				addAsFullObstTile(tid, o)

			// right middle of pipe
			case ']':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_MID)
				o := NewSingleTileObject(res, currentPos, ZINDEX_4)
				addAsFullObstTile(tid, o)

			// left top of pipe that will jump level
			case '{':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_TOP)
				o := NewSingleTileObject(res, currentPos, ZINDEX_4)
				addAsFullObstTile(tid, o)

				// also needs to add level jumper
				nextLevelName := spec.NextLevelNames[nextLevelJumperIdx]
				enemies = append(enemies, NewLevelJumper(tid, nextLevelName))
				nextLevelJumperIdx++

			// right top of pipe that will jump level
			case '}':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_TOP)
				o := NewSingleTileObject(res, currentPos, ZINDEX_4)
				addAsFullObstTile(tid, o)

			// left bottom of pipe
			case '<':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_BOTTOM)
				o := NewSingleTileObject(res, currentPos, ZINDEX_4)
				addAsFullObstTile(tid, o)

			// right bottom of pipe
			case '>':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_BOTTOM)
				o := NewSingleTileObject(res, currentPos, ZINDEX_4)
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
				hero = NewHero(currentPos, 0.2, 0.2)
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
			switch spec.DecArr[tidY][tidX] {
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
		Spec:          spec,
		BGRes:         bgRes,
		Decorations:   decorations,
		TileObjects:   tileObjs,
		Enemies:       enemies,
		VolatileObjs:  list.New(),
		ObstMngr:      obstMngr,
		EnemyObstMngr: enemyObstMngr,
		TheHero:       hero,
		InitHeroRect:  hero.levelRect,
		BGColor:       spec.BgColor,
		NumTiles:      numTiles,
		effects:       list.New(),
	}
}

func ParseLevelSpec(levelFile string) *LevelSpec {
	conf, err := toml.LoadFile(levelFile)
	if err != nil {
		log.Fatal(err)
	}

	name := conf.Get("basic.name").(string)
	var nextLevelNames []string
	for _, name := range conf.Get("transfer.next-levels").([]interface{}) {
		nextLevelNames = append(nextLevelNames, name.(string))
	}

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

	// check if num of next levels equals to '{' in level definition
	leftBracketCnt := strings.Count(conf.Get("level.def").(string), "{")
	if leftBracketCnt != len(nextLevelNames) {
		log.Fatalf("failed to parse level %s: there are %d next levels but %d '{'",
			name, len(nextLevelNames), leftBracketCnt)
	}

	return &LevelSpec{
		Name:           name,
		NextLevelNames: nextLevelNames,
		BgFilename:     bgFilename,
		BgColor:        bgColor,
		LevelArr:       levelDef,
		DecArr:         levelDecDef,
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

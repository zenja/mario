package level

import (
	"container/list"
	"log"

	"strings"

	"strconv"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/audio"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type LevelSpec struct {
	Name           string
	NextLevelNames []string
	BgFilename     string // file name of background file
	BgColor        sdl.Color
	BgMusicID      audio.MusicID
	LevelArr       [][]byte
	DecArr         [][]byte // decoration array
}

func BuildLevel(spec *LevelSpec) *Level {
	graphic.RegisterBackgroundResource(spec.BgFilename, graphic.RESOURCE_TYPE_CURR_BG, len(spec.LevelArr))
	bgRes := graphic.Res(graphic.RESOURCE_TYPE_CURR_BG)

	// NOTE: index is tid.X, tid.Y
	var tileObjs [][]Object

	var enemies []Enemy

	numTiles := vector.Vec2D{int32(len(spec.LevelArr[0])), int32(len(spec.LevelArr))}
	obstMngr := NewObstacleManager(len(spec.LevelArr[0]), len(spec.LevelArr))
	var hero *Hero

	// init tileObjs array
	for i := 0; i < int(numTiles.X); i++ {
		tileObjs = append(tileObjs, make([]Object, numTiles.Y))
	}

	addAsNormalObstTile := func(tid vector.TileID, o Object) {
		tileObjs[tid.X][tid.Y] = o
		obstMngr.AddNormalTileObst(tid)
	}

	addAsEnemyOnlyObstTile := func(tid vector.TileID, o Object) {
		tileObjs[tid.X][tid.Y] = o
		obstMngr.AddEnemyOnlyTileObst(tid)
	}

	addAsUpThruObstTile := func(tid vector.TileID, o Object) {
		tileObjs[tid.X][tid.Y] = o
		obstMngr.AddUpThruTileObst(tid)
	}

	addAsNoObstTile := func(tid vector.TileID, o Object) {
		tileObjs[tid.X][tid.Y] = o
	}

	needAddGroundLeft := func(tid vector.TileID) bool {
		leftSpec := spec.LevelArr[tid.Y][tid.X-1]
		if tid.X-1 > 0 && (leftSpec == 'l' || leftSpec == 'L' || leftSpec == 'g') {
			return true
		}
		return false
	}

	needAddGroundRight := func(tid vector.TileID) bool {
		rightSpec := spec.LevelArr[tid.Y][tid.X+1]
		if tid.X+1 < numTiles.X && (rightSpec == 'r' || rightSpec == 'R' || rightSpec == 'g') {
			return true
		}
		return false
	}

	var decorations []Object
	addDecoration := func(d *decoration) {
		decorations = append(decorations, d)
	}
	addTextDecoration := func(d *textDecoration) {
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
				addAsNormalObstTile(tid, NewInvisibleTileObject(tid))

			// Invisible block only to enemies
			case '\'':
				addAsEnemyOnlyObstTile(tid, NewInvisibleTileObject(tid))

			// Brick: yellow
			case 'B':
				mainRes := graphic.Res(graphic.RESOURCE_TYPE_BRICK_YELLOW)
				pieceRes := graphic.Res(graphic.RESOURCE_TYPE_BRICK_PIECE_YELLOW)
				o := NewBreakableTileObject(mainRes, pieceRes, currentPos, ZINDEX_0)
				addAsNormalObstTile(tid, o)

			// Brick: red
			case 'D':
				mainRes := graphic.Res(graphic.RESOURCE_TYPE_BRICK_RED)
				pieceRes := graphic.Res(graphic.RESOURCE_TYPE_BRICK_PIECE_RED)
				o := NewBreakableTileObject(mainRes, pieceRes, currentPos, ZINDEX_0)
				addAsNormalObstTile(tid, o)

			// Ground with left grass
			case 'L':
				resID := graphic.RESOURCE_TYPE_GRASS_GROUD_LEFT
				res := graphic.Res(resID)
				var o Object
				if needAddGroundLeft(tid) {
					o = NewOverlapTilesObject(
						[]graphic.ResourceID{graphic.RESOURCE_TYPE_GROUD_MID, resID}, tid, ZINDEX_1)
				} else {
					o = NewSingleTileObject(res, tid, ZINDEX_1)
				}
				addAsUpThruObstTile(tid, o)

			// Ground with mid grass
			case 'G':
				res := graphic.Res(graphic.RESOURCE_TYPE_GRASS_GROUD_MID)
				o := NewSingleTileObject(res, tid, ZINDEX_0)
				addAsUpThruObstTile(tid, o)

			// Ground with right grass
			case 'R':
				resID := graphic.RESOURCE_TYPE_GRASS_GROUD_RIGHT
				res := graphic.Res(resID)
				var o Object
				if needAddGroundRight(tid) {
					o = NewOverlapTilesObject(
						[]graphic.ResourceID{graphic.RESOURCE_TYPE_GROUD_MID, resID}, tid, ZINDEX_1)
				} else {
					o = NewSingleTileObject(res, tid, ZINDEX_1)
				}
				addAsUpThruObstTile(tid, o)

			// Inner ground in left
			case 'l':
				resID := graphic.RESOURCE_TYPE_GROUD_LEFT
				res := graphic.Res(resID)
				var o Object
				if needAddGroundLeft(tid) {
					o = NewOverlapTilesObject(
						[]graphic.ResourceID{graphic.RESOURCE_TYPE_GROUD_MID, resID}, tid, ZINDEX_1)
				} else {
					o = NewSingleTileObject(res, tid, ZINDEX_1)
				}
				addAsNoObstTile(tid, o)

			// Inner ground in middle
			case 'g':
				res := graphic.Res(graphic.RESOURCE_TYPE_GROUD_MID)
				o := NewSingleTileObject(res, tid, ZINDEX_0)
				addAsNoObstTile(tid, o)

			// Inner ground in right
			case 'r':
				resID := graphic.RESOURCE_TYPE_GROUD_RIGHT
				res := graphic.Res(resID)
				var o Object
				if needAddGroundRight(tid) {
					o = NewOverlapTilesObject(
						[]graphic.ResourceID{graphic.RESOURCE_TYPE_GROUD_MID, resID}, tid, ZINDEX_1)
				} else {
					o = NewSingleTileObject(res, tid, ZINDEX_1)
				}
				addAsNoObstTile(tid, o)

			// Myth box for coins
			case 'C':
				addAsNormalObstTile(tid, NewCoinMythBox(currentPos, 3))

			// Myth box for mushrooms
			case 'M':
				addAsNormalObstTile(tid, NewMushroomMythBox(currentPos))

			// left middle of pipe
			case '[':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_MID)
				o := NewSingleTileObject(res, tid, ZINDEX_4)
				addAsNormalObstTile(tid, o)

			// right middle of pipe
			case ']':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_MID)
				o := NewSingleTileObject(res, tid, ZINDEX_4)
				addAsNormalObstTile(tid, o)

			// right middle of pipe, with eater flower
			case 'E':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_MID)
				o := NewSingleTileObject(res, tid, ZINDEX_4)
				addAsNormalObstTile(tid, o)

				// add eater
				enemies = append(enemies, NewEaterFlower(tid))

			// left top of pipe that will jump level
			case '{':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_TOP)
				o := NewSingleTileObject(res, tid, ZINDEX_4)
				addAsNormalObstTile(tid, o)

				// also needs to add level jumper
				nextLevelName := spec.NextLevelNames[nextLevelJumperIdx]
				enemies = append(enemies, NewLevelJumper(tid, nextLevelName))
				nextLevelJumperIdx++

			// right top of pipe that will jump level
			case '}':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_TOP)
				o := NewSingleTileObject(res, tid, ZINDEX_4)
				addAsNormalObstTile(tid, o)

			// left top of normal pipe
			case '(':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_TOP)
				o := NewSingleTileObject(res, tid, ZINDEX_4)
				addAsNormalObstTile(tid, o)

			// right top of normal pipe
			case ')':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_TOP)
				o := NewSingleTileObject(res, tid, ZINDEX_4)
				addAsNormalObstTile(tid, o)

			// left bottom of pipe
			case '<':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_LEFT_BOTTOM)
				o := NewSingleTileObject(res, tid, ZINDEX_4)
				addAsNormalObstTile(tid, o)

			// right bottom of pipe
			case '>':
				res := graphic.Res(graphic.RESOURCE_TYPE_PIPE_RIGHT_BOTTOM)
				o := NewSingleTileObject(res, tid, ZINDEX_4)
				addAsNormalObstTile(tid, o)

			// water surface
			case 'W':
				o := NewWaterSurfaceAnimationObject(tid)
				addAsNoObstTile(tid, o)

			// water inside
			case 'w':
				res := graphic.Res(graphic.RESOURCE_TYPE_WATER_FULL)
				o := NewSingleTileObject(res, tid, ZINDEX_1)
				addAsNoObstTile(tid, o)

			// coin
			case 'c':
				enemies = append(enemies, NewCoinEnemy(tid))

			// Enemy 1: mushroom enemy
			case '1':
				enemies = append(enemies, NewMushroomEnemy(currentPos))

			// Enemy 2: tortoise enemy
			case '2':
				enemies = append(enemies, NewRandomICTortoiseEnemy(currentPos))

			// Enemy 3: tortoise enemy: richard direct report
			case '3':
				enemies = append(enemies, NewRandomRichardLeadershipTortoiseEnemy(currentPos))

			// Boss A
			case 'X':
				enemies = append(enemies, NewBossA(currentPos))

			// Boss B
			case 'Y':
				enemies = append(enemies, NewRandomBossB(currentPos))

			// Boss C
			case 'Z':
				enemies = append(enemies, NewBossC(currentPos))

			// Boss D
			case 'x':
				enemies = append(enemies, NewBossD(currentPos))

			// Boss E
			case 'y':
				enemies = append(enemies, NewBossE(currentPos))

			// Boss F
			case 'z':
				enemies = append(enemies, NewBossF(currentPos))

			// Boss G
			case 'S':
				enemies = append(enemies, NewBossG(currentPos))

			// Boss H
			case 'T':
				enemies = append(enemies, NewBossH(currentPos))

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
			case '2':
				resIds := []graphic.ResourceID{
					graphic.RESOURCE_TYPE_DEC_TREE_0,
				}
				addDecoration(NewDecoration(tid, resIds, 1000))
			case '3':
				resIds := []graphic.ResourceID{
					graphic.RESOURCE_TYPE_PAYPAL_IS_NEW_MONEY,
				}
				addDecoration(NewDecoration(tid, resIds, 1000))
			case '4':
				resIds := []graphic.ResourceID{
					graphic.RESOURCE_TYPE_DEC_HIGH_ENERGY_AHEAD,
				}
				addDecoration(NewDecoration(tid, resIds, 1000))
			case '5':
				resIds := []graphic.ResourceID{
					graphic.RESOURCE_TYPE_DEC_PRINCESS_0,
					graphic.RESOURCE_TYPE_DEC_PRINCESS_1,
				}
				addDecoration(NewDecoration(tid, resIds, 1000))
			case '6':
				resIds := []graphic.ResourceID{
					graphic.RESOURCE_TYPE_DEC_PRINCESS_IS_WAITING,
				}
				addDecoration(NewDecoration(tid, resIds, 1000))
			case '7':
				resIds := []graphic.ResourceID{
					graphic.RESOURCE_TYPE_DEC_SUPER_MARIO_PAYPAL,
				}
				addDecoration(NewDecoration(tid, resIds, 1000))
			case 'a':
				resIds := []graphic.ResourceID{
					graphic.RESOURCE_TYPE_DEC_FAT_TREE_GREEN,
				}
				addDecoration(NewDecoration(tid, resIds, 1000))
			case 'b':
				resIds := []graphic.ResourceID{
					graphic.RESOURCE_TYPE_DEC_FAT_TREE_RED,
				}
				addDecoration(NewDecoration(tid, resIds, 1000))
			case 'c':
				resIds := []graphic.ResourceID{
					graphic.RESOURCE_TYPE_DEC_FAT_TREE_PINK,
				}
				addDecoration(NewDecoration(tid, resIds, 1000))
			case 'd':
				resIds := []graphic.ResourceID{
					graphic.RESOURCE_TYPE_DEC_FAT_TREE_WHITE,
				}
				addDecoration(NewDecoration(tid, resIds, 1000))
			case 'e':
				addTextDecoration(NewPrincessTextDecoration(tid))
			}
		}
	}

	if hero == nil {
		log.Fatal("no hero found when parsing level")
	}

	return &Level{
		Spec:        spec,
		BGRes:       bgRes,
		BGColor:     spec.BgColor,
		BGMusicID:   spec.BgMusicID,
		Decorations: decorations,
		TileObjects: tileObjs,
		Enemies:     enemies,
		Bullets:     list.New(),
		ObstMngr:    obstMngr,
		TheHero:     hero,
		InitHeroPos: vector.Pos{hero.levelRect.X, hero.levelRect.Y},
		NumTiles:    numTiles,
		effects:     list.New(),
	}
}

func ParseLevelSpec(levelFile string) *LevelSpec {
	conf, err := toml.LoadFile(levelFile)
	if err != nil {
		log.Fatalf("failed to load level file %s: %v", levelFile, err)
	}

	name := conf.Get("basic.name").(string)
	var nextLevelNames []string
	for _, name := range conf.Get("transfer.next-levels").([]interface{}) {
		nextLevelNames = append(nextLevelNames, name.(string))
	}

	bgFilename := conf.Get("graphic.bg-file").(string)

	bgMusicID, err := parseMusicID(conf.Get("music.bg-music-id").(string))
	if err != nil {
		log.Fatal(err)
	}

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
		BgMusicID:      bgMusicID,
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

func parseMusicID(str string) (audio.MusicID, error) {
	id, err := strconv.Atoi(str)
	if err != nil {
		return audio.MusicID(-1), errors.Wrap(err, "error parsing music ID")
	}
	return audio.MusicID(id), nil
}

package level_test

import (
	"testing"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/level"
	"github.com/zenja/mario/vector"
)

func TestGetTileID(t *testing.T) {
	TS := int32(graphic.TILE_SIZE)

	var expected vector.TileID
	var actual vector.TileID

	expected = vector.TileID{1, 2}
	actual = level.GetTileID(vector.Pos{TS * 2, TS * 3}, true, true)
	if expected != actual {
		t.Errorf("expected tile id %v but was %v", expected, actual)
	}

	expected = vector.TileID{2, 3}
	actual = level.GetTileID(vector.Pos{TS * 2, TS * 3}, false, false)
	if expected != actual {
		t.Errorf("expected tile id %v but was %v", expected, actual)
	}

	expected = vector.TileID{2, 3}
	actual = level.GetTileID(vector.Pos{TS*2 + 1, TS*3 + 1}, true, true)
	if expected != actual {
		t.Errorf("expected tile id %v but was %v", expected, actual)
	}
}

func TestGetTileRect(t *testing.T) {
	TS := int32(graphic.TILE_SIZE)

	var expected sdl.Rect
	var actual sdl.Rect

	expected = sdl.Rect{TS, TS * 2, TS, TS}
	actual = level.GetTileRect(vector.TileID{1, 2})
	if expected != actual {
		t.Errorf("expected tile rect %v but was %v", expected, actual)
	}
}

func TestGetSurroundingTileIDs(t *testing.T) {
	TS := int32(graphic.TILE_SIZE)

	var tiles [8]vector.TileID

	tiles = level.GetSurroundingTileIDs(sdl.Rect{TS * 3, TS * 5, TS, TS})
	if tiles[0] != (vector.TileID{3, 6}) {
		t.Errorf("expected tile id {3, 6} but was %v", tiles[0])
	}
	if tiles[1] != (vector.TileID{3, 4}) {
		t.Errorf("expected tile id {3, 4} but was %v", tiles[1])
	}
	if tiles[2] != (vector.TileID{2, 5}) {
		t.Errorf("expected tile id {2, 5} but was %v", tiles[2])
	}
	if tiles[3] != (vector.TileID{4, 5}) {
		t.Errorf("expected tile id {4, 5} but was %v", tiles[3])
	}
	if tiles[4] != (vector.TileID{2, 4}) {
		t.Errorf("expected tile id {2, 4} but was %v", tiles[4])
	}
	if tiles[5] != (vector.TileID{4, 4}) {
		t.Errorf("expected tile id {4, 4} but was %v", tiles[5])
	}
	if tiles[6] != (vector.TileID{2, 6}) {
		t.Errorf("expected tile id {2, 6} but was %v", tiles[6])
	}
	if tiles[7] != (vector.TileID{4, 6}) {
		t.Errorf("expected tile id {4, 6} but was %v", tiles[7])
	}

	tiles = level.GetSurroundingTileIDs(sdl.Rect{TS*3 + TS/5, TS*5 + TS/5, TS, TS})
	if tiles[0] != (vector.TileID{3, 6}) {
		t.Errorf("expected tile id {3, 6} but was %v", tiles[0])
	}
	if tiles[1] != (vector.TileID{3, 5}) {
		t.Errorf("expected tile id {3, 5} but was %v", tiles[1])
	}
	if tiles[2] != (vector.TileID{3, 5}) {
		t.Errorf("expected tile id {3, 5} but was %v", tiles[2])
	}
	if tiles[3] != (vector.TileID{4, 5}) {
		t.Errorf("expected tile id {4, 5} but was %v", tiles[3])
	}
	if tiles[4] != (vector.TileID{2, 5}) {
		t.Errorf("expected tile id {2, 5} but was %v", tiles[4])
	}
	if tiles[5] != (vector.TileID{4, 5}) {
		t.Errorf("expected tile id {4, 5} but was %v", tiles[5])
	}
	if tiles[6] != (vector.TileID{2, 6}) {
		t.Errorf("expected tile id {2, 6} but was %v", tiles[6])
	}
	if tiles[7] != (vector.TileID{4, 6}) {
		t.Errorf("expected tile id {4, 6} but was %v", tiles[7])
	}
}

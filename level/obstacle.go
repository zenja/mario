package level

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

// ObstacleManager know where is obstacle
// "Obst" means obstacle
type ObstacleManager struct {
	tilesInRow    int
	tilesInColumn int
	isObstTile    [][]bool
}

func NewObstacleManager(tilesInRow, tilesInColumn int) *ObstacleManager {
	var isObstTile [][]bool
	for i := 0; i < tilesInRow; i++ {
		var row []bool
		for j := 0; j < tilesInColumn; j++ {
			row = append(row, false)
		}
		isObstTile = append(isObstTile, row)
	}
	return &ObstacleManager{
		tilesInRow:    tilesInRow,
		tilesInColumn: tilesInColumn,
		isObstTile:    isObstTile,
	}
}

func (om *ObstacleManager) AddTileObst(tileID vector.TileID) {
	om.assertLegalTilePos(tileID)
	om.isObstTile[tileID.X][tileID.Y] = true
}

func (om *ObstacleManager) RemoveTileObst(tileID vector.TileID) {
	om.assertLegalTilePos(tileID)
	om.isObstTile[tileID.X][tileID.Y] = false
}

func (om *ObstacleManager) IsObstTile(tileID vector.TileID) bool {
	om.assertLegalTilePos(tileID)
	return om.isObstTile[tileID.X][tileID.Y]
}

// CalcCollisionSize returns how much a given rect collides obstacles in four directions: top, right, bottom, left
func (om *ObstacleManager) CalcCollisionSize(rect *sdl.Rect) (top, right, bottom, left int32) {
	leftTopPos := vector.Pos{rect.X, rect.Y}
	rightTopPos := vector.Pos{rect.X + rect.W, rect.Y}
	leftBottomPos := vector.Pos{rect.X, rect.Y + rect.H}
	rightBottomPos := vector.Pos{rect.X + rect.W, rect.Y + rect.H}

	//leftTopTildID := getTileID(leftTopPos)
	//rightTopTildID := getTileID(rightTopPos)
	//leftBottomTildID := getTileID(leftBottomPos)
	//rightBottomTildID := getTileID(rightBottomPos)
	//
	//for x := leftTopTildID.X; x <= rightTopTildID.X; x++ {
	//	for y := leftBottomTildID.Y; y <= rightBottomTildID.Y; y++ {
	//		if om.IsObstTile(vector.TileID{x, y}) {
	//		}
	//	}
	//}

	// check top edge
	var tmpLTPos = leftTopPos
	for ; tmpLTPos.X < rect.X+rect.W; tmpLTPos.X += graphic.TILE_SIZE {
		if om.IsObstTile(getTileID(tmpLTPos)) {
			tileRect := getTileRectByPos(tmpLTPos)
			top = tileRect.H - (rect.Y - tileRect.Y)
			left = tileRect.W - (rect.X - tileRect.X)
		}
	}
	// check left edge
	tmpLTPos = leftTopPos
	for ; tmpLTPos.Y < rect.Y+rect.H; tmpLTPos.Y += graphic.TILE_SIZE {
		if om.IsObstTile(getTileID(tmpLTPos)) {
			tileRect := getTileRectByPos(tmpLTPos)
			top = tileRect.H - (rect.Y - tileRect.Y)
			left = tileRect.W - (rect.X - tileRect.X)
		}
	}
	// check right edge
	tmpRTPos := rightTopPos
	for ; tmpRTPos.Y < rect.Y+rect.H; tmpRTPos.Y += graphic.TILE_SIZE {
		if om.IsObstTile(getTileID(tmpRTPos)) {
			tileRect := getTileRectByPos(tmpRTPos)
			top = tileRect.H - (rect.Y - tileRect.Y)
			right = rect.X + rect.W - tileRect.X
		}
	}
	// check bottom edge
	tmpLBPos := leftBottomPos
	for ; tmpLBPos.X < rect.X+rect.W; tmpLBPos.X += graphic.TILE_SIZE {
		if om.IsObstTile(getTileID(tmpLBPos)) {
			tileRect := getTileRectByPos(tmpLBPos)
			left = tileRect.W - (rect.X - tileRect.X)
			bottom = rect.Y + rect.H - tileRect.Y
		}
	}

	// check four corners
	if om.IsObstTile(getTileID(leftTopPos)) {
		tileRect := getTileRectByPos(leftTopPos)
		top = tileRect.H - (rect.Y - tileRect.Y)
		left = tileRect.W - (rect.X - tileRect.X)
	}
	if om.IsObstTile(getTileID(rightTopPos)) {
		tileRect := getTileRectByPos(rightTopPos)
		top = tileRect.H - (rect.Y - tileRect.Y)
		right = rect.X + rect.W - tileRect.X
	}
	if om.IsObstTile(getTileID(leftBottomPos)) {
		tileRect := getTileRectByPos(leftBottomPos)
		left = tileRect.W - (rect.X - tileRect.X)
		bottom = rect.Y + rect.H - tileRect.Y
	}
	if om.IsObstTile(getTileID(rightBottomPos)) {
		tileRect := getTileRectByPos(rightBottomPos)
		right = rect.X + rect.W - tileRect.X
		bottom = rect.Y + rect.H - tileRect.Y
	}
	return
}

func (om *ObstacleManager) assertLegalTilePos(tileID vector.TileID) {
	if tileID.X < 0 || tileID.Y < 0 {
		log.Fatalf("Tile position cannot be negative: (%d, %d)", tileID.X, tileID.Y)
	}
	if int(tileID.X) > om.tilesInRow {
		log.Fatalf("Tile X (%d) exceeds max width (%d)", tileID.X, om.tilesInRow)
	}
	if int(tileID.Y) > om.tilesInColumn {
		log.Fatalf("Tile Y (%d) exceeds max height (%d)", tileID.Y, om.tilesInColumn)
	}
}

// getTileID returns tile ID for a given position
// Note that getTileID won't check tile id nor given position
func getTileID(levelPos vector.Pos) vector.TileID {
	return vector.TileID{
		X: levelPos.X / graphic.TILE_SIZE,
		Y: levelPos.Y / graphic.TILE_SIZE,
	}
}

func getTileRectByPos(levelPos vector.Pos) *sdl.Rect {
	return getTileRect(getTileID(levelPos))
}

func getTileRect(tileID vector.TileID) *sdl.Rect {
	return &sdl.Rect{
		X: tileID.X * graphic.TILE_SIZE,
		Y: tileID.Y * graphic.TILE_SIZE,
		W: graphic.TILE_SIZE,
		H: graphic.TILE_SIZE,
	}
}

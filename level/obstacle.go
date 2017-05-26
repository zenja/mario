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
	if om.IsObstTile(getTileID(leftTopPos)) {
		tileRect := getTileRectByPos(leftTopPos)
		top = rect.Y - tileRect.Y
		left = rect.X - tileRect.X
	}
	if om.IsObstTile(getTileID(rightTopPos)) {
		tileRect := getTileRectByPos(rightTopPos)
		top = rect.Y - tileRect.Y
		right = (tileRect.X + tileRect.W) - (rect.X + rect.W)
	}
	if om.IsObstTile(getTileID(leftBottomPos)) {
		tileRect := getTileRectByPos(leftBottomPos)
		left = rect.X - tileRect.Y
		bottom = (tileRect.Y + tileRect.H) - (rect.Y + rect.H)
	}
	if om.IsObstTile(getTileID(rightBottomPos)) {
		tileRect := getTileRectByPos(rightBottomPos)
		right = (tileRect.X + tileRect.W) - (rect.X + rect.W)
		bottom = (tileRect.Y + tileRect.H) - (rect.Y + rect.H)
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

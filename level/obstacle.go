package level

import (
	"log"

	"fmt"

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

func (om *ObstacleManager) SolveCollision(desiredRect *sdl.Rect) (hitTop bool, hitRight bool, hitBottom bool, hitLeft bool) {
	tiles := GetSurroundingTileIDs(*desiredRect)
	for i, tid := range tiles {
		if tid.X < 0 || tid.Y < 0 {
			continue
		}

		if !om.IsObstTile(tid) {
			continue
		}

		interRect, isIntersect := desiredRect.Intersect(GetTileRect(tid))
		if !isIntersect {
			continue
		}

		switch i {
		case 0:
			desiredRect.Y -= interRect.H
			hitBottom = true
			fmt.Printf("i: %d, hit bottom! Y -= %d (interRect: %v) \n", i, interRect.H, interRect)
		case 1:
			desiredRect.Y += interRect.H
			fmt.Printf("i: %d, hit top! Y += %d (interRect: %v) \n", i, interRect.H, interRect)
			hitTop = true
		case 2:
			desiredRect.X += interRect.W
			hitLeft = true
		case 3:
			desiredRect.X -= interRect.W
			hitRight = true
		default:
			if interRect.W > interRect.H {
				// tile is diagonal, but resolving collision vertically
				if i > 5 {
					hitBottom = true
					desiredRect.Y -= interRect.H
					fmt.Printf("i: %d, hit bottom! Y += %d (interRect: %v) \n", i, interRect.H, interRect)
				} else {
					hitTop = true
					desiredRect.Y += interRect.H
					fmt.Printf("i: %d, hit top! Y += %d (interRect: %v) \n", i, interRect.H, interRect)
				}
			} else {
				// tile is diagonal, but resolving horizontally
				if i == 4 || i == 6 {
					hitLeft = true
					desiredRect.X += interRect.W
				} else {
					hitRight = true
					desiredRect.X -= interRect.W
				}
			}
		}
	}

	return
}

// GetSurroundingTileIDs returns the 8 surrounding tiles of a given rect
// The order is:
//
//     4     1     5
//         ______
//        |     |
//     2  |     |  3
//        |     |
//        |_____|
//     6     0     7
//
// NOTE that the tile id returned can be invalid, like negative X or Y
func GetSurroundingTileIDs(rect sdl.Rect) (tids [8]vector.TileID) {
	if rect.W >= graphic.TILE_SIZE*3 || rect.H >= graphic.TILE_SIZE*3 {
		log.Fatalf("rect width or height cannot exceed 3 * TILE_SIZE (%d) but is (%d, %d)",
			graphic.TILE_SIZE, rect.W, rect.H)
	}

	topMid := vector.Pos{rect.X + rect.W/2, rect.Y}
	bottomMid := vector.Pos{rect.X + rect.W/2, rect.Y + rect.H}
	leftMid := vector.Pos{rect.X, rect.Y + rect.H/2}
	rightMid := vector.Pos{rect.X + rect.W, rect.Y + rect.H/2}

	topMidTID := GetTileID(topMid, true, true)
	bottomMidTID := GetTileID(bottomMid, false, true)
	leftMidTID := GetTileID(leftMid, true, true)
	rightMidTID := GetTileID(rightMid, true, false)

	tids[0] = bottomMidTID
	tids[1] = topMidTID
	tids[2] = leftMidTID
	tids[3] = rightMidTID

	tids[4] = vector.TileID{topMidTID.X - 1, topMidTID.Y}
	tids[5] = vector.TileID{topMidTID.X + 1, topMidTID.Y}
	tids[6] = vector.TileID{bottomMidTID.X - 1, bottomMidTID.Y}
	tids[7] = vector.TileID{bottomMidTID.X + 1, bottomMidTID.Y}

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

// GetTileID returns tile ID for a given position
// Note that GetTileID won't check tile id nor given position
func GetTileID(levelPos vector.Pos, preferTop bool, preferLeft bool) vector.TileID {
	x := levelPos.X / graphic.TILE_SIZE
	if levelPos.X%graphic.TILE_SIZE == 0 && preferLeft {
		x--
	}
	y := levelPos.Y / graphic.TILE_SIZE
	if levelPos.Y%graphic.TILE_SIZE == 0 && preferTop {
		y--
	}
	return vector.TileID{
		X: x,
		Y: y,
	}
}

func GetTileRect(tileID vector.TileID) *sdl.Rect {
	return &sdl.Rect{
		X: tileID.X * graphic.TILE_SIZE,
		Y: tileID.Y * graphic.TILE_SIZE,
		W: graphic.TILE_SIZE,
		H: graphic.TILE_SIZE,
	}
}

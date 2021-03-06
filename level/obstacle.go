package level

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/vector"
)

type obstType uint8

const (
	not_obst obstType = iota // default: no obst
	normal_obst
	enemy_only_obst
	up_thru_obst // obst that can pass through from its bottom
)

type SolveCollisionType uint8

const (
	SOLVE_COLLISION_NORMAL SolveCollisionType = iota
	// also consider enemy-only obst as obst
	SOLVE_COLLISION_ENEMY
)

// ObstacleManager know where is obstacle
// "Obst" means obstacle
type ObstacleManager struct {
	tilesInRow    int
	tilesInColumn int

	// obstType[tilesInRow][tilesInColumn], so that we can use obsts[TID.X][TID.Y]
	// so its shape is a rotation of level's
	obsts [][]obstType
}

func NewObstacleManager(tilesInRow, tilesInColumn int) *ObstacleManager {
	var obsts [][]obstType
	for i := 0; i < tilesInRow; i++ {
		row := make([]obstType, tilesInColumn)
		obsts = append(obsts, row)
	}
	return &ObstacleManager{
		tilesInRow:    tilesInRow,
		tilesInColumn: tilesInColumn,
		obsts:         obsts,
	}
}

func (om *ObstacleManager) AddNormalTileObst(tileID vector.TileID) {
	om.assertLegalTilePos(tileID)
	om.obsts[tileID.X][tileID.Y] = normal_obst
}

func (om *ObstacleManager) AddEnemyOnlyTileObst(tileID vector.TileID) {
	om.assertLegalTilePos(tileID)
	om.obsts[tileID.X][tileID.Y] = enemy_only_obst
}

func (om *ObstacleManager) AddUpThruTileObst(tileID vector.TileID) {
	om.assertLegalTilePos(tileID)
	om.obsts[tileID.X][tileID.Y] = up_thru_obst
}

func (om *ObstacleManager) RemoveTileObst(tileID vector.TileID) {
	om.assertLegalTilePos(tileID)
	om.obsts[tileID.X][tileID.Y] = not_obst
}

func (om *ObstacleManager) SolveCollision(desiredRect *sdl.Rect, sctype SolveCollisionType) (
	hitTop bool,
	hitRight bool,
	hitBottom bool,
	hitLeft bool,
	tilesHit []vector.TileID) {

	tiles := GetSurroundingTileIDs(*desiredRect)
	for i, tid := range tiles {
		if tid.X < 0 || tid.Y < 0 {
			continue
		}

		if !om.isObstTile(tid, *desiredRect, sctype) {
			continue
		}

		tileRect := GetTileRect(tid)
		interRect, isIntersect := desiredRect.Intersect(&tileRect)
		if !isIntersect {
			continue
		}

		tilesHit = append(tilesHit, tid)

		switch i {
		case 0:
			desiredRect.Y -= interRect.H
			hitBottom = true
			//log.Printf("i: %d, hit bottom! Y -= %d (interRect: %v) \n", i, interRect.H, interRect)
		case 1:
			desiredRect.Y += interRect.H
			//log.Printf("i: %d, hit top! Y += %d (interRect: %v) \n", i, interRect.H, interRect)
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
					//log.Printf("i: %d, hit bottom! Y += %d (interRect: %v) \n", i, interRect.H, interRect)
				} else {
					hitTop = true
					desiredRect.Y += interRect.H
					//log.Printf("i: %d, hit top! Y += %d (interRect: %v) \n", i, interRect.H, interRect)
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

// TODO return value instead of pointer
func GetTileRect(tileID vector.TileID) sdl.Rect {
	return sdl.Rect{
		X: tileID.X * graphic.TILE_SIZE,
		Y: tileID.Y * graphic.TILE_SIZE,
		W: graphic.TILE_SIZE,
		H: graphic.TILE_SIZE,
	}
}

func GetTileStartPos(tileID vector.TileID) vector.Pos {
	rect := GetTileRect(tileID)
	return vector.Pos{rect.X, rect.Y}
}

func GetRectStartPos(rect sdl.Rect) vector.Vec2D {
	return vector.Vec2D{rect.X, rect.Y}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Helper methods
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (om *ObstacleManager) isLegalTilePos(tileID vector.TileID) bool {
	if tileID.X < 0 || tileID.Y < 0 {
		return false
	}
	if int(tileID.X) >= om.tilesInRow {
		return false
	}
	if int(tileID.Y) >= om.tilesInColumn {
		return false
	}
	return true
}

func (om *ObstacleManager) assertLegalTilePos(tileID vector.TileID) {
	if tileID.X < 0 || tileID.Y < 0 {
		log.Fatalf("Tile position cannot be negative: (%d, %d)", tileID.X, tileID.Y)
	}
	if int(tileID.X) >= om.tilesInRow {
		log.Fatalf("Tile X (%d) exceeds max width (%d)", tileID.X, om.tilesInRow)
	}
	if int(tileID.Y) >= om.tilesInColumn {
		log.Fatalf("Tile Y (%d) exceeds max height (%d)", tileID.Y, om.tilesInColumn)
	}
}

func (om *ObstacleManager) isObstTile(tileID vector.TileID, desiredRect sdl.Rect, sctype SolveCollisionType) bool {
	if !om.isLegalTilePos(tileID) {
		// all tiles out of scope are considered not obstacles
		// so objects can actually update itself to go out of scope, and it is easy to detect this
		return false
	}

	obstType := om.obsts[tileID.X][tileID.Y]

	if obstType == normal_obst {
		return true
	}

	if obstType == enemy_only_obst && sctype == SOLVE_COLLISION_ENEMY {
		return true
	}

	if obstType == up_thru_obst && GetTileRect(tileID).Y+graphic.TILE_SIZE/2 >= desiredRect.Y+desiredRect.H {
		return true
	}

	return false
}

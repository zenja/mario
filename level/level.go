package level

import (
	"github.com/zenja/mario/graphic"
	"github.com/zenja/mario/object"
)

type Level struct {
	Objects [][]object.Object
}

func ParseLevel(arr [][]byte, tileRegistry map[graphic.TileID]*graphic.Tile) *Level {
	var objRows [][]object.Object
	var currentX, currentY int32
	for i, arrRow := range arr {
		currentX = 0
		var objRow []object.Object
		for j, _ := range arrRow {
			switch arr[i][j] {
			// Ground
			case 'G':
				objRow = append(objRow, object.NewSingleTileObject(graphic.TILE_TYPE_GROUD, currentX, currentY))
				//fmt.Printf("%d, %d\n", currentX, currentY) fixme
			}
			currentX += graphic.TILE_SIZE
		}
		objRows = append(objRows, objRow)
		currentY += graphic.TILE_SIZE
	}
	return &Level{Objects: objRows}
}

func (l *Level) Draw(g *graphic.Graphic, xCamStart, yCamStart int32) {
	for i := range l.Objects {
		for j := range l.Objects[i] {
			l.Objects[i][j].Draw(g, xCamStart, yCamStart)
		}
	}
}

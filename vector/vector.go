package vector

type Vec2D struct {
	X, Y int32
}

func (v *Vec2D) Add(other Vec2D) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vec2D) Multiply(m int32) {
	v.X *= m
	v.Y *= m
}

func (v *Vec2D) Divide(m int32) {
	v.X /= m
	v.Y /= m
}

// Vec2D aliases
type Pos Vec2D
type TileID Vec2D // (0, 0) is left top tile

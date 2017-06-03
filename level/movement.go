package level

import (
	"log"

	mutils "github.com/zenja/mario/math_utils"
	"github.com/zenja/mario/vector"
)

func CalcVelocityStep(velocity vector.Vec2D, currTicks uint32, lastTicks uint32, maxVel *vector.Vec2D) vector.Vec2D {
	// calculate velocity step
	velocityStep := velocity
	velocityStep.Multiply(int32(currTicks - lastTicks))
	velocityStep.Divide(1000)

	// limit max velocity step if there is a limit (non-nil maxVel)
	if maxVel != nil {
		if mutils.Abs(velocityStep.X) > maxVel.X {
			log.Printf("warning: velocity step's |X| is %d > %d, limited to %d", velocityStep.X, maxVel.X, maxVel.X)
			if velocityStep.X > 0 {
				velocityStep.X = maxVel.X
			} else {
				velocityStep.X = -maxVel.X
			}
		}
		if mutils.Abs(velocityStep.Y) > maxVel.Y {
			log.Printf("warning: velocity step's |Y| is %d > %d, limited to %d", velocityStep.Y, maxVel.Y, maxVel.Y)
			if velocityStep.Y > 0 {
				velocityStep.Y = maxVel.Y
			} else {
				velocityStep.Y = -maxVel.Y
			}
		}
	}

	return velocityStep
}

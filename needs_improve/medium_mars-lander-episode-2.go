package main

import (
	"fmt"
	"math"
)

const (
	maxHorizontalSpeed = 30
	maxVerticalSpeed   = 38
	maxRotation        = 60
)

var (
	maxTerrainHeight   int
	startFlat, endFlat int
	flatHeight         int
)

func main() {
	// surfaceN: the number of points used to draw the surface of Mars.
	var surfaceN int
	fmt.Scan(&surfaceN)

	lastTerrainY := math.MaxInt32
	var lastTerrainX int
	for i := 0; i < surfaceN; i++ {
		// landX: X coordinate of a surface point. (0 to 6999)
		// landY: Y coordinate of a surface point. By linking all the points together in a sequential fashion, you form the surface of Mars.
		var landX, landY int
		fmt.Scan(&landX, &landY)

		if landY > maxTerrainHeight {
			maxTerrainHeight = landY
		}

		if landY == lastTerrainY {
			startFlat = lastTerrainX
			endFlat = landX
			flatHeight = landY
		}

		lastTerrainX = landX
		lastTerrainY = landY
	}

	for {
		// hSpeed: the horizontal speed (in m/s), can be negative.
		// vSpeed: the vertical speed (in m/s), can be negative.
		// fuel: the quantity of remaining fuel in liters.
		// rotate: the rotation angle in degrees (-90 to 90).
		// power: the thrust power (0 to 4).
		var X, Y, hSpeed, vSpeed, fuel, rotate, power int
		fmt.Scan(&X, &Y, &hSpeed, &vSpeed, &fuel, &rotate, &power)

		overFlatTerrain := X > startFlat && X < endFlat

		rotation := 0
		power = 2
		desiredHorizontal := 0
		desiredVertical := -1 * maxVerticalSpeed

		if X < startFlat {
			desiredHorizontal = maxHorizontalSpeed
		} else if X > endFlat {
			desiredHorizontal = -1 * maxHorizontalSpeed
		}

		horizontalDelta := desiredHorizontal - hSpeed
		rotation = int(-1.0 * float64(maxRotation) * math.Min(1.0, (float64(horizontalDelta)/float64(maxHorizontalSpeed))))

		if !overFlatTerrain {
			if Y < maxTerrainHeight+200 {
				desiredVertical = 5
				desiredHorizontal = 0
			} else {
				desiredVertical = -15
			}
		} else if math.Abs(float64(rotation)) > 0.0 {
			desiredVertical = -5
		}

		if vSpeed < desiredVertical {
			power = 4
			rotationModifier := math.Max(0.4, 1.0/(float64(vSpeed)/float64(desiredVertical)))
			rotation = int(float64(rotation) * rotationModifier)
		}

		if Y < flatHeight+50 {
			rotation = 0
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")

		// rotate power. rotate is the desired rotation angle. power is the desired thrust power.
		fmt.Printf("%d %d\n", rotation, power)
	}
}

package main

import (
	"fmt"
	"strings"
)

type bounds struct {
	X1, Y1, X2, Y2 int
}

var (
	searchArea     bounds
	myPosX, myPosY int
)

func main() {
	// W: width of the building.
	// H: height of the building.
	var W, H int
	fmt.Scan(&W, &H)

	searchArea = bounds{0, 0, W, H}

	// N: maximum number of turns before game over.
	var N int
	fmt.Scan(&N)

	fmt.Scan(&myPosX, &myPosY)

	for {
		// bombDir: the direction of the bombs from batman's current location (U, UR, R, DR, D, DL, L or UL)
		var bombDir string
		fmt.Scan(&bombDir)

		if strings.Contains(bombDir, "U") {
			searchArea.Y2 = myPosY
		}

		if strings.Contains(bombDir, "D") {
			searchArea.Y1 = myPosY
		}

		if strings.Contains(bombDir, "R") {
			searchArea.X1 = myPosX
		}

		if strings.Contains(bombDir, "L") {
			searchArea.X2 = myPosX
		}

		myPosX = searchArea.X1 + (searchArea.X2-searchArea.X1)/2
		myPosY = searchArea.Y1 + (searchArea.Y2-searchArea.Y1)/2

		// the location of the next window Batman should jump to.
		fmt.Printf("%d %d\n",
			myPosX,
			myPosY,
		)
	}
}

package main

import "fmt"

var (
	elevators       = map[int]int{}
	blockersInPlace = map[int]bool{}
)

func main() {
	// nbFloors: number of floors
	// width: width of the area
	// nbRounds: maximum number of rounds
	// exitFloor: floor on which the exit is found
	// exitPos: position of the exit on its floor
	// nbTotalClones: number of generated clones
	// nbAdditionalElevators: ignore (always zero)
	// nbElevators: number of elevators
	var nbFloors, width, nbRounds, exitFloor, exitPos, nbTotalClones, nbAdditionalElevators, nbElevators int
	fmt.Scan(&nbFloors, &width, &nbRounds, &exitFloor, &exitPos, &nbTotalClones, &nbAdditionalElevators, &nbElevators)

	for i := 0; i < nbElevators; i++ {
		// elevatorFloor: floor on which this elevator is found
		// elevatorPos: position of the elevator on its floor
		var elevatorFloor, elevatorPos int
		fmt.Scan(&elevatorFloor, &elevatorPos)
		elevators[elevatorFloor] = elevatorPos
	}
	for {
		// cloneFloor: floor of the leading clone
		// clonePos: position of the leading clone on its floor
		// direction: direction of the leading clone: LEFT or RIGHT
		var cloneFloor, clonePos int
		var direction string
		fmt.Scan(&cloneFloor, &clonePos, &direction)

		switch {
		case clonePos == 0 || clonePos == width-1:
			fmt.Println("BLOCK")
			blockersInPlace[cloneFloor] = true
		case direction == "RIGHT" && elevators[cloneFloor] < clonePos && !blockersInPlace[cloneFloor]:
			fmt.Println("BLOCK")
			blockersInPlace[cloneFloor] = true
		case direction == "LEFT" && elevators[cloneFloor] > clonePos && !blockersInPlace[cloneFloor]:
			fmt.Println("BLOCK")
			blockersInPlace[cloneFloor] = true
		case direction == "RIGHT" && exitPos < clonePos && exitFloor == cloneFloor && !blockersInPlace[cloneFloor]:
			fmt.Println("BLOCK")
			blockersInPlace[cloneFloor] = true
		case direction == "LEFT" && exitPos > clonePos && exitFloor == cloneFloor && !blockersInPlace[cloneFloor]:
			fmt.Println("BLOCK")
			blockersInPlace[cloneFloor] = true
		default:
			fmt.Println("WAIT")
		}

	}
}

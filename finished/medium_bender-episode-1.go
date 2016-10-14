package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type direction string

const (
	SOUTH direction = "SOUTH"
	WEST            = "WEST"
	EAST            = "EAST"
	NORTH           = "NORTH"
)

const (
	MAP_TILE_OBSTACLE      byte = 'X'
	MAP_TILE_WALL               = '#'
	MAP_TILE_EAST               = 'E'
	MAP_TILE_WEST               = 'W'
	MAP_TILE_NORTH              = 'N'
	MAP_TILE_SOUTH              = 'S'
	MAP_TILE_BEER               = 'B'
	MAP_TILE_INVERT             = 'I'
	MAP_TILE_START              = '@'
	MAP_TILE_TELEPORTER         = 'T'
	MAP_TILE_FLOOR              = ' '
	MAP_TILE_DEATH              = '$'
	MAP_TILE_OUT_OF_BOUNDS      = '?' // Internal error-reporter
)

var directionChange = map[bool][]direction{
	false: []direction{SOUTH, EAST, NORTH, WEST},
	true:  []direction{WEST, NORTH, EAST, SOUTH},
}

type point struct{ X, Y int }

func (p point) Next(d direction) point {
	switch d {
	case SOUTH:
		return point{p.X, p.Y + 1}
	case WEST:
		return point{p.X - 1, p.Y}
	case NORTH:
		return point{p.X, p.Y - 1}
	case EAST:
		return point{p.X + 1, p.Y}
	}
	return p
}

type playground struct {
	Source        string
	Width, Height int
}

func (p playground) Get(pt point) byte {
	pos := pt.Y*p.Width + pt.X
	if pos >= len(p.Source) {
		return MAP_TILE_OUT_OF_BOUNDS
	}
	return p.Source[pos]
}

func (p *playground) Set(pt point, mapReplace byte) {
	pos := pt.Y*p.Width + pt.X
	tmp := []byte(p.Source)
	tmp[pos] = mapReplace
	p.Source = string(tmp)
}

func (p playground) FindUniquePOI(poi byte) point {
	pos := strings.IndexByte(p.Source, poi)
	return point{X: pos % p.Width, Y: pos / p.Width}
}

func (p playground) FindMultiplePOI(poi byte) []point {
	out := []point{}

	for pos := 0; pos < len(p.Source); pos++ {
		if p.Source[pos] == poi {
			out = append(out, point{X: pos % p.Width, Y: pos / p.Width})
		}
	}

	return out
}

type bender struct {
	City playground

	Position      point
	BeerMode      bool
	InvertMode    bool
	MoveDirection direction

	path []string
}

func (b *bender) Init() {
	b.Position = b.City.FindUniquePOI(MAP_TILE_START)
	b.MoveDirection = SOUTH
}

func (b bender) Path() string {
	return strings.Join(b.path, "\n")
}

func (b *bender) Trace(maxSteps int) error {
	madeSteps := 0

	for madeSteps < maxSteps {
		nextPos := b.Position.Next(b.MoveDirection)

		nd := 0
		for b.City.Get(nextPos) == MAP_TILE_OBSTACLE || b.City.Get(nextPos) == MAP_TILE_WALL {
			if b.BeerMode && b.City.Get(nextPos) == MAP_TILE_OBSTACLE {
				b.City.Set(nextPos, MAP_TILE_FLOOR)
			} else {
				b.MoveDirection = directionChange[b.InvertMode][nd]
				nextPos = b.Position.Next(b.MoveDirection)
				nd++
			}
		}

		b.Position = nextPos
		log.Printf("Moved to %#v", b.Position)
		b.path = append(b.path, string(b.MoveDirection))

		switch b.City.Get(b.Position) {
		case MAP_TILE_START:
			return errors.New("Returned to start, no cool.")
		case MAP_TILE_OUT_OF_BOUNDS:
			return errors.New("Bender left the map. How?!?")
		case MAP_TILE_WEST:
			b.MoveDirection = WEST
		case MAP_TILE_EAST:
			b.MoveDirection = EAST
		case MAP_TILE_SOUTH:
			b.MoveDirection = SOUTH
		case MAP_TILE_NORTH:
			b.MoveDirection = NORTH
		case MAP_TILE_BEER:
			b.BeerMode = !b.BeerMode
		case MAP_TILE_INVERT:
			b.InvertMode = !b.InvertMode
		case MAP_TILE_DEATH:
			// We got killed. Bye cruel world!
			return nil
		case MAP_TILE_TELEPORTER:
			teleporters := b.City.FindMultiplePOI(MAP_TILE_TELEPORTER)
			for i := range teleporters {
				t := teleporters[i]
				if t.X == b.Position.X && t.Y == b.Position.Y {
					continue
				}
				b.Position = t
				log.Printf("Bzzzzp. Now at: %#v", b.Position)
				break
			}
		}

		madeSteps++
	}

	return errors.New("Max steps reached. (LOOP)")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var L, C int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &L, &C)

	b := &bender{}
	b.City.Width = C
	b.City.Height = L

	for i := 0; i < L; i++ {
		scanner.Scan()
		b.City.Source = b.City.Source + scanner.Text()
	}

	b.Init()

	log.Printf("%#v", b)

	if err := b.Trace(L * C); err != nil {
		log.Printf("Trace error: %s", err)
		fmt.Println("LOOP")
	} else {
		fmt.Println(b.Path())
	}
}

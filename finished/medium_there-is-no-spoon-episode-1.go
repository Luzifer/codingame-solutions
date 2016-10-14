package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	CHAR_NODE = '0'
	CHAR_FREE = '.'
)

type playground struct {
	Source        string
	Width, Height int
}

func (p playground) Scan() <-chan string {
	log.Printf("%#v", p)
	c := make(chan string, 100)

	go func() {
		defer close(c)
		wg := sync.WaitGroup{}
		for i := 0; i < len(p.Source); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				x := i % p.Width
				y := i / p.Width
				if p.get(x, y) != CHAR_NODE {
					return
				}
				log.Printf("X: %d Y: %d", x, y)

				c <- p.GetNeighborNotation(x, y)
			}(i)
		}
		wg.Wait()
	}()

	return c
}

func (p playground) get(x, y int) byte {
	pos := y*p.Width + x
	log.Printf("Getting x=%d y=%d (str=%d, char=%c", x, y, pos, p.Source[pos])
	if pos >= len(p.Source) {
		return '?'
	}
	return p.Source[pos]
}

func (p playground) GetNeighborNotation(x, y int) string {
	bx, by := p.getBottomNext(x, y+1)
	rx, ry := p.getRightNext(x+1, y)
	return fmt.Sprintf("%d %d %d %d %d %d", x, y, rx, ry, bx, by)
}

func (p playground) getBottomNext(x, y int) (int, int) {
	if y >= p.Height {
		return -1, -1
	}

	if p.get(x, y) == CHAR_NODE {
		return x, y
	}

	return p.getBottomNext(x, y+1)
}

func (p playground) getRightNext(x, y int) (int, int) {
	if x >= p.Width {
		return -1, -1
	}

	if p.get(x, y) == CHAR_NODE {
		return x, y
	}

	return p.getRightNext(x+1, y)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// width: the number of cells on the X axis
	var width int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &width)

	// height: the number of cells on the Y axis
	var height int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &height)

	pg := playground{Width: width, Height: height}

	for i := 0; i < height; i++ {
		scanner.Scan()
		pg.Source = pg.Source + scanner.Text()
	}

	for out := range pg.Scan() {
		fmt.Println(out)
	}

}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type calculationResult struct{ Seated, Next int }

var (
	groups           [1000]int
	calculationCache = map[int]calculationResult{}
)

func main() {
	defer timeTrack(time.Now(), "main")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)
	scanner.Split(bufio.ScanWords)

	var L, C, N, i int
	scanner.Scan()
	L = toInt(scanner.Bytes())
	scanner.Scan()
	C = toInt(scanner.Bytes())
	scanner.Scan()
	N = toInt(scanner.Bytes())

	log.Printf("C: %d", C)

	for ; i < N; i++ {
		scanner.Scan()
		groups[i] = toInt(scanner.Bytes())
	}

	var earned, next, seated, j int
	for i = 0; i < C; i++ {
		key := next
		seated = 0

		if cr, ok := calculationCache[key]; ok {
			seated = cr.Seated
			next = cr.Next
		} else {
			j = 0
			for j < N && groups[next] <= L-seated {
				seated += groups[next]
				j++
				next++
				if next == N {
					next = 0
				}
			}
			calculationCache[key] = calculationResult{Seated: seated, Next: next}
		}

		earned += seated
	}

	fmt.Printf("%d\n", earned)
}

func toInt(buf []byte) (n int) {
	for _, v := range buf {
		n = n*10 + int(v-'0')
	}
	return
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

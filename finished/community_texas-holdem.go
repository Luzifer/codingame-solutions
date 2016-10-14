package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	WINNING_UNKNOWN uint = iota
	WINNING_HIGH_CARD
	WINNING_PAIR
	WINNING_TWO_PAIR
	WINNING_THREE_OF_A_KIND
	WINNING_STRAIGHT
	WINNING_FLUSH
	WINNING_FULL_HOUSE
	WINNING_FOUR_OF_A_KIND
	WINNING_STRAIGHT_FLUSH
)

var (
	winningOutputs = map[uint]string{
		WINNING_HIGH_CARD:       "HIGH_CARD",
		WINNING_PAIR:            "PAIR",
		WINNING_TWO_PAIR:        "TWO_PAIR",
		WINNING_THREE_OF_A_KIND: "THREE_OF_A_KIND",
		WINNING_STRAIGHT:        "STRAIGHT",
		WINNING_FLUSH:           "FLUSH",
		WINNING_FULL_HOUSE:      "FULL_HOUSE",
		WINNING_FOUR_OF_A_KIND:  "FOUR_OF_A_KIND",
		WINNING_STRAIGHT_FLUSH:  "STRAIGHT_FLUSH",
	}
	handChecks = map[uint]handCheck{
		WINNING_HIGH_CARD:       checkHighCard,
		WINNING_PAIR:            checkPair,
		WINNING_TWO_PAIR:        checkTwoPair,
		WINNING_THREE_OF_A_KIND: checkThreeOfAKind,
		WINNING_STRAIGHT:        checkStraight,
		WINNING_FLUSH:           checkFlush,
		WINNING_FULL_HOUSE:      checkFullHouse,
		WINNING_FOUR_OF_A_KIND:  checkFourOfAKind,
		WINNING_STRAIGHT_FLUSH:  checkStraightFlush,
	}
	cardValues = map[byte]int{'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'J': 11, 'Q': 12, 'K': 13, 'A': 14}
)

type card struct{ Value, Suit byte }
type handCheck func(hand) (bool, string)
type hand []card

func parseHand(player, board string) hand {
	all := strings.Split(player, " ")
	all = append(all, strings.Split(board, " ")...)

	out := hand{}
	for i := 0; i < 7; i++ {
		out = append(out, card{all[i][0], all[i][1]})
	}

	return out
}

type cardValueList []byte

func (c cardValueList) Len() int           { return len(c) }
func (c cardValueList) Less(i, j int) bool { return cardValues[c[i]] > cardValues[c[j]] }
func (c cardValueList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	scanner.Scan()
	holeCardsPlayer1 := scanner.Text()
	scanner.Scan()
	holeCardsPlayer2 := scanner.Text()
	scanner.Scan()
	communityCards := scanner.Text()

	fmt.Fprintln(os.Stderr, "+-------- DBG -------+")
	fmt.Fprintf(os.Stderr, "| P1: %s          |\n| P2: %s          |\n| CO: %s |\n", holeCardsPlayer1, holeCardsPlayer2, communityCards)
	fmt.Fprintln(os.Stderr, "+-------- DBG -------+\n")

	fmt.Println(executeChecks(holeCardsPlayer1, holeCardsPlayer2, communityCards))
}

func executeChecks(holeCardsPlayer1, holeCardsPlayer2, communityCards string) string {
	player1Hand := parseHand(holeCardsPlayer1, communityCards)
	player2Hand := parseHand(holeCardsPlayer2, communityCards)

	player1Result := WINNING_UNKNOWN
	player2Result := WINNING_UNKNOWN

	var player1Out, player2Out string

	for w, hc := range handChecks {
		if ok, out := hc(player1Hand); ok && w > player1Result {
			player1Result = w
			player1Out = out
		}

		if ok, out := hc(player2Hand); ok && w > player2Result {
			player2Result = w
			player2Out = out
		}
	}

	switch {
	case player1Result > player2Result:
		return fmt.Sprintf("1 %s %s", winningOutputs[player1Result], player1Out)
	case player2Result > player1Result:
		return fmt.Sprintf("2 %s %s", winningOutputs[player2Result], player2Out)
	case player2Result == player1Result:
		switch compareResultOutputHighCard(player1Out, player2Out) {
		case 1:
			return fmt.Sprintf("1 %s %s", winningOutputs[player1Result], player1Out)
		case 2:
			return fmt.Sprintf("2 %s %s", winningOutputs[player2Result], player2Out)
		}
	}

	return "DRAW"
}

func compareResultOutputHighCard(player1, player2 string) int {
	for i := 0; i < 5; i++ {
		switch {
		case cardValues[player1[i]] > cardValues[player2[i]]:
			return 1
		case cardValues[player1[i]] < cardValues[player2[i]]:
			return 2
		}
	}

	return 0
}

func (h hand) OrderCards() []byte {
	cvs := []byte{}
	for _, c := range h {
		cvs = append(cvs, c.Value)
	}
	sort.Sort(cardValueList(cvs))
	return cvs
}

func (h hand) CountValues() map[byte]int {
	out := map[byte]int{}

	for _, c := range h {
		if _, ok := out[c.Value]; !ok {
			out[c.Value] = 0
		}
		out[c.Value]++
	}

	return out
}

func (h hand) CountSuits() map[byte]int {
	out := map[byte]int{}

	for _, c := range h {
		if _, ok := out[c.Suit]; !ok {
			out[c.Suit] = 0
		}
		out[c.Suit]++
	}

	return out
}

func cardOutput(cardValues []byte, winning []byte) string {
	kickers := []byte{}

	for _, c := range cardValues {
		found := false
		for _, w := range winning {
			if w == c {
				found = true
			}
		}
		if !found {
			kickers = append(kickers, c)
		}
	}

	sort.Sort(cardValueList(kickers))

	return strings.Join(cardValuesToList(winning...), "") + strings.Join(cardValuesToList(kickers[0:5-len(winning)]...), "")
}

func cardValuesToList(vals ...byte) []string {
	out := []string{}
	for _, b := range vals {
		out = append(out, string(b))
	}
	return out
}

func checkHighCard(cards hand) (bool, string) { return true, cardOutput(cards.OrderCards(), []byte{}) }

func checkPair(cards hand) (bool, string) {
	for val, count := range cards.CountValues() {
		if count == 2 {
			return true, cardOutput(cards.OrderCards(), []byte{val, val})
		}
	}

	return false, ""
}

func checkTwoPair(cards hand) (bool, string) {
	var foundPairs = []byte{}

	for val, count := range cards.CountValues() {
		if count == 2 {
			foundPairs = append(foundPairs, val)
		}
	}

	if len(foundPairs) < 2 {
		return false, ""
	}

	sort.Sort(cardValueList(foundPairs))

	pairs := []byte{foundPairs[0], foundPairs[0], foundPairs[1], foundPairs[1]}
	return true, cardOutput(cards.OrderCards(), pairs)

}

func checkThreeOfAKind(cards hand) (bool, string) {
	for val, count := range cards.CountValues() {
		if count == 3 {
			return true, cardOutput(cards.OrderCards(), []byte{val, val, val})
		}
	}

	return false, ""
}

func checkStraight(cards hand) (bool, string) {
	order := cards.OrderCards()

	if order[0] == 'A' { // Special case: Despite having value 14 an Ace can have value 1
		order = append(order, 'A')
	}

	straightLength := 1
	for i := 1; i < len(order); i++ {
		if cardValues[order[i-1]] == cardValues[order[i]]+1 || (cardValues[order[i-1]] == 2 && order[i] == 'A') {
			straightLength++
		} else {
			straightLength = 1
		}

		if straightLength == 5 {
			return true, cardOutput(cards.OrderCards(), order[i-4:i+1])
		}
	}

	return false, ""
}

func checkFlush(cards hand) (bool, string) {
	var fiver byte

	for suit, count := range cards.CountSuits() {
		if count >= 5 {
			fiver = suit
		}
	}

	if fiver == 0x0 {
		return false, ""
	}

	matchedSuitValues := []byte{}
	for _, c := range cards {
		if c.Suit == fiver {
			matchedSuitValues = append(matchedSuitValues, c.Value)
		}
	}
	sort.Sort(cardValueList(matchedSuitValues))

	return true, cardOutput(cards.OrderCards(), matchedSuitValues[0:5])
}

func checkFullHouse(cards hand) (bool, string) {
	var (
		foundPairs        = []byte{}
		foundThreeOfAKind = []byte{}
	)

	for val, count := range cards.CountValues() {
		if count == 2 {
			foundPairs = append(foundPairs, val)
		} else if count == 3 {
			foundThreeOfAKind = append(foundThreeOfAKind, val)
		}
	}

	if len(foundPairs) == 0 || len(foundThreeOfAKind) == 0 {
		return false, ""
	}

	sort.Sort(cardValueList(foundPairs))

	winning := []byte{foundThreeOfAKind[0], foundThreeOfAKind[0], foundThreeOfAKind[0], foundPairs[0], foundPairs[0]}
	return true, cardOutput(cards.OrderCards(), winning)
}

func checkFourOfAKind(cards hand) (bool, string) {
	for val, count := range cards.CountValues() {
		if count == 4 {
			return true, cardOutput(cards.OrderCards(), []byte{val, val, val, val})
		}
	}

	return false, ""
}

func checkStraightFlush(cards hand) (bool, string) {
	var fiver byte

	for suit, count := range cards.CountSuits() {
		if count >= 5 {
			fiver = suit
		}
	}

	if fiver == 0x0 {
		return false, ""
	}

	matchedCards := hand{}
	for i := range cards {
		c := cards[i]
		if c.Suit == fiver {
			matchedCards = append(matchedCards, c)
		}
	}

	return checkStraight(matchedCards)
}

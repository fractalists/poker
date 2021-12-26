package src

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScoreWithFourOfAKind(t *testing.T) {
	card1 := Card{Suit: DIAMONDS, Rank: JACK}
	card2 := Card{Suit: CLUBS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: KING}
	card4 := Card{Suit: HEARTS, Rank: ACE}
	card5 := Card{Suit: SPADES, Rank: ACE}
	card6 := Card{Suit: DIAMONDS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: ACE}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

	handType, fiveCards, score := Score(cards)
	fmt.Printf("%s: %v %d", handType, fiveCards, score)

	assert.Equal(t, handType, FourOfAKind)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Contains(t, fiveCards, card7)
	assert.Equal(t, score, 7978669)
}

func TestScoreWithStraight(t *testing.T) {
	card1 := Card{Suit: DIAMONDS, Rank: TEN}
	card2 := Card{Suit: CLUBS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: QUEEN}
	card4 := Card{Suit: HEARTS, Rank: KING}
	card5 := Card{Suit: SPADES, Rank: ACE}
	card6 := Card{Suit: DIAMONDS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: ACE}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

	handType, fiveCards, score := Score(cards)
	fmt.Printf("%s: %v %d", handType, fiveCards, score)

	assert.Equal(t, Straight, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Equal(t, 4974010, score)
}

func TestScoreWithRoyalFlush(t *testing.T) {
	card1 := Card{Suit: HEARTS, Rank: TEN}
	card2 := Card{Suit: HEARTS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: QUEEN}
	card4 := Card{Suit: HEARTS, Rank: KING}
	card5 := Card{Suit: SPADES, Rank: ACE}
	card6 := Card{Suit: HEARTS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: ACE}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

	handType, fiveCards, score := Score(cards)
	fmt.Printf("%s: %v %d", handType, fiveCards, score)

	assert.Equal(t, RoyalFlush, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card6)
	assert.Equal(t, 9974010, score)
}

func TestScoreWithStraightFlush(t *testing.T) {
	card1 := Card{Suit: HEARTS, Rank: NINE}
	card2 := Card{Suit: HEARTS, Rank: TEN}
	card3 := Card{Suit: HEARTS, Rank: JACK}
	card4 := Card{Suit: HEARTS, Rank: QUEEN}
	card5 := Card{Suit: SPADES, Rank: KING}
	card6 := Card{Suit: HEARTS, Rank: KING}
	card7 := Card{Suit: CLUBS, Rank: KING}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

	handType, fiveCards, score := Score(cards)
	fmt.Printf("%s: %v %d", handType, fiveCards, score)

	assert.Equal(t, StraightFlush, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card6)
	assert.Equal(t, 8904105, score)
}
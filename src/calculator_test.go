package src

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScoreWithFourOfAKind(t *testing.T) {
	cards := Cards{
		{Suit: DIAMONDS, Rank: JACK},
		{Suit: CLUBS, Rank: JACK},
		{Suit: HEARTS, Rank: KING},
		{Suit: HEARTS, Rank: ACE},
		{Suit: SPADES, Rank: ACE},
		{Suit: DIAMONDS, Rank: ACE},
		{Suit: CLUBS, Rank: ACE},
	}

	handType, fiveCards, score := Score(cards)
	fmt.Printf("%s: %v %d", handType, fiveCards, score)

	assert.Equal(t, handType, FourOfAKind)
	assert.Contains(t, fiveCards, Card{Suit: HEARTS, Rank: ACE})
	assert.Contains(t, fiveCards, Card{Suit: SPADES, Rank: ACE})
	assert.Contains(t, fiveCards, Card{Suit: DIAMONDS, Rank: ACE})
	assert.Contains(t, fiveCards, Card{Suit: CLUBS, Rank: ACE})
	assert.Contains(t, fiveCards, Card{Suit: HEARTS, Rank: KING})
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

	assert.Equal(t, handType, Straight)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
}

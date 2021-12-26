package src

import (
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
	assert.Equal(t, handType, FourOfAKind)
	assert.Contains(t, fiveCards, Card{Suit: HEARTS, Rank: ACE})
	assert.Contains(t, fiveCards, Card{Suit: SPADES, Rank: ACE})
	assert.Contains(t, fiveCards, Card{Suit: DIAMONDS, Rank: ACE})
	assert.Contains(t, fiveCards, Card{Suit: CLUBS, Rank: ACE})
	assert.Contains(t, fiveCards, Card{Suit: HEARTS, Rank: KING})
	assert.Equal(t, score, 7978669)
}
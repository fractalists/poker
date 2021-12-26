package test

import (
	"github.com/stretchr/testify/assert"
	"holdem/entity"
	"holdem/util"
	"testing"
)

func TestScoreWithFourOfAKind(t *testing.T) {
	cards := entity.Cards{
		{Suit: entity.DIAMONDS, Rank: entity.JACK},
		{Suit: entity.CLUBS, Rank: entity.JACK},
		{Suit: entity.HEARTS, Rank: entity.KING},
		{Suit: entity.HEARTS, Rank: entity.ACE},
		{Suit: entity.SPADES, Rank: entity.ACE},
		{Suit: entity.DIAMONDS, Rank: entity.ACE},
		{Suit: entity.CLUBS, Rank: entity.ACE},
	}

	handType, fiveCards, score := util.Score(cards)
	assert.Equal(t, handType, util.FourOfAKind)
	assert.Contains(t, fiveCards, entity.Card{Suit: entity.HEARTS, Rank: entity.ACE})
	assert.Contains(t, fiveCards, entity.Card{Suit: entity.SPADES, Rank: entity.ACE})
	assert.Contains(t, fiveCards, entity.Card{Suit: entity.DIAMONDS, Rank: entity.ACE})
	assert.Contains(t, fiveCards, entity.Card{Suit: entity.CLUBS, Rank: entity.ACE})
	assert.Contains(t, fiveCards, entity.Card{Suit: entity.HEARTS, Rank: entity.KING})
	assert.Equal(t, score, 7978669)
}
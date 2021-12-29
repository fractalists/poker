package util

import (
	"fmt"
	"holdem/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScoreAll(t *testing.T) {
	t.Run("RoyalFlush", testScoreWithRoyalFlush)
	t.Run("StraightFlush", testScoreWithStraightFlush)
	t.Run("FourOfAKind", testScoreWithFourOfAKind)
	t.Run("FullHouse", testScoreWithFullHouse)
	t.Run("Flush", testScoreWithFlush)
	t.Run("Straight", testScoreWithStraight)
	t.Run("Straight2", testScoreWithStraight2)
	t.Run("ThreeOfAKind", testScoreWithThreeOfAKind)
	t.Run("TwoPair", testScoreWithTwoPair)
	t.Run("OnePair", testScoreWithOnePair)
	t.Run("HighCard", testScoreWithHighCard)
}

func testScoreWithRoyalFlush(t *testing.T) {
	card1 := model.Card{Suit: model.HEARTS, Rank: model.TEN}
	card2 := model.Card{Suit: model.HEARTS, Rank: model.JACK}
	card3 := model.Card{Suit: model.HEARTS, Rank: model.QUEEN}
	card4 := model.Card{Suit: model.HEARTS, Rank: model.KING}
	card5 := model.Card{Suit: model.SPADES, Rank: model.ACE}
	card6 := model.Card{Suit: model.HEARTS, Rank: model.ACE}
	card7 := model.Card{Suit: model.CLUBS, Rank: model.ACE}
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, RoyalFlush, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card6)
	assert.Equal(t, 9974010, score)
}

func testScoreWithStraightFlush(t *testing.T) {
	card1 := model.Card{Suit: model.HEARTS, Rank: model.NINE}
	card2 := model.Card{Suit: model.HEARTS, Rank: model.TEN}
	card3 := model.Card{Suit: model.HEARTS, Rank: model.JACK}
	card4 := model.Card{Suit: model.HEARTS, Rank: model.QUEEN}
	card5 := model.Card{Suit: model.SPADES, Rank: model.KING}
	card6 := model.Card{Suit: model.HEARTS, Rank: model.KING}
	card7 := model.Card{Suit: model.CLUBS, Rank: model.KING}
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, StraightFlush, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card6)
	assert.Equal(t, 8904105, score)
}

func testScoreWithFourOfAKind(t *testing.T) {
	card1 := model.Card{Suit: model.DIAMONDS, Rank: model.JACK}
	card2 := model.Card{Suit: model.CLUBS, Rank: model.JACK}
	card3 := model.Card{Suit: model.HEARTS, Rank: model.KING}
	card4 := model.Card{Suit: model.HEARTS, Rank: model.ACE}
	card5 := model.Card{Suit: model.SPADES, Rank: model.ACE}
	card6 := model.Card{Suit: model.DIAMONDS, Rank: model.ACE}
	card7 := model.Card{Suit: model.CLUBS, Rank: model.ACE}
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, FourOfAKind, handType)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Contains(t, fiveCards, card7)
	assert.Equal(t, 7978669, score)
}

func testScoreWithFullHouse(t *testing.T) {
	card1 := model.Card{Suit: model.HEARTS, Rank: model.NINE}
	card2 := model.Card{Suit: model.HEARTS, Rank: model.TEN}
	card3 := model.Card{Suit: model.CLUBS, Rank: model.QUEEN}
	card4 := model.Card{Suit: model.HEARTS, Rank: model.QUEEN}
	card5 := model.Card{Suit: model.SPADES, Rank: model.KING}
	card6 := model.Card{Suit: model.HEARTS, Rank: model.KING}
	card7 := model.Card{Suit: model.CLUBS, Rank: model.KING}
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, FullHouse, handType)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Contains(t, fiveCards, card7)
	assert.Equal(t, 6908748, score)
}

func testScoreWithFlush(t *testing.T) {
	card1 := model.Card{Suit: model.HEARTS, Rank: model.NINE}
	card2 := model.Card{Suit: model.HEARTS, Rank: model.TEN}
	card3 := model.Card{Suit: model.HEARTS, Rank: model.JACK}
	card4 := model.Card{Suit: model.HEARTS, Rank: model.QUEEN}
	card5 := model.Card{Suit: model.SPADES, Rank: model.FIVE}
	card6 := model.Card{Suit: model.HEARTS, Rank: model.TWO}
	card7 := model.Card{Suit: model.CLUBS, Rank: model.KING}
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, Flush, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card6)
	assert.Equal(t, 5834194, score)
}

func testScoreWithStraight(t *testing.T) {
	card1 := model.Card{Suit: model.DIAMONDS, Rank: model.TEN}
	card2 := model.Card{Suit: model.CLUBS, Rank: model.JACK}
	card3 := model.Card{Suit: model.HEARTS, Rank: model.QUEEN}
	card4 := model.Card{Suit: model.HEARTS, Rank: model.KING}
	card5 := model.Card{Suit: model.SPADES, Rank: model.ACE}
	card6 := model.Card{Suit: model.DIAMONDS, Rank: model.ACE}
	card7 := model.Card{Suit: model.CLUBS, Rank: model.ACE}
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, Straight, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Equal(t, 4974010, score)
}

func testScoreWithStraight2(t *testing.T) {
	card1 := model.Card{Suit: model.DIAMONDS, Rank: model.FIVE}
	card2 := model.Card{Suit: model.CLUBS, Rank: model.FOUR}
	card3 := model.Card{Suit: model.HEARTS, Rank: model.THREE}
	card4 := model.Card{Suit: model.HEARTS, Rank: model.TWO}
	card5 := model.Card{Suit: model.SPADES, Rank: model.ACE}
	card6 := model.Card{Suit: model.DIAMONDS, Rank: model.ACE}
	card7 := model.Card{Suit: model.CLUBS, Rank: model.ACE}
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, Straight, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Equal(t, 4939058, score)
}

func testScoreWithThreeOfAKind(t *testing.T) {
	card1 := model.Card{Suit: model.DIAMONDS, Rank: model.TEN}
	card2 := model.Card{Suit: model.CLUBS, Rank: model.JACK}
	card3 := model.Card{Suit: model.HEARTS, Rank: model.KING}
	card4 := model.Card{Suit: model.HEARTS, Rank: model.TWO}
	card5 := model.Card{Suit: model.SPADES, Rank: model.ACE}
	card6 := model.Card{Suit: model.DIAMONDS, Rank: model.ACE}
	card7 := model.Card{Suit: model.CLUBS, Rank: model.ACE}
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)
	assert.Equal(t, ThreeOfAKind, handType)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Contains(t, fiveCards, card7)
	assert.Equal(t, 3978651, score)
}

func testScoreWithTwoPair(t *testing.T) {
	card1 := model.Card{Suit: model.DIAMONDS, Rank: model.TEN}
	card2 := model.Card{Suit: model.CLUBS, Rank: model.JACK}
	card3 := model.Card{Suit: model.HEARTS, Rank: model.TWO}
	card4 := model.Card{Suit: model.HEARTS, Rank: model.FIVE}
	card5 := model.Card{Suit: model.SPADES, Rank: model.FIVE}
	card6 := model.Card{Suit: model.DIAMONDS, Rank: model.EIGHT}
	card7 := model.Card{Suit: model.CLUBS, Rank: model.EIGHT}
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, TwoPair, handType)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Contains(t, fiveCards, card7)
	assert.Equal(t, 2755797, score)
}

func testScoreWithOnePair(t *testing.T) {
	card1 := model.Card{Suit: model.DIAMONDS, Rank: model.TEN}
	card2 := model.Card{Suit: model.CLUBS, Rank: model.JACK}
	card3 := model.Card{Suit: model.HEARTS, Rank: model.TWO}
	card4 := model.Card{Suit: model.HEARTS, Rank: model.FIVE}
	card5 := model.Card{Suit: model.SPADES, Rank: model.FIVE}
	card6 := model.Card{Suit: model.DIAMONDS, Rank: model.ACE}
	card7 := model.Card{Suit: model.CLUBS, Rank: model.EIGHT}
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, OnePair, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Equal(t, 1965205, score)
}

func testScoreWithHighCard(t *testing.T) {
	card1 := model.Card{Suit: model.DIAMONDS, Rank: model.TEN}
	card2 := model.Card{Suit: model.CLUBS, Rank: model.JACK}
	card3 := model.Card{Suit: model.HEARTS, Rank: model.TWO}
	card4 := model.Card{Suit: model.HEARTS, Rank: model.THREE}
	card5 := model.Card{Suit: model.SPADES, Rank: model.FIVE}
	card6 := model.Card{Suit: model.DIAMONDS, Rank: model.ACE}
	card7 := model.Card{Suit: model.CLUBS, Rank: model.EIGHT}
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, HighCard, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Contains(t, fiveCards, card7)
	assert.Equal(t, 965253, score)
}

package src

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAll(t *testing.T) {
	t.Run("RoyalFlush", TestScoreWithRoyalFlush)
	t.Run("StraightFlush", TestScoreWithStraightFlush)
	t.Run("FourOfAKind", TestScoreWithFourOfAKind)
	t.Run("FullHouse", TestScoreWithFullHouse)
	t.Run("Flush", TestScoreWithFlush)
	t.Run("Straight", TestScoreWithStraight)
	t.Run("Straight2", TestScoreWithStraight2)
	t.Run("ThreeOfAKind", TestScoreWithThreeOfAKind)
	t.Run("TwoPair", TestScoreWithTwoPair)
	t.Run("OnePair", TestScoreWithOnePair)
	t.Run("HighCard", TestScoreWithHighCard)
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
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

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
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, StraightFlush, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card6)
	assert.Equal(t, 8904105, score)
}

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
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, FourOfAKind, handType)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Contains(t, fiveCards, card7)
	assert.Equal(t, 7978669, score)
}

func TestScoreWithFullHouse(t *testing.T) {
	card1 := Card{Suit: HEARTS, Rank: NINE}
	card2 := Card{Suit: HEARTS, Rank: TEN}
	card3 := Card{Suit: CLUBS, Rank: QUEEN}
	card4 := Card{Suit: HEARTS, Rank: QUEEN}
	card5 := Card{Suit: SPADES, Rank: KING}
	card6 := Card{Suit: HEARTS, Rank: KING}
	card7 := Card{Suit: CLUBS, Rank: KING}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

	handType, fiveCards, score := Score(cards)
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, FullHouse, handType)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Contains(t, fiveCards, card7)
	assert.Equal(t, 6908748, score)
}

func TestScoreWithFlush(t *testing.T) {
	card1 := Card{Suit: HEARTS, Rank: NINE}
	card2 := Card{Suit: HEARTS, Rank: TEN}
	card3 := Card{Suit: HEARTS, Rank: JACK}
	card4 := Card{Suit: HEARTS, Rank: QUEEN}
	card5 := Card{Suit: SPADES, Rank: FIVE}
	card6 := Card{Suit: HEARTS, Rank: TWO}
	card7 := Card{Suit: CLUBS, Rank: KING}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

	handType, fiveCards, score := Score(cards)
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, Flush, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card6)
	assert.Equal(t, 5834194, score)
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
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, Straight, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Equal(t, 4974010, score)
}

func TestScoreWithStraight2(t *testing.T) {
	card1 := Card{Suit: DIAMONDS, Rank: FIVE}
	card2 := Card{Suit: CLUBS, Rank: FOUR}
	card3 := Card{Suit: HEARTS, Rank: THREE}
	card4 := Card{Suit: HEARTS, Rank: TWO}
	card5 := Card{Suit: SPADES, Rank: ACE}
	card6 := Card{Suit: DIAMONDS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: ACE}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

	handType, fiveCards, score := Score(cards)
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, Straight, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Equal(t, 4939058, score)
}

func TestScoreWithThreeOfAKind(t *testing.T) {
	card1 := Card{Suit: DIAMONDS, Rank: TEN}
	card2 := Card{Suit: CLUBS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: KING}
	card4 := Card{Suit: HEARTS, Rank: TWO}
	card5 := Card{Suit: SPADES, Rank: ACE}
	card6 := Card{Suit: DIAMONDS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: ACE}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

	handType, fiveCards, score := Score(cards)
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)
	assert.Equal(t, ThreeOfAKind, handType)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Contains(t, fiveCards, card7)
	assert.Equal(t, 3978651, score)
}

func TestScoreWithTwoPair(t *testing.T) {
	card1 := Card{Suit: DIAMONDS, Rank: TEN}
	card2 := Card{Suit: CLUBS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: TWO}
	card4 := Card{Suit: HEARTS, Rank: FIVE}
	card5 := Card{Suit: SPADES, Rank: FIVE}
	card6 := Card{Suit: DIAMONDS, Rank: EIGHT}
	card7 := Card{Suit: CLUBS, Rank: EIGHT}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

	handType, fiveCards, score := Score(cards)
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, TwoPair, handType)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Contains(t, fiveCards, card7)
	assert.Equal(t, 2755797, score)
}

func   TestScoreWithOnePair(t *testing.T) {
	card1 := Card{Suit: DIAMONDS, Rank: TEN}
	card2 := Card{Suit: CLUBS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: TWO}
	card4 := Card{Suit: HEARTS, Rank: FIVE}
	card5 := Card{Suit: SPADES, Rank: FIVE}
	card6 := Card{Suit: DIAMONDS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: EIGHT}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

	handType, fiveCards, score := Score(cards)
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, OnePair, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Equal(t, 1965205, score)
}

func TestScoreWithHighCard(t *testing.T) {
	card1 := Card{Suit: DIAMONDS, Rank: TEN}
	card2 := Card{Suit: CLUBS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: TWO}
	card4 := Card{Suit: HEARTS, Rank: THREE}
	card5 := Card{Suit: SPADES, Rank: FIVE}
	card6 := Card{Suit: DIAMONDS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: EIGHT}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

	handType, fiveCards, score := Score(cards)
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, HighCard, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Contains(t, fiveCards, card7)
	assert.Equal(t, 965253, score)
}

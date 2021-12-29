package src

import (
	"fmt"
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
	card1 := Card{Suit: HEARTS, Rank: TEN}
	card2 := Card{Suit: HEARTS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: QUEEN}
	card4 := Card{Suit: HEARTS, Rank: KING}
	card5 := Card{Suit: SPADES, Rank: ACE}
	card6 := Card{Suit: HEARTS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: ACE}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

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
	card1 := Card{Suit: HEARTS, Rank: NINE}
	card2 := Card{Suit: HEARTS, Rank: TEN}
	card3 := Card{Suit: HEARTS, Rank: JACK}
	card4 := Card{Suit: HEARTS, Rank: QUEEN}
	card5 := Card{Suit: SPADES, Rank: KING}
	card6 := Card{Suit: HEARTS, Rank: KING}
	card7 := Card{Suit: CLUBS, Rank: KING}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

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
	card1 := Card{Suit: DIAMONDS, Rank: JACK}
	card2 := Card{Suit: CLUBS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: KING}
	card4 := Card{Suit: HEARTS, Rank: ACE}
	card5 := Card{Suit: SPADES, Rank: ACE}
	card6 := Card{Suit: DIAMONDS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: ACE}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

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
	card1 := Card{Suit: HEARTS, Rank: NINE}
	card2 := Card{Suit: HEARTS, Rank: TEN}
	card3 := Card{Suit: CLUBS, Rank: QUEEN}
	card4 := Card{Suit: HEARTS, Rank: QUEEN}
	card5 := Card{Suit: SPADES, Rank: KING}
	card6 := Card{Suit: HEARTS, Rank: KING}
	card7 := Card{Suit: CLUBS, Rank: KING}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

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
	card1 := Card{Suit: HEARTS, Rank: NINE}
	card2 := Card{Suit: HEARTS, Rank: TEN}
	card3 := Card{Suit: HEARTS, Rank: JACK}
	card4 := Card{Suit: HEARTS, Rank: QUEEN}
	card5 := Card{Suit: SPADES, Rank: FIVE}
	card6 := Card{Suit: HEARTS, Rank: TWO}
	card7 := Card{Suit: CLUBS, Rank: KING}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

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
	card1 := Card{Suit: DIAMONDS, Rank: TEN}
	card2 := Card{Suit: CLUBS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: QUEEN}
	card4 := Card{Suit: HEARTS, Rank: KING}
	card5 := Card{Suit: SPADES, Rank: ACE}
	card6 := Card{Suit: DIAMONDS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: ACE}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

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
	card1 := Card{Suit: DIAMONDS, Rank: FIVE}
	card2 := Card{Suit: CLUBS, Rank: FOUR}
	card3 := Card{Suit: HEARTS, Rank: THREE}
	card4 := Card{Suit: HEARTS, Rank: TWO}
	card5 := Card{Suit: SPADES, Rank: ACE}
	card6 := Card{Suit: DIAMONDS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: ACE}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

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
	card1 := Card{Suit: DIAMONDS, Rank: TEN}
	card2 := Card{Suit: CLUBS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: KING}
	card4 := Card{Suit: HEARTS, Rank: TWO}
	card5 := Card{Suit: SPADES, Rank: ACE}
	card6 := Card{Suit: DIAMONDS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: ACE}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

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
	card1 := Card{Suit: DIAMONDS, Rank: TEN}
	card2 := Card{Suit: CLUBS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: TWO}
	card4 := Card{Suit: HEARTS, Rank: FIVE}
	card5 := Card{Suit: SPADES, Rank: FIVE}
	card6 := Card{Suit: DIAMONDS, Rank: EIGHT}
	card7 := Card{Suit: CLUBS, Rank: EIGHT}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

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
	card1 := Card{Suit: DIAMONDS, Rank: TEN}
	card2 := Card{Suit: CLUBS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: TWO}
	card4 := Card{Suit: HEARTS, Rank: FIVE}
	card5 := Card{Suit: SPADES, Rank: FIVE}
	card6 := Card{Suit: DIAMONDS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: EIGHT}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

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
	card1 := Card{Suit: DIAMONDS, Rank: TEN}
	card2 := Card{Suit: CLUBS, Rank: JACK}
	card3 := Card{Suit: HEARTS, Rank: TWO}
	card4 := Card{Suit: HEARTS, Rank: THREE}
	card5 := Card{Suit: SPADES, Rank: FIVE}
	card6 := Card{Suit: DIAMONDS, Rank: ACE}
	card7 := Card{Suit: CLUBS, Rank: EIGHT}
	cards := Cards{card1, card2, card3, card4, card5, card6, card7}

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

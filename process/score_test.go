package process

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
	t.Run("TwoPair2", testScoreWithTwoPair2)
	t.Run("TwoPair3", testScoreWithTwoPair3)
	t.Run("OnePair", testScoreWithOnePair)
	t.Run("HighCard", testScoreWithHighCard)
}

func testScoreWithRoyalFlush(t *testing.T) {
	card1 := model.NewCustomCard(model.HEARTS, model.TEN, true)
	card2 := model.NewCustomCard(model.HEARTS, model.JACK, true)
	card3 := model.NewCustomCard(model.HEARTS, model.QUEEN, true)
	card4 := model.NewCustomCard(model.HEARTS, model.KING, true)
	card5 := model.NewCustomCard(model.SPADES, model.ACE, true)
	card6 := model.NewCustomCard(model.HEARTS, model.ACE, true)
	card7 := model.NewCustomCard(model.CLUBS, model.ACE, true)
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
	card1 := model.NewCustomCard(model.HEARTS, model.NINE, true)
	card2 := model.NewCustomCard(model.HEARTS, model.TEN, true)
	card3 := model.NewCustomCard(model.HEARTS, model.JACK, true)
	card4 := model.NewCustomCard(model.HEARTS, model.QUEEN, true)
	card5 := model.NewCustomCard(model.SPADES, model.KING, true)
	card6 := model.NewCustomCard(model.HEARTS, model.KING, true)
	card7 := model.NewCustomCard(model.CLUBS, model.KING, true)
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
	card1 := model.NewCustomCard(model.DIAMONDS, model.JACK, true)
	card2 := model.NewCustomCard(model.CLUBS, model.JACK, true)
	card3 := model.NewCustomCard(model.HEARTS, model.KING, true)
	card4 := model.NewCustomCard(model.HEARTS, model.ACE, true)
	card5 := model.NewCustomCard(model.SPADES, model.ACE, true)
	card6 := model.NewCustomCard(model.DIAMONDS, model.ACE, true)
	card7 := model.NewCustomCard(model.CLUBS, model.ACE, true)
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
	card1 := model.NewCustomCard(model.HEARTS, model.NINE, true)
	card2 := model.NewCustomCard(model.HEARTS, model.TEN, true)
	card3 := model.NewCustomCard(model.CLUBS, model.QUEEN, true)
	card4 := model.NewCustomCard(model.HEARTS, model.QUEEN, true)
	card5 := model.NewCustomCard(model.SPADES, model.KING, true)
	card6 := model.NewCustomCard(model.HEARTS, model.KING, true)
	card7 := model.NewCustomCard(model.CLUBS, model.KING, true)
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
	card1 := model.NewCustomCard(model.HEARTS, model.NINE, true)
	card2 := model.NewCustomCard(model.HEARTS, model.TEN, true)
	card3 := model.NewCustomCard(model.HEARTS, model.JACK, true)
	card4 := model.NewCustomCard(model.HEARTS, model.QUEEN, true)
	card5 := model.NewCustomCard(model.SPADES, model.FIVE, true)
	card6 := model.NewCustomCard(model.HEARTS, model.TWO, true)
	card7 := model.NewCustomCard(model.CLUBS, model.KING, true)
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
	card1 := model.NewCustomCard(model.DIAMONDS, model.TEN, true)
	card2 := model.NewCustomCard(model.CLUBS, model.JACK, true)
	card3 := model.NewCustomCard(model.HEARTS, model.QUEEN, true)
	card4 := model.NewCustomCard(model.HEARTS, model.KING, true)
	card5 := model.NewCustomCard(model.SPADES, model.ACE, true)
	card6 := model.NewCustomCard(model.DIAMONDS, model.ACE, true)
	card7 := model.NewCustomCard(model.CLUBS, model.ACE, true)
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
	card1 := model.NewCustomCard(model.DIAMONDS, model.FIVE, true)
	card2 := model.NewCustomCard(model.CLUBS, model.FOUR, true)
	card3 := model.NewCustomCard(model.HEARTS, model.THREE, true)
	card4 := model.NewCustomCard(model.HEARTS, model.TWO, true)
	card5 := model.NewCustomCard(model.SPADES, model.ACE, true)
	card6 := model.NewCustomCard(model.DIAMONDS, model.ACE, true)
	card7 := model.NewCustomCard(model.CLUBS, model.ACE, true)
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, Straight, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Equal(t, 4344878, score)
}

func testScoreWithThreeOfAKind(t *testing.T) {
	card1 := model.NewCustomCard(model.DIAMONDS, model.TEN, true)
	card2 := model.NewCustomCard(model.CLUBS, model.JACK, true)
	card3 := model.NewCustomCard(model.HEARTS, model.KING, true)
	card4 := model.NewCustomCard(model.HEARTS, model.TWO, true)
	card5 := model.NewCustomCard(model.SPADES, model.ACE, true)
	card6 := model.NewCustomCard(model.DIAMONDS, model.ACE, true)
	card7 := model.NewCustomCard(model.CLUBS, model.ACE, true)
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
	card1 := model.NewCustomCard(model.DIAMONDS, model.TEN, true)
	card2 := model.NewCustomCard(model.CLUBS, model.JACK, true)
	card3 := model.NewCustomCard(model.HEARTS, model.TWO, true)
	card4 := model.NewCustomCard(model.HEARTS, model.FIVE, true)
	card5 := model.NewCustomCard(model.SPADES, model.FIVE, true)
	card6 := model.NewCustomCard(model.DIAMONDS, model.EIGHT, true)
	card7 := model.NewCustomCard(model.CLUBS, model.EIGHT, true)
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
	assert.Equal(t, 2558427, score)
}

func testScoreWithTwoPair2(t *testing.T) {
	card1 := model.NewCustomCard(model.HEARTS, model.SIX, true)
	card2 := model.NewCustomCard(model.DIAMONDS, model.QUEEN, true)
	card3 := model.NewCustomCard(model.DIAMONDS, model.TWO, true)
	card4 := model.NewCustomCard(model.SPADES, model.TEN, true)
	card5 := model.NewCustomCard(model.CLUBS, model.TEN, true)
	card6 := model.NewCustomCard(model.HEARTS, model.TWO, true)
	card7 := model.NewCustomCard(model.HEARTS, model.THREE, true)
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, TwoPair, handType)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card3)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card6)
	assert.Equal(t, 2696876, score)
}

func testScoreWithTwoPair3(t *testing.T) {
	card1 := model.NewCustomCard(model.DIAMONDS, model.EIGHT, true)
	card2 := model.NewCustomCard(model.CLUBS, model.EIGHT, true)
	card3 := model.NewCustomCard(model.DIAMONDS, model.TWO, true)
	card4 := model.NewCustomCard(model.SPADES, model.TEN, true)
	card5 := model.NewCustomCard(model.CLUBS, model.TEN, true)
	card6 := model.NewCustomCard(model.HEARTS, model.TWO, true)
	card7 := model.NewCustomCard(model.HEARTS, model.THREE, true)
	cards := model.Cards{card1, card2, card3, card4, card5, card6, card7}

	scoreResult := Score(cards)
	handType, fiveCards, score := scoreResult.HandType, scoreResult.FinalCards, scoreResult.Score
	fmt.Printf("%s: %v %d\n", handType, fiveCards, score)

	assert.Equal(t, TwoPair, handType)
	assert.Contains(t, fiveCards, card1)
	assert.Contains(t, fiveCards, card2)
	assert.Contains(t, fiveCards, card4)
	assert.Contains(t, fiveCards, card5)
	assert.Contains(t, fiveCards, card7)
	assert.Equal(t, 2698499, score)
}

func testScoreWithOnePair(t *testing.T) {
	card1 := model.NewCustomCard(model.DIAMONDS, model.TEN, true)
	card2 := model.NewCustomCard(model.CLUBS, model.JACK, true)
	card3 := model.NewCustomCard(model.HEARTS, model.TWO, true)
	card4 := model.NewCustomCard(model.HEARTS, model.FIVE, true)
	card5 := model.NewCustomCard(model.SPADES, model.FIVE, true)
	card6 := model.NewCustomCard(model.DIAMONDS, model.ACE, true)
	card7 := model.NewCustomCard(model.CLUBS, model.EIGHT, true)
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
	assert.Equal(t, 1351930, score)
}

func testScoreWithHighCard(t *testing.T) {
	card1 := model.NewCustomCard(model.DIAMONDS, model.TEN, true)
	card2 := model.NewCustomCard(model.CLUBS, model.JACK, true)
	card3 := model.NewCustomCard(model.HEARTS, model.TWO, true)
	card4 := model.NewCustomCard(model.HEARTS, model.THREE, true)
	card5 := model.NewCustomCard(model.SPADES, model.FIVE, true)
	card6 := model.NewCustomCard(model.DIAMONDS, model.ACE, true)
	card7 := model.NewCustomCard(model.CLUBS, model.EIGHT, true)
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

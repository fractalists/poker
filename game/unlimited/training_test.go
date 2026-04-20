package unlimited

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectWinnerCountsAggregatesWorkerResults(t *testing.T) {
	winners := []int{0, 1, 0, 2}
	var mu sync.Mutex
	next := 0

	counts := collectWinnerCounts(len(winners), func() int {
		mu.Lock()
		defer mu.Unlock()
		winner := winners[next]
		next++
		return winner
	})

	assert.Equal(t, 2, counts[0])
	assert.Equal(t, 1, counts[1])
	assert.Equal(t, 1, counts[2])
}

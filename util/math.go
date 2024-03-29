package util

import (
	"math/rand"
	"time"
)

func Min(a, b int) int {
	if a <= b {
		return a
	}

	return b
}

func MaxFloat32(a, b float32) float32 {
	if a >= b {
		return a
	}

	return b
}

func Max(a, b int) int {
	if a >= b {
		return a
	}

	return b
}

func NewRng() *rand.Rand {
	source := rand.NewSource(time.Now().UnixNano())
	return rand.New(source)
}

func Shuffle(n int, rng *rand.Rand, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to Shuffle")
	}

	// Fisher-Yates shuffle: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	// Shuffle really ought not be called with n that doesn't fit in 32 bits.
	// Not only will it take a very long time, but with 2³¹! possible permutations,
	// there's no way that any PRNG can have a big enough internal state to
	// generate even a minuscule percentage of the possible permutations.
	// Nevertheless, the right API signature accepts an int n, so handle it as best we can.
	i := n - 1
	for ; i > 1<<31-1-1; i-- {
		j := int(rng.Int63n(int64(i + 1)))
		swap(i, j)
	}
	for ; i > 0; i-- {
		j := int(rng.Int31n(int32(i + 1)))
		swap(i, j)
	}
}

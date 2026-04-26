package table

import "strings"

const (
	AIStyleMixed        = "mixed"
	AIStyleRandom       = "random"
	AIStyleSmart        = "smart"
	AIStyleConservative = "conservative"
	AIStyleAggressive   = "aggressive"
	AIStyleGTO          = "gto"
)

func NormalizeAIStyle(style string) string {
	switch strings.ToLower(strings.TrimSpace(style)) {
	case "", AIStyleMixed, AIStyleRandom, "random-seats", "random_seats", "seat-random", "seat_random":
		return AIStyleRandom
	case AIStyleSmart, "odds", "odds-warrior", "odds_warrior":
		return AIStyleSmart
	case AIStyleConservative, "tight", "tight-conservative", "tight_conservative":
		return AIStyleConservative
	case AIStyleAggressive, "loose", "loose-aggressive", "loose_aggressive":
		return AIStyleAggressive
	case AIStyleGTO, "solver", "solver-ish", "solver_ish", "balanced":
		return AIStyleGTO
	default:
		return AIStyleRandom
	}
}

func NormalizeSeatAIStyle(style string) string {
	normalized := NormalizeAIStyle(style)
	switch normalized {
	case AIStyleSmart, AIStyleConservative, AIStyleAggressive, AIStyleGTO:
		return normalized
	default:
		return ""
	}
}

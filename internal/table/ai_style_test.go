package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeAIStyleTreatsMixedAsRandomAlias(t *testing.T) {
	assert.Equal(t, AIStyleRandom, NormalizeAIStyle(""))
	assert.Equal(t, AIStyleRandom, NormalizeAIStyle("mixed"))
	assert.Equal(t, AIStyleRandom, NormalizeAIStyle("random-seats"))
	assert.Equal(t, AIStyleRandom, NormalizeAIStyle("unknown-style"))
}

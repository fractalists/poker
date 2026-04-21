package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOptionsUsesPortableDefaults(t *testing.T) {
	opts, err := parseOptions([]string{})

	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1:8080", opts.addr)
	assert.Equal(t, "info", opts.logLevel)
	assert.Equal(t, "web/dist", opts.webDist)
}

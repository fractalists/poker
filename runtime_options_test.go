package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRuntimeOptionsUsesPortableDefaults(t *testing.T) {
	now := time.Unix(1710000000, 0)

	opts, err := parseRuntimeOptions([]string{}, now)

	require.NoError(t, err)
	assert.Equal(t, "unlimited", opts.mode)
	assert.Equal(t, "generated/log/poker_log_1710000000.log", opts.logPath)
	assert.Equal(t, "", opts.profilePath)
}

func TestParseRuntimeOptionsSupportsTrainOverrides(t *testing.T) {
	now := time.Unix(1710000000, 0)

	opts, err := parseRuntimeOptions([]string{
		"-mode=train",
		"-profile-path=tmp/train.pprof",
		"-log-level=warn",
	}, now)

	require.NoError(t, err)
	assert.Equal(t, "train", opts.mode)
	assert.Equal(t, "tmp/train.pprof", opts.profilePath)
	assert.Equal(t, "warn", opts.logLevel)
}

package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"time"
)

type runtimeOptions struct {
	mode        string
	logPath     string
	profilePath string
	logLevel    string
}

func parseRuntimeOptions(args []string, now time.Time) (runtimeOptions, error) {
	fs := flag.NewFlagSet("poker", flag.ContinueOnError)
	mode := fs.String("mode", "unlimited", "unlimited|train|gui|colosseum")
	logPath := fs.String("log-path", "", "log output path")
	profilePath := fs.String("profile-path", "", "cpu profile output path")
	logLevel := fs.String("log-level", "", "debug|info|warn|error")

	if err := fs.Parse(args); err != nil {
		return runtimeOptions{}, err
	}

	opts := runtimeOptions{
		mode:        *mode,
		logPath:     *logPath,
		profilePath: *profilePath,
		logLevel:    *logLevel,
	}

	switch opts.mode {
	case "unlimited":
		if opts.logPath == "" {
			opts.logPath = generatedLogPath(now)
		}
		if opts.logLevel == "" {
			opts.logLevel = "debug"
		}
	case "train":
		if opts.profilePath == "" {
			opts.profilePath = generatedProfilePath(now)
		}
		if opts.logLevel == "" {
			opts.logLevel = "warn"
		}
	case "colosseum":
		if opts.profilePath == "" {
			opts.profilePath = generatedProfilePath(now)
		}
		if opts.logLevel == "" {
			opts.logLevel = "warn"
		}
	case "gui":
		if opts.logLevel == "" {
			opts.logLevel = "debug"
		}
	default:
		return runtimeOptions{}, fmt.Errorf("unknown mode: %s", opts.mode)
	}

	return opts, nil
}

func generatedLogPath(now time.Time) string {
	return filepath.ToSlash(filepath.Join("generated", "log", fmt.Sprintf("poker_log_%d.log", now.Unix())))
}

func generatedProfilePath(now time.Time) string {
	return filepath.ToSlash(filepath.Join("generated", "pprof", fmt.Sprintf("poker_pprof_%d.pprof", now.Unix())))
}

func parseLogLevelValue(value string) (logrus.Level, error) {
	return logrus.ParseLevel(strings.ToLower(value))
}

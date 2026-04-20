package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fmt"
	"poker/config"
	"poker/game/colosseum"
	"poker/game/unlimited"
	"poker/process"
	"os"
	"time"
)

func main() {
	opts, err := parseRuntimeOptions(os.Args[1:], time.Now())
	if err != nil {
		panic(err)
	}

	switch opts.mode {
	case "unlimited":
		playUnlimited(opts)
	case "train":
		trainWithProfiler(opts)
	case "gui":
		tryFyne()
	case "colosseum":
		playColosseum(opts)
	default:
		panic(fmt.Sprintf("unknown mode: %s", opts.mode))
	}
}

func tryFyne() {
	a := app.New()
	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	w.ShowAndRun()
}

func playUnlimited(opts runtimeOptions) {
	logLevel, err := parseLogLevelValue(opts.logLevel)
	if err != nil {
		panic(err)
	}

	process.Start(
		false,
		config.ZhCn,
		logLevel,
		opts.logPath,
		opts.profilePath,
		unlimited.PlayPoker)
}

func trainWithProfiler(opts runtimeOptions) {
	logLevel, err := parseLogLevelValue(opts.logLevel)
	if err != nil {
		panic(err)
	}

	process.Start(
		true,
		config.ZhCn,
		logLevel,
		opts.logPath,
		opts.profilePath,
		unlimited.Train)
}

func playColosseum(opts runtimeOptions) {
	logLevel, err := parseLogLevelValue(opts.logLevel)
	if err != nil {
		panic(err)
	}

	process.Start(
		false,
		config.ZhCn,
		logLevel,
		opts.logPath,
		opts.profilePath,
		colosseum.PlayPoker)
}

package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"poker/config"
	"poker/game/colosseum"
	"poker/game/unlimited"
	"poker/process"
	"time"
)

func main() {
	switch 1 {
	case 1:
		playUnlimited()
	case 2:
		trainWithProfiler()
	case 3:
		tryFyne()
	case 4:
		playColosseum()
	default:
		playUnlimited()
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

func playUnlimited() {
	process.Start(
		false,
		config.ZhCn,
		logrus.DebugLevel,
		fmt.Sprintf("D:/Git/go/src/poker/generated/log/poker_log_%d.log", time.Now().Unix()), //filepath.Join("generated", "log", fmt.Sprintf("poker_log_%d.log", time.Now().Unix())),
		"",
		unlimited.PlayPoker)
}

func trainWithProfiler() {
	process.Start(
		true,
		config.ZhCn,
		logrus.WarnLevel,
		"",
		fmt.Sprintf("D:/Git/go/src/poker/generated/pprof/poker_pprof_%d.pprof", time.Now().Unix()), //filepath.Join("generated", "pprof", fmt.Sprintf("poker_pprof_%d.pprof", time.Now().Unix())),
		unlimited.Train)
}

func playColosseum() {
	process.Start(
		false,
		config.ZhCn,
		logrus.WarnLevel,
		"",
		fmt.Sprintf("D:/Git/go/src/poker/generated/pprof/poker_pprof_%d.pprof", time.Now().Unix()), //filepath.Join("generated", "pprof", fmt.Sprintf("poker_pprof_%d.pprof", time.Now().Unix())),
		colosseum.PlayPoker)
}

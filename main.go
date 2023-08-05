package main

import (
	"poker/config"
	"poker/game/colosseum"
	"poker/game/unlimited"
	"poker/process"

	"github.com/sirupsen/logrus"
)

func main() {
	switch 1 {
	case 1:
		playUnlimited()
	case 2:
		trainWithProfiler()
	case 3:
		playColosseum()
	default:
		playUnlimited()
	}
}

func playUnlimited() {
	process.Start(
		false,
		false,
		false,
		config.ZhCn,
		logrus.DebugLevel,
		unlimited.PlayPoker)
}

func trainWithProfiler() {
	process.Start(
		true,
		false,
		true,
		config.ZhCn,
		logrus.DebugLevel,
		unlimited.Train)
}

func playColosseum() {
	process.Start(
		false,
		false,
		false,
		config.ZhCn,
		logrus.DebugLevel,
		colosseum.PlayPoker)
}

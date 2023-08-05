package main

import (
	"fmt"
	"poker/config"
	"poker/game/colosseum"
	"poker/game/unlimited"
	"poker/process"
	"time"

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
		true,
		false,
		config.ZhCn,
		logrus.DebugLevel,
		fmt.Sprintf("./generated/log/poker_log_%d.log", time.Now().Unix()),
		unlimited.PlayPoker)
}

func trainWithProfiler() {
	process.Start(
		true,
		true,
		config.ZhCn,
		logrus.WarnLevel,
		"",
		unlimited.Train)
}

func playColosseum() {
	process.Start(
		false,
		false,
		config.ZhCn,
		logrus.WarnLevel,
		"",
		colosseum.PlayPoker)
}

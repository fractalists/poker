package config

import (
	"github.com/panjf2000/ants/v2"
)

var DebugMode = false
var TrainMode = false
var GoroutineLimit = 1
var Pool *ants.Pool

var Language = EnUs

const EnUs string = "en-US"
const ZhCn string = "zh-CN"

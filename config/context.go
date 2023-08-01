package config

import "github.com/panjf2000/ants/v2"

var DebugMode = false
var TrainMode = false
var GoroutineLimit = 1
var Pool *ants.Pool

var Language = EN_US

const EN_US string = "en-US"
const ZH_CN string = "zh-CN"

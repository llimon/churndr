package common

import (
	"go.uber.org/zap"
)

var Logger, _ = zap.NewProduction()
var Sugar = Logger.Sugar()

var PodCache = make(map[string]PodDB)

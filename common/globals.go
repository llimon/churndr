package common

import (
	"go.uber.org/zap"
)

var Logger, _ = zap.NewProduction()
var Sugar = Logger.Sugar()

var PodCache = make(map[string]PodDB)

var PodLogs = make(map[string]PodLogsDB)
var PodLogHistory = make(map[string][]PodLogHistoryDB)

// Config - Holds application configuation Passed down by cobra
var Config Configuration

// DevelopmentMode - If true server with listen to localhost
var DevelopmentMode bool

// ListenPort - Port where rest service listens
var ListenPort int

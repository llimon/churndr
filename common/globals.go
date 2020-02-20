package common

import (
	churndrv1alpha1 "github.com/llimon/churndr/pkg/apis/churndrcontroller/v1alpha1"
	"go.uber.org/zap"
)

var Logger, _ = zap.NewProduction()
var Sugar = Logger.Sugar()

var PodCache = make(map[string]PodDB)

var PodLogs = make(map[string]PodLogsDB)
var PodLogHistory = make(map[string][]PodLogHistoryDB)

// podChurnList - contains list of resources bo monitor by namespace.
var PodChurnList = make(map[string]*churndrv1alpha1.Podchurn)

// Config - Holds application configuation Passed down by cobra
var Config Configuration

// DevelopmentMode - If true server with listen to localhost
var DevelopmentMode bool

// ListenPort - Port where rest service listens
var ListenPort int

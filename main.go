package main

import (
	"time"

	"github.com/llimon/page-on-pod-restarts/common"
	"github.com/llimon/page-on-pod-restarts/server"
)

func main() {

	defer common.Logger.Sync() // flushes buffer, if any

	common.Sugar.Infof("Creating notified ticker....")

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for t := range ticker.C {
			//Call the periodic function here.
			KubeNotifierFunc(t)
		}
	}()

	quit := make(chan bool, 1)

	go server.RESTServer()

	KubeGetPods()

	// main will continue to wait untill there is an entry in quit
	<-quit
}

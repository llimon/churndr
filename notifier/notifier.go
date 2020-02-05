package notifier

import (
	"fmt"
	"time"

	"github.com/llimon/churndr/common"
	"github.com/llimon/churndr/util"
)

func KubeNotifierFunc(tick time.Time) {

	var now = time.Now()
	LookBackTime := 15
	var then = now.Add(time.Duration(-LookBackTime) * time.Minute)
	anythingToReport := false

	common.Sugar.Infow("Notifier", "started at", util.GetDateString(tick))
	out := fmt.Sprintf("List of Pods with Issues in last %v minutes\n", LookBackTime)

	for _, currNamespaceToReport := range common.Config.Namespaces {
		hasPodsToReportInNamespace := func() bool {
			for _, currPod := range common.PodCache {
				if currPod.LastTimeReported.Unix() < then.Unix() && currPod.Reported == false && currPod.Namespace == currNamespaceToReport {
					return true
				}
			}
			return false
		}
		if hasPodsToReportInNamespace() {
			out += ">>Namespace [" + currNamespaceToReport + "]<<\n"
			for _, currPod := range common.PodCache {

				// Reported==true is only set onse we send a notification. if the object is reconstructed we only persist LastTimeReported for noise reduction
				// Intention is to only send notifications for pod that have new Statuses.
				if currPod.LastTimeReported.Unix() < then.Unix() && currPod.Reported == false && currPod.Namespace == currNamespaceToReport {
					out += fmt.Sprintf("Pod: %v\n", currPod.Name)
					for _, currContainer := range currPod.Container {

						out += fmt.Sprintf("\tcontainer: %v\n", currContainer.Name)
						out += fmt.Sprintf("\t\trestarts [%v]\n", currContainer.RestartCount)

						if currContainer.Terminated != nil {
							out += fmt.Sprintf("\t\tterminated Finished At [%v] ExitCode [%v], Reason [%v]\n", currContainer.Terminated.FinishedAt, currContainer.Terminated.ExitCode, currContainer.Terminated.Reason)

						}
						//if currContainer.Running != &b {
						if currContainer.Running != nil {
							out += fmt.Sprintf("\t\tRunning Started At [%v]\n", util.GetDateString(currContainer.Running.StartedAt.Time))
						}
					}
					out += fmt.Sprintf("\n")
					// Mark this pod as reported and don't alert unless it repeats restarts after reporting time window.
					currPod.LastTimeReported = now
					currPod.Reported = true
					common.PodCache[currPod.Name] = currPod
					// Yep we have stuff to report and anoy people
					anythingToReport = true
				}
			}

		}

	}

	if anythingToReport {
		// Bombs away
		email := common.Email{
			From:    "luislimon@gmail.com",
			To:      "luislimon@gmail.com",
			Subject: fmt.Sprintf("ChurnDr: Pods with errors and not yet reported since %v minutes\n", LookBackTime),
			Body:    out,
		}
		fmt.Println(email)
		if !common.Config.DissableEmailNotifications {
			util.SendTLSMail(email)
		}
	}
}

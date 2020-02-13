package util

import (
	"fmt"

	"github.com/llimon/churndr/common"
	corev1 "k8s.io/api/core/v1"
)

func PersistPodLogs(pod *corev1.Pod, containerName string, restartCount int32) {
	log, err := GetPreviousPodLogs(pod, containerName, 200, 4096)
	if err != nil {
		common.Sugar.Infow("Unable to get logs for pod container",
			"pod", pod.Name,
			"namespace", pod.Namespace,
			"container", containerName,
			"err", err.Error())
	} else {
		// Persist lastTimeReported
		p := common.PodLogsDB{
			Name:      pod.ObjectMeta.Name,
			Namespace: pod.ObjectMeta.Namespace,
			Log:       log,
		}
		p.SetLink("self", "/pod/log/"+string(restartCount)+"/"+pod.ObjectMeta.Namespace+"/"+pod.ObjectMeta.Name+"/"+fmt.Sprint(restartCount), "")
		common.PodLogs[pod.ObjectMeta.Name+"/"+containerName+"/"+fmt.Sprint(restartCount)] = p

		// Persist History of executions.
		pHistory := common.PodLogHistory[pod.ObjectMeta.Name+"/"+containerName]

		pHistory = append(pHistory, common.PodLogHistoryDB{
			Name:         pod.ObjectMeta.Name,
			Namespace:    pod.ObjectMeta.Namespace,
			RestartCount: restartCount,
		})
		if len(pHistory) > 3 {
			common.Sugar.Infof("Deleting old Pod run history")
			pHistory = pHistory[1:]
		}
		common.PodLogHistory[pod.ObjectMeta.Name+"/"+containerName] = pHistory

		common.Sugar.Infow("Got Logs for container",
			"pod", pod.Name,
			"container", containerName,
			"namespace", pod.Namespace,
			"container", containerName,
			"restartCount", restartCount,
			"log", log)

	}
}

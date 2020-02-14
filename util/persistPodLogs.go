package util

import (
	"fmt"
	"time"

	"github.com/llimon/churndr/common"
	corev1 "k8s.io/api/core/v1"
)

func PersistPodLogs(pod *corev1.Pod, containerStatus corev1.ContainerStatus) { //containerName string, restartCount int32) {
	log, err := GetPreviousPodLogs(pod, containerStatus.Name, 200, 4096)
	if err != nil {
		common.Sugar.Infow("Unable to get logs for pod container",
			"pod", pod.Name,
			"namespace", pod.Namespace,
			"container", containerStatus.Name,
			"err", err.Error())
	} else {
		// Persist lastTimeReported
		p := common.PodLogsDB{
			Name:      pod.ObjectMeta.Name,
			Namespace: pod.ObjectMeta.Namespace,
			Log:       log,
		}
		p.SetLink("self", "/pod/log/container/"+pod.ObjectMeta.Namespace+"/"+pod.ObjectMeta.Name+"/"+fmt.Sprint(containerStatus.RestartCount), "")
		common.PodLogs[pod.ObjectMeta.Name+"/"+containerStatus.Name+"/"+fmt.Sprint(containerStatus.RestartCount)] = p

		// Persist History of executions.
		pHistory := common.PodLogHistory[pod.ObjectMeta.Namespace+"/"+pod.ObjectMeta.Name]

		c := common.PodLogHistoryDB{
			Name:         pod.ObjectMeta.Name,
			Container:    containerStatus.Name,
			Namespace:    pod.ObjectMeta.Namespace,
			RestartCount: containerStatus.RestartCount,
			FinishedAt:   time.Time(containerStatus.State.Terminated.FinishedAt.Time),
			ExitCode:     containerStatus.State.Terminated.ExitCode,
			Signal:       containerStatus.State.Terminated.Signal,
			Reason:       containerStatus.State.Terminated.Reason,
		}
		c.SetLink("log", "/pod/log/container/"+pod.ObjectMeta.Namespace+"/"+pod.ObjectMeta.Name+"/"+containerStatus.Name+"/"+fmt.Sprint(containerStatus.RestartCount), "")

		pHistory = append(pHistory, c)
		if len(pHistory) > common.Config.MaxPodHistory {
			common.Sugar.Infof("Deleting old Pod run history")
			pHistory = pHistory[1:]
		}
		common.PodLogHistory[pod.ObjectMeta.Namespace+"/"+pod.ObjectMeta.Name] = pHistory

		common.Sugar.Infow("Got Logs for container",
			"pod", pod.Name,
			"container", containerStatus.Name,
			"namespace", pod.Namespace,
			"container", containerStatus.Name,
			"restartCount", containerStatus.RestartCount,
			"log", log)

	}
}

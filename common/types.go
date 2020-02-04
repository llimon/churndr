package common

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

type ContainerDB struct {
	Name         string                           `json:"name"`
	Running      *corev1.ContainerStateRunning    `json:"running,omitempty"`
	Waiting      *corev1.ContainerStateWaiting    `json:"waiting,omitempty"`
	Terminated   *corev1.ContainerStateTerminated `json:"terminated,omitempty"`
	RestartCount int32
}

type PodDB struct {
	Name             string        `json:"name"`
	Namespace        string        `json:"namespace"`
	LastTimeReported time.Time     `json:"lasttimereported,omitempty"`
	Container        []ContainerDB `json:"container,omitempty"`
}

type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

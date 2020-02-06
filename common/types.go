package common

import (
	"time"

	"github.com/RichardKnop/jsonhal"
	corev1 "k8s.io/api/core/v1"
)

// Configuration is the configuration needed to run API server
type Configuration struct {
	InClusterConfiguration        bool
	Port                          int
	Development                   bool
	MonitorCurrentNamespace       bool
	Namespaces                    []string // List of namespaces to monitor if MonitorCurrentNamespace==false
	DissableEmailNotifications    bool
	NoiseReductionLookBackMinutes int
	NotificationFrequency         int
	NoAPIServer                   bool
	EmailSMTPServer               string
	EmailFrom                     string
	EmailTo                       string
	EmailSubject                  string
	EmailLogin                    string
	EmailPassword                 string
}

type ContainerDB struct {
	Name         string                           `json:"name"`
	Running      *corev1.ContainerStateRunning    `json:"running,omitempty"`
	Waiting      *corev1.ContainerStateWaiting    `json:"waiting,omitempty"`
	Terminated   *corev1.ContainerStateTerminated `json:"terminated,omitempty"`
	RestartCount int32
}

type Status struct {
	jsonhal.Hal
	Name string `json:"Name"`
}

type PodDB struct {
	jsonhal.Hal
	Name             string        `json:"name"`
	Namespace        string        `json:"namespace"`
	LastTimeReported time.Time     `json:"lasttimereported,omitempty"`
	Reported         bool          `json:"reported,omitempty"`
	Container        []ContainerDB `json:"container,omitempty"`
}

type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

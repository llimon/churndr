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
	jsonhal.Hal
	Name           string                           `json:"name"`
	Running        *corev1.ContainerStateRunning    `json:"running,omitempty"`
	Waiting        *corev1.ContainerStateWaiting    `json:"waiting,omitempty"`
	Terminated     *corev1.ContainerStateTerminated `json:"terminated,omitempty"`
	Logs           string                           `json:"logs,omitempty"`
	RestartCount   int32
	TerminationLog string
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
	RecoveredAt      time.Time     `json:"recoveredat,omitempty"`
	IsHealthy        bool          `json:"ishealthy"`
}

type PodLogsDB struct {
	jsonhal.Hal
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	RestartCount int32  `json:"restartcount"`
	Log          string `json:"log"`
}

type PodLogHistoryDB struct {
	Name         string    `json:"name"`
	Namespace    string    `json:"namespace"`
	RestartCount int32     `json:"restartCount"`
	Reason       string    `json:"reason"`
	ExitCode     int32     `json:"exitcCode"`
	terminatedAt time.Time `json:"terminatedAt"`
}

type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

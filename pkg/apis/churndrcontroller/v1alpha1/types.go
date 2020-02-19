/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/robfig/cron"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Podchurn is a specification for a Podchurn resource
type Podchurn struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodchurnSpec   `json:"spec"`
	Status PodchurnStatus `json:"status"`
}

// PodchurnSpec is the spec for a Podchurn resource
type PodchurnSpec struct {
	PodFilters             []*PodFilters `json:"podFilters"`
	PodLogs                PodLogs       `json:"podLogs"`
	IgnoreFinalTermination bool          `json:"ignoreFinalTermination"`
	DeploymentName         string        `json:"deploymentName"`
	Replicas               *int32        `json:"replicas"`
}

// PodFilters - Contains filtering used to select pods
type PodFilters struct {
	Name       string `json:"name"`
	MatchRegEx string `json:"matchRegEx"`
}

// AWSKeys - Contains Struct for storing AWS AccessKey and SecretKey
type AWSKeys struct {
	Secret      string `json:"secret"`
	AccessKeyID string `json:"accessKeyID"` // AccessKeyID in Kubernetes Secret
	SecretKeyID string `json:"secretKeyID"` // SecretKeyID in kubernetes secret
	AccessKey   string `json:"-"`           // Decripted secret AccessKey
	SecretKey   string `json:"-"`           // Decrypted secret SecretKey

}

// PodLogs - Contains information about how to handdle logs
type PodLogs struct {
	SaveLogs    bool    `json:"saveLogs`
	TailLines   *int32  `json:"tailLines"`
	MaxBytes    *int32  `json:"tailLines"`
	StorageType string  `json:"storageType"`
	bucket      string  `json:"bucket"`
	AwsKeys     AWSKeys `json:"awsKeys"`
}

// FileWatcher - Defines File watchers used for launching Workflows
// BUG/FIX: mode Triggers to their own struct
type FileWatcher struct {
	Type           string  `json:"type"`           // Type = "S3"
	Bucket         string  `json:"bucket"`         // Bucket = s3://hello-bucket
	Region         string  `json:"region"`         // AWS region of bucket
	AWSKeys        AWSKeys `json:"AWSKeys"`        // AWS AccessKeys and SecretKey
	Path           string  `json:"path"`           // Path = my-trigger.todo
	LaunchExt      string  `json:"launchExt"`      // LaunchExt = `TODO`  = file used to start a job
	FinishExt      string  `json:"finishExt"`      // FinishExt = 'DONE'  = to mark work as complete
	ProcessingExt  string  `json:"processingExt"`  // ProcessingExt = 'PROCESSING' used to mark work in progress.
	RenameOnStart  bool    `json:"renameOnStart"`  // File will be renamed from LanchExt to ProcessingExt when a job is triggered
	DeleteOnStart  bool    `json:"deleteOnStart"`  // Delete LaunchExt on Start of job
	RenameOnFinish bool    `json:"renameOnFinish"` // File will be renamed from ProcessingExt to FinishExt when job completes
	DeleteOnFinish bool    `json:"deleteOnFinish"` // Delete ProcessingExt on Finish of job
	Frequency      string  `json:"frequency"`      // When to run filewatcher "@every 5m", "@hourly", "@every 1h30m"
}

// WorkFlow - Defines information needed for geting workflow definitions and how to launch them
// BUG/FIX: Should read GITHub UserId and Password from a kubernetes Secret
// BUG/FIX: Implement Workflow parameters. Should be rendered by golang text/template (Use my template.go from kubeam)
//           Add to template.go:  ${.today()}, ${.yesterday()}, ${.firstdayOfMonth()}, ${.lastDayOfMonth{}}
//                                ${.prevWorkDay()}, ${.IsLeapyear()}, ${.daysInFebruary()}
//           Functions will be usefull seting parameters like REPORTDATE=${.yesterday()}.
//           This is a powerfull functionality expected in Enterprise schedulers.
type WorkFlow struct {
	SourceType         string   `json:"sourceType"` // SourceType = "github"
	Repo               string   `json:"repo"`       // Repo = "github.com/llimon/schedule-samples"
	Path               string   `json:"path"`       // Path = my-workflow.yaml
	UserID             string   `json:"userID"`     //usedID = llimon
	Password           string   `json:"password"`   // password = quack
	WorkflowParameters []string `json:"workflowParameters"`
}

type Activity struct {
	Status              string       `json:"status"`              // status = active, failed, running, launching
	CronEntryID         cron.EntryID `json:"_,omitempty"`         // Schedule IS assigned by cron
	LastExecutedJobName string       `json:"lastExecutedJobName"` // Name of the last job create launched by ArgoWorkflows
	FileWatcherEntryID  cron.EntryID `json:"_,omitempty"`         // Job used to running Filewatcher
}

// PodchurnStatus is the status for a Podchurn resource
type PodchurnStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodchurnList is a list of Podchurn resources
type PodchurnList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Podchurn `json:"items"`
}

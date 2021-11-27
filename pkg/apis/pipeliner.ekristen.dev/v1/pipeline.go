package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rancher/wrangler/pkg/genericcondition"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PipelineSpec   `json:"spec,omitempty"`   // Spec defines the desired state of the underlying device
	Status            PipelineStatus `json:"status,omitempty"` // Status defines the current state of the underlying device
}

type PipelineSpec struct {
	SourceRef ObjectReference        `json:"sourceRef"`
	Env       []corev1.EnvVar        `json:"env,omitempty"`
	EnvFrom   []corev1.EnvFromSource `json:"envFrom,omitempty"`
	DependsOn []ObjectReference      `json:"dependsOn,omitempty`
}

type PipelineStatus struct {
	Conditions []genericcondition.GenericCondition `json:"conditions,omitempty"`

	State      State           `json:"state" wrangler:"default=initializing"`
	StartedAt  metav1.Time     `json:"startedAt,omitempty"`
	FinishedAt metav1.Time     `json:"finishedAt,omitempty"`
	Duration   metav1.Duration `json:"duration,omitempty"`
	Stages     []PipelineStage `json:"stages"`
}

type PipelineStage struct {
	Name  string `json:"name"`
	State State  `json:"state"`
	Jobs  int    `json:"jobs"`
}

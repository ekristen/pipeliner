package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rancher/wrangler/pkg/genericcondition"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Workflow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WorkflowSpec   `json:"spec,omitempty"`   // Spec defines the desired state of the underlying device
	Status            WorkflowStatus `json:"status,omitempty"` // Status defines the current state of the underlying device
}

type WorkflowSpec struct {
	Raw     string                 `json:"raw"`
	Env     []corev1.EnvVar        `json:"env,omitempty"`
	EnvFrom []corev1.EnvFromSource `json:"envFrom,omitempty"`
}

type WorkflowStatus struct {
	Conditions []genericcondition.GenericCondition `json:"conditions,omitempty"`
}

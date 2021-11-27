package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rancher/wrangler/pkg/genericcondition"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Job struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              JobSpec   `json:"spec,omitempty"`   // Spec defines the desired state of the underlying device
	Status            JobStatus `json:"status,omitempty"` // Status defines the current state of the underlying device
}

type JobSpec struct {
	SecretRef ObjectReference   `json:"secretRef"`
	DependsOn []ObjectReference `json:"dependsOn,omitempty`
}

type JobStatus struct {
	Conditions []genericcondition.GenericCondition `json:"conditions,omitempty"`
}

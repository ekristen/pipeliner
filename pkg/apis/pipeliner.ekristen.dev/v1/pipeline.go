package v1

import (
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
	SecretRef ObjectReference `json:"secretRef"`
}

type PipelineStatus struct {
	Conditions []genericcondition.GenericCondition `json:"conditions,omitempty"`
}

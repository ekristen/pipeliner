package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/rancher/wrangler/pkg/genericcondition"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Artifact struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ArtifactSpec   `json:"spec,omitempty"`   // Spec defines the desired state of the underlying device
	Status            ArtifactStatus `json:"status,omitempty"` // Status defines the current state of the underlying device
}

type ArtifactSpec struct {
	Job ObjectReference `json:"job"`
}

type ArtifactStatus struct {
	Conditions []genericcondition.GenericCondition `json:"conditions,omitempty"`

	Type      string      `json:"type,omitempty"`
	Format    string      `json:"format,omitempty"`
	File      string      `json:"file,omitempty"`
	Size      int64       `json:"size,omitempty"`
	ExpiresAt metav1.Time `json:"expires_at,omitempty"`
}

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type State string

func (s State) String() string {
	return string(s)
}

const (
	VALID State = "valid"

	EMPTY        State = ""
	INITIALIZING State = "initializing"
	INITIALIZED  State = "initalized"
	PENDING      State = "pending"
	RUNNING      State = "running"
	SUCCESS      State = "success"
	FAILED       State = "failed"
	SKIPPED      State = "skipped"
	CANCELED     State = "canceled"
)

// ObjectReference is a reference to an object with a given name, kind and group.
type ObjectReference struct {
	// Name of the resource being referred to.
	Name string `json:"name"`
	// Namespace of the resource being referred to.
	Namespace string `json:"namespace"`
	// Kind of the resource being referred to.
	// +optional
	Kind string `json:"kind,omitempty"`
	// Group of the resource being referred to.
	// +optional
	Group string `json:"group,omitempty"`
	// Selector for matching resources, conflicts with Name, only useful for DependsOn
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

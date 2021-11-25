package v1

type State string

const (
	VALID State = "valid"
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
}

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Setting struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Value   string `json:"value" wrangler:"required" column:"name=value,type=string,jsonpath=.value"`
	Default string `json:"default" wrangler:"nocreate,noupdate" column:"name=default/fallback,type=string,jsonpath=.default"`
}

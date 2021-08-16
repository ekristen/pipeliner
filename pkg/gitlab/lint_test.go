package gitlab

import (
	"io/ioutil"
	"testing"
)

func TestWorkflowInvalidStage(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/missing-stage.yaml")
	if err != nil {
		t.Error(err)
	}

	validator := NewYAMLValidator(data)
	validator.Parse()
	if err := validator.Validate(); err != nil {
		t.Error(err)
	}
}

package controllers

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"gitlab.com/gitlab-org/gitlab-runner/helpers/gitlab_ci_yaml_parser"
	"gopkg.in/yaml.v1"
)

func TestGenerateBuilds1(t *testing.T) {
	expectedBuilds := []MatrixBuild{
		{
			Name: "testing 1/3",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "0",
				"CI_NODE_TOTAL": "3",
			},
		},
		{
			Name: "testing 2/3",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "1",
				"CI_NODE_TOTAL": "3",
			},
		},
		{
			Name: "testing 3/3",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "2",
				"CI_NODE_TOTAL": "3",
			},
		},
	}

	data, err := ioutil.ReadFile("testdata/parallel.yaml")
	if err != nil {
		t.Error(err)
	}

	config := make(gitlab_ci_yaml_parser.DataBag)
	if err := yaml.Unmarshal(data, config); err != nil {
		t.Error(err)
	}

	if err := config.Sanitize(); err != nil {
		t.Error(err)
	}

	d, _ := config.GetSubOptions("testing")

	builds := generateBuilds("testing", d)
	if !reflect.DeepEqual(builds, expectedBuilds) {
		t.Error("builds do not match")
		fmt.Println("Actual:", builds)
		fmt.Println("Expected:", expectedBuilds)
	}
}

func TestGenerateBuilds2(t *testing.T) {
	expectedBuilds := []MatrixBuild{
		{
			Name: "testing: [aws, monitoring]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "0",
				"CI_NODE_TOTAL": "3",
				"PROVIDER":      "aws",
				"STACK":         "monitoring",
			},
		},
		{
			Name: "testing: [aws, app1]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "1",
				"CI_NODE_TOTAL": "3",
				"PROVIDER":      "aws",
				"STACK":         "app1",
			},
		},
		{
			Name: "testing: [aws, app2]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "2",
				"CI_NODE_TOTAL": "3",
				"PROVIDER":      "aws",
				"STACK":         "app2",
			},
		},
	}

	data, err := ioutil.ReadFile("testdata/matrix.yaml")
	if err != nil {
		t.Error(err)
	}

	config := make(gitlab_ci_yaml_parser.DataBag)
	if err := yaml.Unmarshal(data, config); err != nil {
		t.Error(err)
	}

	if err := config.Sanitize(); err != nil {
		t.Error(err)
	}

	d, _ := config.GetSubOptions("testing")

	actualBuilds := generateBuilds("testing", d)

	if !reflect.DeepEqual(actualBuilds, expectedBuilds) {
		t.Errorf("builds are not equal")
	}
}

func TestGenerateBuilds3(t *testing.T) {
	expectedBuilds := []MatrixBuild{
		{
			Name: "testing: [aws, monitoring]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "0",
				"CI_NODE_TOTAL": "3",
				"PROVIDER":      "aws",
				"STACK":         "monitoring",
			},
		},
		{
			Name: "testing: [aws, app1]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "1",
				"CI_NODE_TOTAL": "3",
				"PROVIDER":      "aws",
				"STACK":         "app1",
			},
		},
		{
			Name: "testing: [aws, app2]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "2",
				"CI_NODE_TOTAL": "3",
				"PROVIDER":      "aws",
				"STACK":         "app2",
			},
		},
		{
			Name: "testing: [ovh, monitoring]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "0",
				"CI_NODE_TOTAL": "3",
				"PROVIDER":      "ovh",
				"STACK":         "monitoring",
			},
		},
		{
			Name: "testing: [ovh, backup]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "1",
				"CI_NODE_TOTAL": "3",
				"PROVIDER":      "ovh",
				"STACK":         "backup",
			},
		},
		{
			Name: "testing: [ovh, app]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "2",
				"CI_NODE_TOTAL": "3",
				"PROVIDER":      "ovh",
				"STACK":         "app",
			},
		},
		{
			Name: "testing: [gcp, data]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "0",
				"CI_NODE_TOTAL": "4",
				"PROVIDER":      "gcp",
				"STACK":         "data",
			},
		},
		{
			Name: "testing: [gcp, processing]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "1",
				"CI_NODE_TOTAL": "4",
				"PROVIDER":      "gcp",
				"STACK":         "processing",
			},
		},
		{
			Name: "testing: [vultr, data]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "2",
				"CI_NODE_TOTAL": "4",
				"PROVIDER":      "vultr",
				"STACK":         "data",
			},
		},
		{
			Name: "testing: [vultr, processing]",
			Job:  "testing",
			Env: map[string]string{
				"CI_NODE_INDEX": "3",
				"CI_NODE_TOTAL": "4",
				"PROVIDER":      "vultr",
				"STACK":         "processing",
			},
		},
	}

	data, err := ioutil.ReadFile("testdata/matrix2.yaml")
	if err != nil {
		t.Error(err)
	}

	config := make(gitlab_ci_yaml_parser.DataBag)
	if err := yaml.Unmarshal(data, config); err != nil {
		t.Error(err)
	}

	if err := config.Sanitize(); err != nil {
		t.Error(err)
	}

	d, _ := config.GetSubOptions("testing")

	actualBuilds := generateBuilds("testing", d)
	if !reflect.DeepEqual(actualBuilds, expectedBuilds) {
		t.Errorf("builds are not equal")
	}
}

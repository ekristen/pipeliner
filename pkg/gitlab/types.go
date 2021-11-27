package gitlab

import (
	"strings"
	"time"
)

var DefaultNames = []string{"after_script", "artifact", "before_script", "cache", "image", "interruptible", "retry", "services", "tags", "timeout", "variables"}
var ReservedNames = append([]string{"default", "stages", "include"}, DefaultNames...)

var DefaultStages = []string{"build", "test", "publish"}

type Default struct {
	AfterScript   []string      `json:"after_script,omitempty"`
	Artifact      []Artifact    `json:"artifacts,omitempty"`
	BeforeScript  []string      `json:"before_script,omitempty"`
	Cache         Cache         `json:"cache,omitempty"`
	Image         Image         `json:"image,omitempty`
	Interruptible bool          `json:"interruptible"`
	Retry         Retry         `json:"retry,omitempty"`
	Services      []Service     `json:"services,omitempty"`
	Tags          []string      `json:"tags,omitempty"`
	Timeout       time.Duration `json:"timeout" default:"3600s"`
}

type Artifact struct {
	Name      string         `json:"name,omitempty"`
	Paths     []string       `json:"paths,omitempty"`
	Exclude   []string       `json:"exclude,omitempty"`
	ExpireIn  time.Duration  `json:"expire_in,omitempty"`
	Reports   ArtifactReport `json:"reports,omitempty"`
	Untracked bool           `json:"untracked"`
	When      string         `json:"when" default:"on_success"`
}

type ArtifactReport struct {
	JUnit string `json:"junit,omitempty"`
}

type Cache struct {
	Key       CacheKey `json:"key"`
	Paths     []string `json:"paths,omitempty"`
	Untracked bool     `json:"untracked"`
	When      string   `json:"when" default:"on_success"`
	Policy    string   `json:"policy" default:"pull-push" validate:"(pull|push|pull-push)"`
}

type CacheKey struct {
	Name   string   `json:"name" default:"cache"` // this is normally just Cache.Key
	Files  []string `json:"files,omitempty"`
	Prefix string   `json:"prefix,omitempty"`
}

type Image struct {
	Name       string   `json:"name"` // this is normally just Image.Name
	Entrypoint []string `json:"entrypoint,omitempty"`
}

type Retry struct {
	Max  int    `json:"max" default:"0"`
	When string `json:"when" default:"always"` // validate with https://docs.gitlab.com/ee/ci/yaml/#retry
}

func (r *Retry) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var retry Retry
	if err := unmarshal(&retry); err != nil {
		var retrySingle int
		if err := unmarshal(&retrySingle); err != nil {
			return nil
		}
		*r = Retry{
			Max: retrySingle,
		}
	}

	*r = retry

	return nil
}

type Service struct {
	Name       string            `json:"name"`
	Alias      string            `json:"alias,omitempty"`
	Entrypoint []string          `json:"entrypoint,omitempty"`
	Command    []string          `json:"command,omitempty"`
	Variables  map[string]string `json:"variables,omitempty"`
}

type Pipeline struct {
	Stages []string `json:"stages,omitempty"`
	Jobs   []Job    `json:"jobs,omitempty"`
}

type Job struct {
	Default

	Name         string       `json:"name"`
	Token        string       `json:"token"`
	Stage        string       `json:"stage,omitempty"`
	State        string       `json:"state" default:"created"`
	AllowFailure AllowFailure `json:"allow_failure,omitempty"`
	Parallel     Parallel     `json:"parallel,omitempty`
	Script       []string     `json:"script,omitempty"`
	Variables    []Variable   `json:"variables,omitempty"`
	When         string       `json:"when,omitempty" default:"on_success"`
	Dependencies []string     `json:"dependencies,omitempty"`
	Needs        []Need       `json:"needs,omitempty"`

	Data []byte `json:"data"`
}

type AllowFailure struct {
	Enabled   bool  `json:"enabled"`
	ExitCodes []int `json:"exit_codes,omitempty"`
}

type Need struct {
	Job       string `json:"job,omitempty"`
	Artifacts bool   `json:"artifacts"`
	Pipeline  string `json:"pipeline,omitempty"`
}

type Parallel struct {
	Count  int                 `json:"count,omitempty"`
	Matrix []map[string]string `json:"matrix,omitempty"`
}

type MatrixBuild struct {
	Name string
	Job  string
	Env  map[string]string
}

type Variable struct {
	Name   string
	Value  string
	Masked bool
	File   bool
}

func IsReserved(name string) bool {
	for _, reservedName := range ReservedNames {
		if reservedName == name {
			return true
		}
	}

	return false
}

func IsHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

package common

// Docs: https://docs.gitlab.com/ee/ci/yaml/

// GitlabGlobalConfig --
var GitlabGlobalConfig = []string{
	"image",
	"services",
	"before_script",
	"after_script",
	"tags",
	"cache",
	"artifacts",
	"retry",
	"timeout",
	"interruptible",
	// Not in the documentation
	"script",
	"stages",
}

// GitLabUnavailableJobNames are used to exclude job names
var GitLabUnavailableJobNames = []string{
	"image",
	"services",
	"stages",
	"types",
	"before_script",
	"after_script",
	"variables",
	"cache",
	"include",
	// Not in Docs (from Defaults)
	"tags",
	"cache",
	"artifacts",
	"retry",
	"timeout",
	"script",
	"interruptible",
}

// GitLabKeywordsGlobalDefaults define the global defaults for a pipeline
var GitLabKeywordsGlobalDefaults = []string{
	"image",
	"services",
	"before_script",
	"after_script",
	"tags",
	"cache",
	"artifacts",
	"retry",
	"timeout",
	"interruptible",
}

// GitLabDefaultStages defines the default stages if none are specified
var GitLabDefaultStages = []string{"build", "test", "deploy"}

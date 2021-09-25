package common

// NAME of the App
var NAME = "pipeliner"

// SUMMARY of the Version
var SUMMARY = "dirty"

// BRANCH of the Version
var BRANCH = "dev"

// VERSION of Release
var VERSION = "0.1.0"

// AppVersion --
var AppVersion AppVersionInfo

// AppVersionInfo --
type AppVersionInfo struct {
	Name    string
	Version string
	Branch  string
	Summary string
}

func init() {
	AppVersion = AppVersionInfo{
		Name:    NAME,
		Version: VERSION,
		Branch:  BRANCH,
		Summary: SUMMARY,
	}
}

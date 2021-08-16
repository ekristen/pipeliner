package common

import "time"

type contextKey string

// ContextKeyDb is used to set context value for db
const ContextKeyDb contextKey = "db"

// ContextKeyNode is used to set context value for node/snowflake
const ContextKeyNode contextKey = "node"

// ContextKeyWebsocket is used to set context value for passing websocket
const ContextKeyWebsocket contextKey = "websocket"

const (
	// BuildRunningOutdatedTimeout --
	BuildRunningOutdatedTimeout = 1 * time.Hour
	// BuildPendingOutdatedTimeout --
	BuildPendingOutdatedTimeout = 24 * time.Hour
	// BuildScheduledOutdatedTimeout --
	BuildScheduledOutdatedTimeout = 1 * time.Hour
	// BuildPendingStuckTimeout --
	BuildPendingStuckTimeout = 1 * time.Hour
)

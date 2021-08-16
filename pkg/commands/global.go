package commands

import (
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func globalFlags() []cli.Flag {
	globalFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "Log Level",
			EnvVars: []string{"LOGLEVEL", "LOG_LEVEL"},
			Value:   "info",
		},
		&cli.BoolFlag{
			Name:    "log-disable-colors",
			Usage:   "Disable Colors for Logging Output",
			EnvVars: []string{"LOG_DISABLE_COLORS"},
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "log-full-timestamp",
			Usage:   "Show Full Timestamp in Log Output",
			EnvVars: []string{"LOG_FULL_TIMESTAMP"},
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "log-caller",
			Usage:   "Log Caller Filename and Line Number",
			EnvVars: []string{"LOG_CALLER"},
			Value:   true,
		},
		&cli.Int64Flag{
			Name:    "node-id",
			Usage:   "Node ID",
			EnvVars: []string{"NODE_ID"},
			Value:   1,
		},
		&cli.StringFlag{
			Name:    "dsn",
			Usage:   "Database Connection String",
			EnvVars: []string{"DSN"},
			Value:   "pipeliner.sqlite",
		},
	}

	return globalFlags
}

func globalBefore(c *cli.Context) error {
	formatter := &logrus.TextFormatter{
		DisableColors: c.Bool("log-disable-color"),
		FullTimestamp: c.Bool("log-full-timestamp"),
	}
	if c.Bool("log-caller") {
		logrus.SetReportCaller(true)

		formatter.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		}
	}

	logrus.SetFormatter(formatter)

	switch c.String("log-level") {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	}

	return nil
}

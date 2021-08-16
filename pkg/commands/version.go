package commands

import (
	"fmt"

	"github.com/ekristen/pipeliner/pkg/common"
	"github.com/urfave/cli/v2"
)

type versionCommand struct{}

func (s *versionCommand) Execute(c *cli.Context) error {
	fmt.Println(common.AppVersion)
	return nil
}

func init() {
	cmd := versionCommand{}

	cliCmd := &cli.Command{
		Name:   "version",
		Usage:  "pipeliner version information",
		Action: cmd.Execute,
		Before: globalBefore,
	}

	common.RegisterCommand(cliCmd)
}

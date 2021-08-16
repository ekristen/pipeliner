package commands

import (
	"context"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/ekristen/pipeliner/pkg/common"
	"github.com/ekristen/pipeliner/pkg/database"
	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/ekristen/pipeliner/pkg/utils"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type seedWorkflowsCommand struct{}

func (s *seedWorkflowsCommand) Execute(c *cli.Context) error {
	nodeID := c.Int64("node-id")

	ctx := signals.SetupSignalHandler(c.Context)

	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		return err
	}

	db, err := database.New("mysql", c.String("dsn"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	var notifyDatabaseChange chan []byte = make(chan []byte)

	dbctx := context.WithValue(ctx, common.ContextKeyNode, node)
	dbctx = context.WithValue(dbctx, common.ContextKeyWebsocket, notifyDatabaseChange)

	db = db.WithContext(dbctx)

	return seedWorkflows(db, c.String("directory"))
}

func init() {
	cmd := seedWorkflowsCommand{}

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "directory",
			Usage:   "workflow directory",
			EnvVars: []string{"WORKFLOW_DIR"},
			Value:   "workflows",
		},
	}

	cliCmd := &cli.Command{
		Name:   "seed-workflows",
		Usage:  "seed-workflows for pipeliner",
		Action: cmd.Execute,
		Flags:  append(flags, globalFlags()...),
		Before: globalBefore,
	}

	common.RegisterCommand(cliCmd)
}

func seedWorkflows(db *gorm.DB, directoryPath string) error {
	matches, err := filepath.Glob(path.Join(directoryPath, "*.yaml"))
	if err != nil {
		return err
	}
	for _, file := range matches {
		fileName := path.Base(file)
		n := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		c, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		if _, err := utils.CreateOrUpdate(db, &models.Workflow{}, db.Where("name = ?", n), &models.Workflow{
			Name: n,
			Data: string(c),
		}); err != nil {
			return err
		}
	}
	return nil
}

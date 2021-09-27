package commands

import (
	"context"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"gorm.io/gorm"

	"github.com/rancher/wrangler/pkg/signals"

	"github.com/ekristen/pipeliner/pkg/api"
	"github.com/ekristen/pipeliner/pkg/common"
	"github.com/ekristen/pipeliner/pkg/controllers"
	"github.com/ekristen/pipeliner/pkg/database"
	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/ekristen/pipeliner/pkg/store"
	"github.com/ekristen/pipeliner/pkg/workers"
)

var db *gorm.DB
var err error
var node *snowflake.Node

var notifyBuilds chan *models.Pipeline = make(chan *models.Pipeline)
var notifyPipelineRun chan *models.Pipeline = make(chan *models.Pipeline)

var pipelineWorker *workers.PipelineWorker
var buildsWorker *workers.BuildsWorker

var pipelinesController *controllers.Pipelines
var buildsController *controllers.Builds

type apiServerCommand struct{}

func (s *apiServerCommand) Execute(c *cli.Context) error {
	cookieSecret := c.String("cookie-secret")
	dsn := c.String("dsn")
	nodeID := c.Int64("node-id")
	port := c.Int("port")
	dbDriver := c.String("database-driver")

	ctx := signals.SetupSignalHandler(c.Context)

	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		return err
	}

	db, err := database.New(dbDriver, dsn, &gorm.Config{
		Logger: database.NewLogger(c.String("log-level")),
	})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&models.Artifact{},
		&models.Runner{},
		&models.RunnerTag{},
		&models.Workflow{},
		&models.Pipeline{},
		&models.PipelineStage{},
		&models.Build{},
		&models.BuildTag{},
		&models.RegisterToken{},
		&models.GlobalVariable{},
		&models.WorkflowVariable{},
		&models.PipelineVariable{},
		&models.BuildVariable{},
		&models.ScheduleVariable{},
		&models.Schedule{},
		&models.Setting{},
		&models.Trace{},
		&models.TracePart{},
	); err != nil {
		return err
	}

	storage, err := store.NewUploader(c.String("storage-driver"), c.String("storage-bucket"), c.String("storage-prefix"), c.String("storage-region"), c.String("storage-endpoint"), c.String("storage-sse"))
	if err != nil {
		return err
	}

	var notifyDatabaseChange chan []byte = make(chan []byte)

	dbctx := context.WithValue(ctx, common.ContextKeyNode, node)
	dbctx = context.WithValue(dbctx, common.ContextKeyWebsocket, notifyDatabaseChange)

	db = db.WithContext(dbctx)

	pipelinesController = controllers.NewPipelines(db)
	buildsController = controllers.NewBuilds(db, pipelinesController, notifyPipelineRun)

	pipelineWorker = workers.NewPipelineWorker(db, pipelinesController, buildsController)
	buildsWorker = workers.NewBuildsWorker(db, pipelineWorker, notifyPipelineRun)

	apiServer := api.NewServer(ctx, db, port, cookieSecret, storage, pipelinesController, buildsController, notifyPipelineRun, notifyBuilds, notifyDatabaseChange)

	go createBuildsWorker(ctx, notifyBuilds)
	go pipelinesWorker(ctx, notifyPipelineRun)

	if err := apiServer.Run(); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func pipelinesWorker(ctx context.Context, notify chan *models.Pipeline) {
	log := logrus.WithField("component", "worker").WithField("worker", "pipelines")
	for {
		select {
		case pipeline := <-notify:
			log.Info("received notification")
			if err := pipelineWorker.Perform(pipeline.ID); err != nil {
				log.WithError(err).Error("error processing pipeline")
			}
		case <-ctx.Done():
			log.Info("Shutting down pipelines worker loop")
			return
		default:
			time.Sleep(time.Second * 1)
		}
	}
}

func createBuildsWorker(ctx context.Context, notify chan *models.Pipeline) {
	log := logrus.WithField("component", "worker").WithField("worker", "builds")
	for {
		select {
		case pipeline := <-notify:
			log.Info("received notification")
			buildsWorker.Perform(pipeline)
		case <-ctx.Done():
			log.Info("Shutting down builds worker loop")
			return
		default:
			time.Sleep(time.Second * 1)
		}
	}
}

func init() {
	cmd := apiServerCommand{}

	flags := []cli.Flag{
		&cli.IntFlag{
			Name:    "port",
			Usage:   "Port for the HTTP Server to Listen On",
			EnvVars: []string{"PORT"},
			Value:   4444,
		},
		&cli.StringFlag{
			Name:    "cookie-secret",
			Usage:   "Cookie Secret",
			EnvVars: []string{"COOKIE_SECRET"},
			Value:   "wD7IesRJ2OhK6ufz-p-nig==",
		},
		&cli.StringFlag{
			Name:    "database-driver",
			Usage:   "Database Driver",
			EnvVars: []string{"DATABASE_DRIVER"},
			Value:   "sqlite",
		},
		&cli.StringFlag{
			Name:    "storage-driver",
			Usage:   "Storage Driver (local, aws, gcp, azure)",
			EnvVars: []string{"STORAGE_DRIVER"},
			Value:   "local",
		},
		&cli.StringFlag{
			Name:    "storage-bucket",
			Usage:   "Storage Bucket (only applies to AWS, Azure, GCP and not local)",
			EnvVars: []string{"STORAGE_BUCKET"},
			Value:   "pipeliner",
		},
		&cli.StringFlag{
			Name:    "storage-prefix",
			Usage:   "Storage Path Prefix",
			EnvVars: []string{"STORAGE_PREFIX"},
			Value:   "artifacts",
		},
		&cli.StringFlag{
			Name:    "storage-region",
			Usage:   "Storage Region (only applies to AWS S3)",
			EnvVars: []string{"STORAGE_REGION"},
			Value:   "us-east-1",
		},
		&cli.StringFlag{
			Name:    "storage-endpoint",
			Usage:   "Storage Endpoint (only applies to AWS S3)",
			EnvVars: []string{"STORAGE_ENDPOINT"},
			Value:   "s3.us-east-1.amazonaws.com",
		},
		&cli.StringFlag{
			Name:    "storage-sse",
			Usage:   "Storage SSE (only applies to AWS S3)",
			EnvVars: []string{"STORAGE_SSE"},
			Value:   "aws:kms",
		},
	}

	cliCmd := &cli.Command{
		Name:   "api-server",
		Usage:  "api-server for pipeliner",
		Action: cmd.Execute,
		Flags:  append(flags, globalFlags()...),
		Before: globalBefore,
	}

	common.RegisterCommand(cliCmd)
}

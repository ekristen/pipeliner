package commands

import (
	"context"

	"github.com/ekristen/pipeliner/pkg/common"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/generated/controllers/core"
	"github.com/rancher/wrangler/pkg/kubeconfig"
	"github.com/rancher/wrangler/pkg/leader"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/rancher/wrangler/pkg/start"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"k8s.io/apimachinery/pkg/util/runtime"
)

type controllerCommand struct{}

func (s *controllerCommand) Execute(c *cli.Context) error {
	// set up signals so we handle the first shutdown signal gracefully
	ctx := signals.SetupSignalHandler(context.Background())

	log := logrus.WithField("command", "controller")

	//go metrics.NewMetricsServer(ctx, c.String("metrics-port"), true, metrics.OdinRegistry)

	cfg, err := kubeconfig.GetNonInteractiveClientConfig(c.String("kubeconfig")).ClientConfig()
	if err != nil {
		return err
	}

	apply, err := apply.NewForConfig(cfg)
	if err != nil {
		return err
	}

	pipeliner, err := pipeliner.NewFactoryFromConfig(cfg)
	if err != nil {
		return err
	}

	if err := crds.Create(ctx, cfg); err != nil {
		return err
	}

	core, err := core.NewFactoryFromConfig(cfg)
	if err != nil {
		return err
	}

	// Register all our controllers
	if err := settings.Register(odin.Odin().V1().Setting()); err != nil {
		return err
	}

	// Become leader, then create CRDS (or update), followed by starting all controllers
	leader.RunOrDie(ctx, c.String("namespace"), c.String("lockname"), kube, func(ctx context.Context) {
		runtime.Must(start.All(ctx, 5, core, pipeliner))

		<-ctx.Done()
	})

	return nil
}

func init() {
	cmd := controllerCommand{}

	cliCmd := &cli.Command{
		Name:   "controller",
		Usage:  "controller for pipeliner k8s crds",
		Action: cmd.Execute,
		Before: globalBefore,
	}

	common.RegisterCommand(cliCmd)
}

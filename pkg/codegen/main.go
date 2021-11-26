package main

import (
	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"

	v1 "github.com/ekristen/pipeliner/pkg/apis/pipeliner.ekristen.dev/v1"
)

func main() {
	controllergen.Run(args.Options{
		OutputPackage: "github.com/ekristen/pipeliner/pkg/generated",
		Boilerplate:   "hack/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"pipeliner.ekristen.dev": {
				Types: []interface{}{
					v1.Setting{},
					v1.GitRepository{},
					v1.Workflow{},
				},
				GenerateTypes:   true,
				GenerateClients: true,
			},
		},
	})
}

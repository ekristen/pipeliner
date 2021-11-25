package crds

import (
	"context"

	v1 "github.com/ekristen/pipeliner/pkg/apis/pipeliner.ekristen.dev/v1"
	"github.com/rancher/wrangler/pkg/crd"
	"k8s.io/client-go/rest"
)

func List() []crd.CRD {
	return []crd.CRD{
		newCRD(v1.Setting{}, func(obj interface{}, c crd.CRD) crd.CRD {
			c.NonNamespace = true

			return c.
				WithCategories("pipeliner").
				WithColumnsFromStruct(obj)
		}),
		newCRD(v1.GitRepository{}, func(obj interface{}, c crd.CRD) crd.CRD {
			return c.WithStatus().
				WithShortNames("git").
				WithCategories("pipeliner").
				WithColumnsFromStruct(obj)
		}),
	}
}

/*
func Objects() (result []runtime.Object, err error) {
	for _, crdDef := range List() {
		crd, err := crdDef.ToCustomResourceDefinition()
		if err != nil {
			return nil, err
		}
		n := crd.DeepCopy()
		result = append(result, n)
	}
	return
}
*/

func Create(ctx context.Context, cfg *rest.Config) error {
	factory, err := crd.NewFactoryFromClient(cfg)
	if err != nil {
		return err
	}

	return factory.BatchCreateCRDs(ctx, List()...).BatchWait()
}

func newCRD(obj interface{}, customize func(interface{}, crd.CRD) crd.CRD) crd.CRD {
	crd := crd.CRD{
		SchemaObject: obj,
	}
	if customize != nil {
		crd = customize(obj, crd)
	}
	return crd
}

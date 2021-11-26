/*
Copyright 2021 Erik Kristensen

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/ekristen/pipeliner/pkg/apis/pipeliner.ekristen.dev/v1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type PipelineHandler func(string, *v1.Pipeline) (*v1.Pipeline, error)

type PipelineController interface {
	generic.ControllerMeta
	PipelineClient

	OnChange(ctx context.Context, name string, sync PipelineHandler)
	OnRemove(ctx context.Context, name string, sync PipelineHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() PipelineCache
}

type PipelineClient interface {
	Create(*v1.Pipeline) (*v1.Pipeline, error)
	Update(*v1.Pipeline) (*v1.Pipeline, error)
	UpdateStatus(*v1.Pipeline) (*v1.Pipeline, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.Pipeline, error)
	List(namespace string, opts metav1.ListOptions) (*v1.PipelineList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Pipeline, err error)
}

type PipelineCache interface {
	Get(namespace, name string) (*v1.Pipeline, error)
	List(namespace string, selector labels.Selector) ([]*v1.Pipeline, error)

	AddIndexer(indexName string, indexer PipelineIndexer)
	GetByIndex(indexName, key string) ([]*v1.Pipeline, error)
}

type PipelineIndexer func(obj *v1.Pipeline) ([]string, error)

type pipelineController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewPipelineController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) PipelineController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &pipelineController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromPipelineHandlerToHandler(sync PipelineHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Pipeline
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Pipeline))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *pipelineController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Pipeline))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdatePipelineDeepCopyOnChange(client PipelineClient, obj *v1.Pipeline, handler func(obj *v1.Pipeline) (*v1.Pipeline, error)) (*v1.Pipeline, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *pipelineController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *pipelineController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *pipelineController) OnChange(ctx context.Context, name string, sync PipelineHandler) {
	c.AddGenericHandler(ctx, name, FromPipelineHandlerToHandler(sync))
}

func (c *pipelineController) OnRemove(ctx context.Context, name string, sync PipelineHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromPipelineHandlerToHandler(sync)))
}

func (c *pipelineController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *pipelineController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *pipelineController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *pipelineController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *pipelineController) Cache() PipelineCache {
	return &pipelineCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *pipelineController) Create(obj *v1.Pipeline) (*v1.Pipeline, error) {
	result := &v1.Pipeline{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *pipelineController) Update(obj *v1.Pipeline) (*v1.Pipeline, error) {
	result := &v1.Pipeline{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *pipelineController) UpdateStatus(obj *v1.Pipeline) (*v1.Pipeline, error) {
	result := &v1.Pipeline{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *pipelineController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *pipelineController) Get(namespace, name string, options metav1.GetOptions) (*v1.Pipeline, error) {
	result := &v1.Pipeline{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *pipelineController) List(namespace string, opts metav1.ListOptions) (*v1.PipelineList, error) {
	result := &v1.PipelineList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *pipelineController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *pipelineController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.Pipeline, error) {
	result := &v1.Pipeline{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type pipelineCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *pipelineCache) Get(namespace, name string) (*v1.Pipeline, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.Pipeline), nil
}

func (c *pipelineCache) List(namespace string, selector labels.Selector) (ret []*v1.Pipeline, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Pipeline))
	})

	return ret, err
}

func (c *pipelineCache) AddIndexer(indexName string, indexer PipelineIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Pipeline))
		},
	}))
}

func (c *pipelineCache) GetByIndex(indexName, key string) (result []*v1.Pipeline, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.Pipeline, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.Pipeline))
	}
	return result, nil
}

type PipelineStatusHandler func(obj *v1.Pipeline, status v1.PipelineStatus) (v1.PipelineStatus, error)

type PipelineGeneratingHandler func(obj *v1.Pipeline, status v1.PipelineStatus) ([]runtime.Object, v1.PipelineStatus, error)

func RegisterPipelineStatusHandler(ctx context.Context, controller PipelineController, condition condition.Cond, name string, handler PipelineStatusHandler) {
	statusHandler := &pipelineStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromPipelineHandlerToHandler(statusHandler.sync))
}

func RegisterPipelineGeneratingHandler(ctx context.Context, controller PipelineController, apply apply.Apply,
	condition condition.Cond, name string, handler PipelineGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &pipelineGeneratingHandler{
		PipelineGeneratingHandler: handler,
		apply:                     apply,
		name:                      name,
		gvk:                       controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterPipelineStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type pipelineStatusHandler struct {
	client    PipelineClient
	condition condition.Cond
	handler   PipelineStatusHandler
}

func (a *pipelineStatusHandler) sync(key string, obj *v1.Pipeline) (*v1.Pipeline, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type pipelineGeneratingHandler struct {
	PipelineGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *pipelineGeneratingHandler) Remove(key string, obj *v1.Pipeline) (*v1.Pipeline, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.Pipeline{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *pipelineGeneratingHandler) Handle(obj *v1.Pipeline, status v1.PipelineStatus) (v1.PipelineStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.PipelineGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}

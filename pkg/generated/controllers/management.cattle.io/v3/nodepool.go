/*
Copyright 2024 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v3

import (
	"context"
	"sync"
	"time"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/shepherd/pkg/wrangler/pkg/generic"
	"github.com/rancher/wrangler/v3/pkg/apply"
	"github.com/rancher/wrangler/v3/pkg/condition"
	"github.com/rancher/wrangler/v3/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// NodePoolController interface for managing NodePool resources.
type NodePoolController interface {
	generic.ControllerInterface[*v3.NodePool, *v3.NodePoolList]
}

// NodePoolClient interface for managing NodePool resources in Kubernetes.
type NodePoolClient interface {
	generic.ClientInterface[*v3.NodePool, *v3.NodePoolList]
}

// NodePoolCache interface for retrieving NodePool resources in memory.
type NodePoolCache interface {
	generic.CacheInterface[*v3.NodePool]
}

// NodePoolStatusHandler is executed for every added or modified NodePool. Should return the new status to be updated
type NodePoolStatusHandler func(obj *v3.NodePool, status v3.NodePoolStatus) (v3.NodePoolStatus, error)

// NodePoolGeneratingHandler is the top-level handler that is executed for every NodePool event. It extends NodePoolStatusHandler by a returning a slice of child objects to be passed to apply.Apply
type NodePoolGeneratingHandler func(obj *v3.NodePool, status v3.NodePoolStatus) ([]runtime.Object, v3.NodePoolStatus, error)

// RegisterNodePoolStatusHandler configures a NodePoolController to execute a NodePoolStatusHandler for every events observed.
// If a non-empty condition is provided, it will be updated in the status conditions for every handler execution
func RegisterNodePoolStatusHandler(ctx context.Context, controller NodePoolController, condition condition.Cond, name string, handler NodePoolStatusHandler) {
	statusHandler := &nodePoolStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, generic.FromObjectHandlerToHandler(statusHandler.sync))
}

// RegisterNodePoolGeneratingHandler configures a NodePoolController to execute a NodePoolGeneratingHandler for every events observed, passing the returned objects to the provided apply.Apply.
// If a non-empty condition is provided, it will be updated in the status conditions for every handler execution
func RegisterNodePoolGeneratingHandler(ctx context.Context, controller NodePoolController, apply apply.Apply,
	condition condition.Cond, name string, handler NodePoolGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &nodePoolGeneratingHandler{
		NodePoolGeneratingHandler: handler,
		apply:                     apply,
		name:                      name,
		gvk:                       controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterNodePoolStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type nodePoolStatusHandler struct {
	client    NodePoolClient
	condition condition.Cond
	handler   NodePoolStatusHandler
}

// sync is executed on every resource addition or modification. Executes the configured handlers and sends the updated status to the Kubernetes API
func (a *nodePoolStatusHandler) sync(key string, obj *v3.NodePool) (*v3.NodePool, error) {
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

type nodePoolGeneratingHandler struct {
	NodePoolGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
	seen  sync.Map
}

// Remove handles the observed deletion of a resource, cascade deleting every associated resource previously applied
func (a *nodePoolGeneratingHandler) Remove(key string, obj *v3.NodePool) (*v3.NodePool, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.NodePool{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	if a.opts.UniqueApplyForResourceVersion {
		a.seen.Delete(key)
	}

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

// Handle executes the configured NodePoolGeneratingHandler and pass the resulting objects to apply.Apply, finally returning the new status of the resource
func (a *nodePoolGeneratingHandler) Handle(obj *v3.NodePool, status v3.NodePoolStatus) (v3.NodePoolStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.NodePoolGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}
	if !a.isNewResourceVersion(obj) {
		return newStatus, nil
	}

	err = generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
	if err != nil {
		return newStatus, err
	}
	a.storeResourceVersion(obj)
	return newStatus, nil
}

// isNewResourceVersion detects if a specific resource version was already successfully processed.
// Only used if UniqueApplyForResourceVersion is set in generic.GeneratingHandlerOptions
func (a *nodePoolGeneratingHandler) isNewResourceVersion(obj *v3.NodePool) bool {
	if !a.opts.UniqueApplyForResourceVersion {
		return true
	}

	// Apply once per resource version
	key := obj.Namespace + "/" + obj.Name
	previous, ok := a.seen.Load(key)
	return !ok || previous != obj.ResourceVersion
}

// storeResourceVersion keeps track of the latest resource version of an object for which Apply was executed
// Only used if UniqueApplyForResourceVersion is set in generic.GeneratingHandlerOptions
func (a *nodePoolGeneratingHandler) storeResourceVersion(obj *v3.NodePool) {
	if !a.opts.UniqueApplyForResourceVersion {
		return
	}

	key := obj.Namespace + "/" + obj.Name
	a.seen.Store(key, obj.ResourceVersion)
}

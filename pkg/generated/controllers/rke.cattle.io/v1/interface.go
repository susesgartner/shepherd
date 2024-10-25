/*
Copyright 2025 Rancher Labs, Inc.

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

package v1

import (
	"github.com/rancher/lasso/pkg/controller"
	v1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/rancher/shepherd/pkg/wrangler/pkg/generic"
	"github.com/rancher/wrangler/v3/pkg/schemes"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	schemes.Register(v1.AddToScheme)
}

type Interface interface {
	CustomMachine() CustomMachineController
	ETCDSnapshot() ETCDSnapshotController
	RKEBootstrap() RKEBootstrapController
	RKEBootstrapTemplate() RKEBootstrapTemplateController
	RKECluster() RKEClusterController
	RKEControlPlane() RKEControlPlaneController
}

func New(controllerFactory controller.SharedControllerFactory, ts *session.Session) Interface {
	return &version{
		controllerFactory: controllerFactory,
		ts:                ts,
	}
}

type version struct {
	controllerFactory controller.SharedControllerFactory
	ts                *session.Session
}

func (v *version) CustomMachine() CustomMachineController {
	return generic.NewController[*v1.CustomMachine, *v1.CustomMachineList](schema.GroupVersionKind{Group: "rke.cattle.io", Version: "v1", Kind: "CustomMachine"}, "custommachines", true, v.controllerFactory, v.ts)
}

func (v *version) ETCDSnapshot() ETCDSnapshotController {
	return generic.NewController[*v1.ETCDSnapshot, *v1.ETCDSnapshotList](schema.GroupVersionKind{Group: "rke.cattle.io", Version: "v1", Kind: "ETCDSnapshot"}, "etcdsnapshots", true, v.controllerFactory, v.ts)
}

func (v *version) RKEBootstrap() RKEBootstrapController {
	return generic.NewController[*v1.RKEBootstrap, *v1.RKEBootstrapList](schema.GroupVersionKind{Group: "rke.cattle.io", Version: "v1", Kind: "RKEBootstrap"}, "rkebootstraps", true, v.controllerFactory, v.ts)
}

func (v *version) RKEBootstrapTemplate() RKEBootstrapTemplateController {
	return generic.NewController[*v1.RKEBootstrapTemplate, *v1.RKEBootstrapTemplateList](schema.GroupVersionKind{Group: "rke.cattle.io", Version: "v1", Kind: "RKEBootstrapTemplate"}, "rkebootstraptemplates", true, v.controllerFactory, v.ts)
}

func (v *version) RKECluster() RKEClusterController {
	return generic.NewController[*v1.RKECluster, *v1.RKEClusterList](schema.GroupVersionKind{Group: "rke.cattle.io", Version: "v1", Kind: "RKECluster"}, "rkeclusters", true, v.controllerFactory, v.ts)
}

func (v *version) RKEControlPlane() RKEControlPlaneController {
	return generic.NewController[*v1.RKEControlPlane, *v1.RKEControlPlaneList](schema.GroupVersionKind{Group: "rke.cattle.io", Version: "v1", Kind: "RKEControlPlane"}, "rkecontrolplanes", true, v.controllerFactory, v.ts)
}

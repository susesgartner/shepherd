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

package v3

import (
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/shepherd/pkg/wrangler/pkg/generic"
)

// ClusterProxyConfigController interface for managing ClusterProxyConfig resources.
type ClusterProxyConfigController interface {
	generic.ControllerInterface[*v3.ClusterProxyConfig, *v3.ClusterProxyConfigList]
}

// ClusterProxyConfigClient interface for managing ClusterProxyConfig resources in Kubernetes.
type ClusterProxyConfigClient interface {
	generic.ClientInterface[*v3.ClusterProxyConfig, *v3.ClusterProxyConfigList]
}

// ClusterProxyConfigCache interface for retrieving ClusterProxyConfig resources in memory.
type ClusterProxyConfigCache interface {
	generic.CacheInterface[*v3.ClusterProxyConfig]
}

package machinepools

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	apisV1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/rancher/shepherd/clients/rancher"
	v1 "github.com/rancher/shepherd/clients/rancher/v1"
	"github.com/rancher/shepherd/extensions/defaults/states"
	"github.com/rancher/shepherd/extensions/defaults/stevetypes"
	"github.com/rancher/shepherd/extensions/defaults/timeouts"
	nodestat "github.com/rancher/shepherd/extensions/nodes"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kwait "k8s.io/apimachinery/pkg/util/wait"
)

const (
	pool = "pool"

	nodeRoleListLength = 4
)

// MatchNodeRolesToMachinePool matches the role of machinePools to the nodeRoles.
func MatchNodeRolesToMachinePool(nodeRoles NodeRoles, machinePools []apisV1.RKEMachinePool) (int, int32) {
	count := int32(0)
	for index, machinePoolConfig := range machinePools {
		if nodeRoles.ControlPlane != machinePoolConfig.ControlPlaneRole {
			continue
		}
		if nodeRoles.Etcd != machinePoolConfig.EtcdRole {
			continue
		}
		if nodeRoles.Worker != machinePoolConfig.WorkerRole {
			continue
		}

		count += *machinePoolConfig.Quantity

		return index, count
	}

	return -1, count
}

// updateMachinePoolQuantity is a helper method that will update the desired machine pool with the latest quantity.
func updateMachinePoolQuantity(client *rancher.Client, cluster *v1.SteveAPIObject, nodeRoles NodeRoles) (*v1.SteveAPIObject, error) {
	updateCluster, err := client.Steve.SteveType(stevetypes.Provisioning).ByID(cluster.ID)
	if err != nil {
		return nil, err
	}

	updatedCluster := new(apisV1.Cluster)
	err = v1.ConvertToK8sType(cluster, &updatedCluster)
	if err != nil {
		return nil, err
	}

	updatedCluster.ObjectMeta.ResourceVersion = updateCluster.ObjectMeta.ResourceVersion
	machineConfig, newQuantity := MatchNodeRolesToMachinePool(nodeRoles, updatedCluster.Spec.RKEConfig.MachinePools)

	newQuantity += nodeRoles.Quantity
	updatedCluster.Spec.RKEConfig.MachinePools[machineConfig].Quantity = &newQuantity

	logrus.Infof("Scaling the machine pool to %v total nodes", newQuantity)
	cluster, err = client.Steve.SteveType(stevetypes.Provisioning).Update(cluster, updatedCluster)
	if err != nil {
		return nil, err
	}

	err = kwait.PollUntilContextTimeout(context.TODO(), 500*time.Millisecond, timeouts.ThirtyMinute, true, func(ctx context.Context) (done bool, err error) {
		client, err = client.ReLogin()
		if err != nil {
			return false, err
		}

		clusterResp, err := client.Steve.SteveType(stevetypes.Provisioning).ByID(cluster.ID)
		if err != nil {
			return false, err
		}

		if clusterResp.ObjectMeta.State.Name == states.Active &&
			nodestat.AllMachineReady(client, cluster.ID, timeouts.ThirtyMinute) == nil {
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

// MatchMachineConfigToRolesIndex will return the index of the matching role for a given machineConfig.
func MatchMachineConfigToRolesIndex(machineConfig *MachinePoolConfig, objectRoles []Roles) int {
	for index, roles := range objectRoles {
		etcdMatch := false
		controlplaneMatch := false
		workerMatch := false
		windowsMatch := false

		for _, role := range roles.Roles {

			if machineConfig.Etcd && role == "etcd" {
				etcdMatch = true
			}

			if machineConfig.ControlPlane && role == "controlplane" {
				controlplaneMatch = true
			}

			if machineConfig.Worker && role == "worker" {
				workerMatch = true
			}

			if machineConfig.Windows && role == "windows" {
				windowsMatch = true
			}
		}

		if etcdMatch == machineConfig.Etcd &&
			controlplaneMatch == machineConfig.ControlPlane &&
			workerMatch == machineConfig.Worker &&
			windowsMatch == machineConfig.Windows {
			return index
		}
	}

	return -1
}

// NewRKEMachinePool is a constructor that sets up a apisV1.RKEMachinePool object to be used to provision a cluster.
func NewRKEMachinePool(machineObject v1.SteveAPIObject, pool Pools, machineConfig *MachinePoolConfig) apisV1.RKEMachinePool {
	machineConfigRef := &corev1.ObjectReference{
		Kind: machineObject.Kind,
		Name: machineObject.Name,
	}

	// windows pools are just worker pools exclusive to windows nodes.
	machineWorkerRole := machineConfig.Worker
	if machineConfig.Windows {
		machineWorkerRole = machineConfig.Windows
	}

	machinePoolQuantity := machineConfig.Quantity
	machinePool := apisV1.RKEMachinePool{
		ControlPlaneRole:     machineConfig.ControlPlane,
		EtcdRole:             machineConfig.Etcd,
		WorkerRole:           machineWorkerRole,
		NodeConfig:           machineConfigRef,
		Name:                 machineConfig.Name,
		Quantity:             &machinePoolQuantity,
		DrainBeforeDelete:    machineConfig.DrainBeforeDelete,
		NodeStartupTimeout:   machineConfig.NodeStartupTimeout,
		UnhealthyNodeTimeout: machineConfig.UnhealthyNodeTimeout,
		MaxUnhealthy:         machineConfig.MaxUnhealthy,
		UnhealthyRange:       machineConfig.UnhealthyRange,
		RKECommonNodeConfig: rkev1.RKECommonNodeConfig{
			Labels: pool.NodeLabels,
			Taints: pool.NodeTaints,
		},
	}

	if machineConfig.Windows {
		if machinePool.Labels != nil {
			machinePool.Labels["cattle.io/os"] = "windows"
		} else {
			machinePool.Labels = map[string]string{
				"cattle.io/os": "windows",
			}
		}
	}

	if machineConfig.HostnameLengthLimit > 0 {
		machinePool.HostnameLengthLimit = machineConfig.HostnameLengthLimit
	}

	return machinePool
}

type Pools struct {
	NodeLabels             map[string]string `json:"nodeLabels,omitempty" yaml:"nodeLabels,omitempty"`
	NodeTaints             []corev1.Taint    `json:"nodeTaints,omitempty" yaml:"nodeTaints,omitempty"`
	SpecifyCustomPrivateIP bool              `json:"specifyPrivateIP,omitempty" yaml:"specifyPrivateIP,omitempty"`
	SpecifyCustomPublicIP  bool              `json:"specifyPublicIP,omitempty" yaml:"specifyPublicIP,omitempty" default:"true"`
	CustomNodeNameSuffix   string            `json:"nodeNameSuffix,omitempty" yaml:"nodeNameSuffix,omitempty"`
}

type NodeRoles struct {
	ControlPlane bool  `json:"controlplane,omitempty" yaml:"controlplane,omitempty"`
	Etcd         bool  `json:"etcd,omitempty" yaml:"etcd,omitempty"`
	Worker       bool  `json:"worker,omitempty" yaml:"worker,omitempty"`
	Windows      bool  `json:"windows,omitempty" yaml:"windows,omitempty"`
	Quantity     int32 `json:"quantity" yaml:"quantity"`
}

type MachinePoolConfig struct {
	NodeRoles
	Name                 string           `json:"name,omitempty" yaml:"name,omitempty"`
	DrainBeforeDelete    bool             `json:"drainBeforeDelete,omitempty" yaml:"drainBeforeDelete,omitempty"`
	HostnameLengthLimit  int              `json:"hostnameLengthLimit" yaml:"hostnameLengthLimit" default:"0"`
	NodeStartupTimeout   *metav1.Duration `json:"nodeStartupTimeout,omitempty" yaml:"nodeStartupTimeout,omitempty"`
	UnhealthyNodeTimeout *metav1.Duration `json:"unhealthyNodeTimeout,omitempty" yaml:"unhealthyNodeTimeout,omitempty"`
	MaxUnhealthy         *string          `json:"maxUnhealthy,omitempty" yaml:"maxUnhealthy,omitempty"`
	UnhealthyRange       *string          `json:"unhealthyRange,omitempty" yaml:"unhealthyRange,omitempty"`
}

type Roles struct {
	Roles []string `json:"roles,omitempty" yaml:"roles,omitempty"`
}

// HostnameTruncation is a struct that is used to set the hostname length limit for a cluster or pool
type HostnameTruncation struct {
	PoolNameLengthLimit    int
	ClusterNameLengthLimit int
	Name                   string
}

func (n NodeRoles) String() string {
	result := make([]string, 0, nodeRoleListLength)
	if n.Quantity < 1 {
		return ""
	}
	if n.ControlPlane {
		result = append(result, "controlplane")
	}
	if n.Etcd {
		result = append(result, "etcd")
	}
	if n.Worker {
		result = append(result, "worker")
	}
	if n.Windows {
		result = append(result, "windows")
	}

	return fmt.Sprintf("%d %s", n.Quantity, strings.Join(result, "+"))
}

// CreateAllMachinePools will setup multiple node pools from a given config.
func CreateAllMachinePools(machineConfigs []MachinePoolConfig, pools []Pools, machineObjects []v1.SteveAPIObject, objectRoles []Roles, hostnameLengthLimits []HostnameTruncation) []apisV1.RKEMachinePool {
	machinePools := make([]apisV1.RKEMachinePool, 0, len(machineConfigs))

	for index, machineConfig := range machineConfigs {
		machineConfig.Name = pool + strconv.Itoa(index)
		if hostnameLengthLimits != nil && len(hostnameLengthLimits) >= index {
			machineConfig.HostnameLengthLimit = hostnameLengthLimits[index].PoolNameLengthLimit
			machineConfig.Name = hostnameLengthLimits[index].Name
		}
		objectIndex := MatchMachineConfigToRolesIndex(&machineConfig, objectRoles)
		machinePool := NewRKEMachinePool(machineObjects[objectIndex], pools[objectIndex], &machineConfig)
		machinePools = append(machinePools, machinePool)
	}

	return machinePools
}

// ScaleMachinePoolNodes is a helper method that will scale the machine pool to the desired quantity.
func ScaleMachinePoolNodes(client *rancher.Client, cluster *v1.SteveAPIObject, nodeRoles NodeRoles) (*v1.SteveAPIObject, error) {
	scaledClusterResp, err := updateMachinePoolQuantity(client, cluster, nodeRoles)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Machine pool has been scaled!")

	return scaledClusterResp, nil
}

// MatchRoleToPool matches the role of a pool to the Roles of a machine. Returns the index of the matching Roles.
func MatchRoleToPool(poolRole string, allRoles []Roles) int {
	hasMatch := false
	for poolIndex, machineRole := range allRoles {
		for _, configRole := range machineRole.Roles {
			if strings.Contains(poolRole, configRole) {
				hasMatch = true
			}
		}
		if hasMatch {
			return poolIndex
		}
	}
	logrus.Warn("unable to match pool to role, likely missing [roles] in machineConfig")
	return -1
}

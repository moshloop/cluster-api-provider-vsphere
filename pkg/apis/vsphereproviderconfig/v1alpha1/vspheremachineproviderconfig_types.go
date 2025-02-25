/*
Copyright 2018 The Kubernetes Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeadmv1beta1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VsphereMachineProviderConfig is the Schema for the vspheremachineproviderconfigs API
// +k8s:openapi-gen=true
type VsphereMachineProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	MachineRef        string             `json:"machineRef,omitempty"`
	MachineSpec       VsphereMachineSpec `json:"machineSpec,omitempty"`

	// KubeadmConfiguration holds the kubeadm configuration options
	// +optional
	KubeadmConfiguration KubeadmConfiguration `json:"kubeadmConfiguration,omitempty"`
}

// KubeadmConfiguration holds the various configurations that kubeadm uses
type KubeadmConfiguration struct {
	// JoinConfiguration is used to customize any kubeadm join configuration
	// parameters.
	Join kubeadmv1beta1.JoinConfiguration `json:"join,omitempty"`

	// InitConfiguration is used to customize any kubeadm init configuration
	// parameters.
	Init kubeadmv1beta1.InitConfiguration `json:"init,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VsphereMachineProviderConfigList contains a list of VsphereMachineProviderConfig
type VsphereMachineProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VsphereMachineProviderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VsphereMachineProviderConfig{}, &VsphereMachineProviderConfigList{})
}

//**** New extensions

type VsphereMachineSpec struct {
	Datacenter       string        `json:"datacenter"`
	Datastore        string        `json:"datastore"`
	ResourcePool     string        `json:"resourcePool,omitempty"`
	VMFolder         string        `json:"vmFolder,omitempty"`
	Networks         []NetworkSpec `json:"networks"`
	NumCPUs          int32         `json:"numCPUs,omitempty"`
	MemoryMB         int64         `json:"memoryMB,omitempty"`
	VMTemplate       string        `json:"template" yaml:"template"`
	Disks            []DiskSpec    `json:"disks"`
	Preloaded        bool          `json:"preloaded,omitempty"`
	VsphereCloudInit bool          `json:"vsphereCloudInit,omitempty"`
	TrustedCerts     []string      `json:"trustedCerts,omitempty"`
	NTPServers       []string      `json:"ntpServers,omitempty"`
}

type NetworkSpec struct {
	NetworkName string   `json:"networkName"`
	IPConfig    IPConfig `json:"ipConfig,omitempty"`
}

type IPConfig struct {
	NetworkType NetworkType `json:"networkType"`
	IP          string      `json:"ip,omitempty"`
	Netmask     string      `json:"netmask,omitempty"`
	Gateway     string      `json:"gateway,omitempty"`
	Dns         []string    `json:"dns,omitempty"`
}

type NetworkType string

const (
	Static NetworkType = "static"
	DHCP   NetworkType = "dhcp"
)

type DiskSpec struct {
	DiskSizeGB int64  `json:"diskSizeGB,omitempty"`
	DiskLabel  string `json:"diskLabel,omitempty"`
}

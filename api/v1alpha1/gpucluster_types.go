/*
Copyright 2024.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type State string

const (
	// Ignored indicates duplicate gpucluster instances and rest are ignored.
	Ignored State = "ignored"

	// Ready indicates all components of GPUCluster are ready
	Ready State = "ready"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GPUClusterSpec defines the desired state of GPUCluster
type GPUClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// DevicePlugin component spec
	DevicePlugin DevicePluginSpec `json:"devicePlugin"`
}

// GPUClusterStatus defines the observed state of GPUCluster
type GPUClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Namespace string `json:"namespace,omitempty"`

	// status of gpucluster
	State State `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GPUCluster is the Schema for the gpuclusters API
type GPUCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GPUClusterSpec   `json:"spec,omitempty"`
	Status GPUClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GPUClusterList contains a list of GPUCluster
type GPUClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GPUCluster `json:"items"`
}

type DevicePluginSpec struct {
	// Enabled indicates whether to deploy xdxct-device-plugin
	Enabled *bool `json:"enabled,omitempty"`

	// Xdxct Device-plugin repository
	Repository string `json:"repository,omitempty"`

	// Xdxct Device-plugin image
	Image string `json:"image,omitempty"`

	// Xdxct Device-plugin image tag
	Version string `json:"version,omitempty"`

	// Device-plugin image pull policy
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`

	// Device-plugin image pull secrets
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`

	// Optional: List of arguments
	Args []string `json:"args,omitempty"`

	// Optional: List of environmemt variables
	Env []EnvVar `json:"env,omitempty"`

	// Optional: resources requests and limits for device plugin pod
	Resources *ResourceRequirements `json:"resources,omitempty"`

	// Optional: Configmap for Device-plugin
	Config *DevicePluginConfig `json:"config,omitempty"`
}

type EnvVar struct {
	// Environment name
	Name string `json:"name"`

	// Environment value
	Value string `json:"value"`
}

type ResourceRequirements struct {
	// Limits describes the maxium amount of compute resource requirements
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Limits corev1.ResourceList `json:"limits,omitempty"`

	// Requests describes the minimum amount of compute resources requirements
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Requests corev1.ResourceList `json:"requests,omitempty"`
}

// DevicePluginSpec properties for device plugin deployment
type DevicePluginConfig struct {
	// ConfigMap name
	Name string `json:"name"`

	// Default config for ConfigMap
	Default string `json:"default"`
}

func init() {
	SchemeBuilder.Register(&GPUCluster{}, &GPUClusterList{})
}

func (c *GPUCluster) SetStatus(s State, ns string) {
	c.Status.State = s
	c.Status.Namespace = ns
}

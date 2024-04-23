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
	"fmt"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type State string
type Runtime string

const (
	// Ignored indicates duplicate gpucluster instances and rest are ignored.
	Ignored State = "ignored"

	// Ready indicates all components of GPUCluster are ready
	Ready State = "ready"

	// notReady indicates the components of GPUCluster are notReady
	NotReady State = "notReady"

	// the components is disabled
	Disabled State = "disabled"
)

const (
	// Docker runtime
	Docker Runtime = "docker"
	// CRIO runtime
	CRIO Runtime = "crio"
	// Containerd runtime
	Containerd Runtime = "containerd"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GPUClusterSpec defines the desired state of GPUCluster
type GPUClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Operator defines configurations for cluster
	Operator OperatorSpec `json:"operator"`

	// Daemonset defines common configuration for all components
	DaemonSets DaemonSetsSpec `json:"daemonSets"`

	// DevicePlugin component spec
	DevicePlugin DevicePluginSpec `json:"devicePlugin"`

	// Kubevirt device plugin component spec
	KubevirtDevicePlugin KubevirtDevicePluginSpec `json:"kubevirtDevicePlugin,omitempty"`

	// VGPUDeviceManager component spec
	VGPUDeviceManager VGPUDeviceManagerSpec `json:"vgpuDeviceManager,omitempty"`

	// VFIOManager for configuration to deploy vfio-pci manager
	VFIOManager VFIOManagerSpec `json:"vfioManager,omitempty"`
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

// OperatorSpec describes configuration options for the operator
type OperatorSpec struct {
	RuntimeClass string `json:"runtimeClass,omitempty"`
}

// DaemonSetsSpec describe configuration for all daemonsets components
type DaemonSetsSpec struct {
	Labels map[string]string `json:"labels,omitempty"`

	Annotations map[string]string `json:"annotations,omitempty"`

	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	UpdateStrategy string `json:"updateStrategy,omitempty"`

	RollingUpdate *RollingUpdateSpec `json:"rollingUpdate,omitempty"`

	PriorityClassName string `json:"priorityClassName,omitempty"`
}

// RollingUpdateSpec indicates configurations for all daemonset pod
type RollingUpdateSpec struct {
	MaxUnavilable string `json:"maxUnavilable,omitempty"`
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

type KubevirtDevicePluginSpec struct {
	// Enabled indicates whether to deploy kubevirt-device-plugin
	Enabled *bool `json:"enabled,omitempty"`

	// Xdxct kubevirt-device-plugin image repository
	Repository string `json:"repository,omitempty"`

	// Xdxct kubevirt-device-plugin image name
	Image string `json:"image,omitempty"`

	// Xdxct kubevirt-device-plugin image tag
	Version string `json:"version,omitempty"`

	// Xdxct kubevirt-device-plugin image Pull Policy
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`

	// Xdxct kubevirt-device-plugin image Pull Secrets
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`

	// Optional: List of arguments
	Args []string `json:"args,omitempty"`

	// Optional: List of environmemt variables
	Env []EnvVar `json:"env,omitempty"`

	// Optional: resources requests and limits for xdxct kubevirt-device-plugin pod
	Resources *ResourceRequirements `json:"resources,omitempty"`
}

type VGPUDeviceManagerSpec struct {
	// Enabled indicates whether to deploy vgpu-device-manager
	Enabled *bool `json:"enabled,omitempty"`

	// Xdxct vgpu-device-manager image repository
	Repository string `json:"repository,omitempty"`

	// Xdxct vgpu-device-manager image name
	Image string `json:"image,omitempty"`

	// Xdxct vgpu-device-manager image tag
	Version string `json:"version,omitempty"`

	// Xdxct vgpu-device-manager image Pull Policy
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`

	// Xdxct vgpu-device-manager image Pull Secrets
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`

	// Optional: List of arguments
	Args []string `json:"args,omitempty"`

	// Optional: List of environmemt variables
	Env []EnvVar `json:"env,omitempty"`

	// Optional: resources requests and limits for xdxct kubevirt-device-plugin pod
	Resources *ResourceRequirements `json:"resources,omitempty"`

	// Xdxct vgpu-device-manager configuration for vGPU Device type
	Config *VGPUDeviceManagerConfigSpec `json:"config,omitempty"`
}

type VGPUDeviceManagerConfigSpec struct {
	// the name of configmap for vgpu-device-config
	Name string `json:"name,omitempty"`

	// config for vgpu devices
	Default string `json:"default,omitempty"`
}

type VFIOManagerSpec struct {
	// Enabled indicates whether to deploy vfio-manager
	Enabled *bool `json:"enabled,omitempty"`

	// Xdxct vfio-manager image repository
	Repository string `json:"repository,omitempty"`

	// Xdxct vfio-manager image name
	Image string `json:"image,omitempty"`

	// Xdxct vfio-manager image tag
	Version string `json:"version,omitempty"`

	// Xdxct vfio-manager image Pull Policy
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`

	// Xdxct vfio-manager image Pull Secrets
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`

	// Optional: List of arguments
	Args []string `json:"args,omitempty"`

	// Optional: List of environmemt variables
	Env []EnvVar `json:"env,omitempty"`

	// Optional: resources requests and limits for xdxct-vfio-manager pod
	Resources *ResourceRequirements `json:"resources,omitempty"`
}

func init() {
	SchemeBuilder.Register(&GPUCluster{}, &GPUClusterList{})
}

func imagePath(repoistory string, image string, version string, imagePathEnvName string) (string, error) {
	// 1. GpuClusterSpec
	var crdImagePath string
	if repoistory == "" && version == "" {
		if image != "" {
			crdImagePath = image
		}
	} else {
		if strings.HasPrefix(version, "sha256:") {
			crdImagePath = repoistory + "/" + image + "@" + version
		} else {
			crdImagePath = repoistory + "/" + image + ":" + version
		}
	}
	if crdImagePath != "" {
		return crdImagePath, nil
	}

	// 2. Env passed to pod
	envImagePath := os.Getenv(imagePathEnvName)
	if envImagePath != "" {
		return envImagePath, nil
	}
	return "", fmt.Errorf("empty image path both GpuClusterSpec or Env: %s", imagePathEnvName)
}

func ImagePath(spec interface{}) (string, error) {
	switch v := spec.(type) {
	case *DevicePluginSpec:
		config := spec.(*DevicePluginSpec)
		return imagePath(config.Repository, config.Image, config.Version, "DEVICE_PLUGIN_IMAGE")
	case *KubevirtDevicePluginSpec:
		config := spec.(*KubevirtDevicePluginSpec)
		return imagePath(config.Repository, config.Image, config.Version, "KUBEVIRT_DEVICE_PLUGIN_IMAGE")
	case *VGPUDeviceManagerSpec:
		config := spec.(*VGPUDeviceManagerSpec)
		return imagePath(config.Repository, config.Image, config.Version, "VGPU_DEVICE_MANAGER")
	default:
		return "", fmt.Errorf("invalid type to construct image type path: %v", v)
	}
}

func ImagePullPolicy(pullPolicy string) corev1.PullPolicy {
	var imagePullPolicy corev1.PullPolicy
	switch pullPolicy {
	case "Always":
		imagePullPolicy = corev1.PullAlways
	case "Never":
		imagePullPolicy = corev1.PullNever
	case "IfNotPresent":
		imagePullPolicy = corev1.PullIfNotPresent
	default:
		imagePullPolicy = corev1.PullIfNotPresent
	}
	return imagePullPolicy
}

func (c *GPUCluster) SetStatus(s State, ns string) {
	c.Status.State = s
	c.Status.Namespace = ns
}

func (d *DevicePluginSpec) IsEnabled() bool {
	if d.Enabled == nil {
		return true
	}
	return *d.Enabled
}

func (k *KubevirtDevicePluginSpec) IsEnabled() bool {
	if k.Enabled == nil {
		return true
	}
	return *k.Enabled
}

func (v *VGPUDeviceManagerSpec) IsEnabled() bool {
	if v.Enabled == nil {
		return true
	}
	return *v.Enabled
}

func (vm *VFIOManagerSpec) IsEnabled() bool {
	if vm.Enabled == nil {
		return true
	}
	return *vm.Enabled
}

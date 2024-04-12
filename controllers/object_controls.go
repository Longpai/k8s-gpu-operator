package controllers

import (
	gpuv1alpha1 "github.com/chen-mao/k8s-gpu-operator.git/api/v1alpha1"
)

type controlFunc []func(c GPUClusterController) (gpuv1alpha1.State, error)

func ServiceAccount(c GPUClusterController) (gpuv1alpha1.State, error) {
	return "", nil
}

func Role(c GPUClusterController) (gpuv1alpha1.State, error) {
	return "", nil
}

func ClusterRole(c GPUClusterController) (gpuv1alpha1.State, error) {
	return "", nil
}

func RoleBinding(c GPUClusterController) (gpuv1alpha1.State, error) {
	return "", nil
}

func ClusterRoleBinding(c GPUClusterController) (gpuv1alpha1.State, error) {
	return "", nil
}

func ConfigMaps(c GPUClusterController) (gpuv1alpha1.State, error) {
	return "", nil
}

func DaemonSet(c GPUClusterController) (gpuv1alpha1.State, error) {
	return "", nil
}

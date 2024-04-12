package controllers

import (
	"fmt"

	gpuv1alpha1 "github.com/chen-mao/k8s-gpu-operator.git/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type controlFunc []func(c GPUClusterController) (gpuv1alpha1.State, error)

// create ServiceAccount resource
func ServiceAccount(c GPUClusterController) (gpuv1alpha1.State, error) {
	saObj := c.resources[0].ServiceAccount.DeepCopy()
	saObj.Namespace = c.namespace

	fmt.Println("ServiceAccount", saObj.Name)

	// 组件被disabled时，清理掉已经存在资源
	if !c.isStateEnabled(c.componentNames[0]) {
		err := c.client.Delete(c.ctx, saObj)
		if err != nil && !apierrors.IsNotFound(err) {
			fmt.Printf("Failed to delete serviceAccount: %v", err)
			return gpuv1alpha1.NotReady, err
		}
		return gpuv1alpha1.Disabled, err
	}

	// 将该资源与控制器相关联
	if err := controllerutil.SetControllerReference(c.singleton, saObj, c.schema); err != nil {
		fmt.Printf("error: %v", err)
		return gpuv1alpha1.NotReady, err
	}

	if err := c.client.Create(c.ctx, saObj); err != nil {
		if apierrors.IsAlreadyExists(err) {
			fmt.Println("Found Resource, Skip update")
			return gpuv1alpha1.Ready, nil
		}

		fmt.Printf("Failed to create: %v", err)
		return gpuv1alpha1.NotReady, err
	}
	return gpuv1alpha1.Ready, nil
}

func Role(c GPUClusterController) (gpuv1alpha1.State, error) {
	return gpuv1alpha1.Ready, nil
}

func ClusterRole(c GPUClusterController) (gpuv1alpha1.State, error) {
	return gpuv1alpha1.Ready, nil
}

func RoleBinding(c GPUClusterController) (gpuv1alpha1.State, error) {
	return gpuv1alpha1.Ready, nil
}

func ClusterRoleBinding(c GPUClusterController) (gpuv1alpha1.State, error) {
	return gpuv1alpha1.Ready, nil
}

func ConfigMaps(c GPUClusterController) (gpuv1alpha1.State, error) {
	return gpuv1alpha1.Ready, nil
}

func DaemonSet(c GPUClusterController) (gpuv1alpha1.State, error) {
	return gpuv1alpha1.Ready, nil
}

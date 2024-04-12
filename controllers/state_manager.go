package controllers

import (
	"context"
	"fmt"

	gpuv1alpha1 "github.com/chen-mao/k8s-gpu-operator.git/api/v1alpha1"
)

// controlFunc: 保存了组件的执行函数
// controls: 保存各个组件
type GPUClusterController struct {
	singleton *gpuv1alpha1.GPUCluster
	resources []Resouces
	controls  []controlFunc
	namespace string
}

func addState(c *GPUClusterController, path string) {
	res, ctrlFunc := addRescourcesControls(path)
	c.resources = append(c.resources, res)
	c.controls = append(c.controls, ctrlFunc)
	// fmt.Println(c.resources)
	fmt.Println(c.controls)
}

func (c *GPUClusterController) init(ctx context.Context, gpuCluster *gpuv1alpha1.GPUCluster) error {
	c.singleton = gpuCluster

	addState(c, "/opt/k8s-gpu-operator/device-plugin")
	return nil
}

func (c *GPUClusterController) step() (gpuv1alpha1.State, error) {
	result := gpuv1alpha1.Ready

	for _, fs := range c.controls[0] {
		stat, err := fs(*c)
		if err != nil {
			return stat, err
		}
		// 成功部署了资源，检查ready.
		// 只要组件中有一个资源没有ready，则该组件就是安装失败，不会继续安装后面的组件
		if stat != gpuv1alpha1.Ready {
			result = stat
		}
	}
	// c.idx++
	return result, nil
}

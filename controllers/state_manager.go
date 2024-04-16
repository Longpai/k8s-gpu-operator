package controllers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	gpuv1alpha1 "github.com/chen-mao/k8s-gpu-operator.git/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// controlFunc: 保存了组件的执行函数
// controls: 保存各个组件
type GPUClusterController struct {
	client client.Client

	ctx       context.Context
	singleton *gpuv1alpha1.GPUCluster
	schema    *runtime.Scheme

	resources      []Resouces
	controls       []controlFunc
	componentNames []string
	namespace      string
	index          int

	runtime gpuv1alpha1.Runtime
}

func addState(c *GPUClusterController, path string) {
	res, ctrlFunc := addRescourcesControls(path)
	c.resources = append(c.resources, res)
	c.controls = append(c.controls, ctrlFunc)
	c.componentNames = append(c.componentNames, filepath.Base(path))

	fmt.Println(c.componentNames)
	fmt.Println(c.controls)
}

func (c *GPUClusterController) init(ctx context.Context, reconciler *GPUClusterReconciler, gpuCluster *gpuv1alpha1.GPUCluster) error {
	c.ctx = ctx
	c.client = reconciler.Client
	c.schema = reconciler.Scheme
	c.singleton = gpuCluster
	c.index = 0
	if len(c.controls) == 0 {
		gpuClusterCtrl.namespace = os.Getenv("OPERATOR_NAMESPACE")
		if gpuClusterCtrl.namespace == "" {
			// 任何操作都是在namespace下，如果没有namespace, 则退出循环
			fmt.Println("namespace environment variable not set, exit.")
			os.Exit(1)
		}

		addState(c, "/opt/k8s-gpu-operator/device-plugin")
		addState(c, "/opt/k8s-gpu-operator/kubevirt-device-plugin")
	}

	return nil
}

func (c *GPUClusterController) step() (gpuv1alpha1.State, error) {
	result := gpuv1alpha1.Ready
	for _, fs := range c.controls[c.index] {
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
	// install the next component
	c.index++
	return result, nil
}

func (c *GPUClusterController) isStateEnabled(name string) bool {
	GPUClusterSpec := &c.singleton.Spec
	switch name {
	case "device-plugin":
		return GPUClusterSpec.DevicePlugin.IsEnabled()
	case "kubevirt-device-plugin":
		return GPUClusterSpec.KubevirtDevicePlugin.IsEnabled()
	default:
		fmt.Println("invalid component name")
		return false
	}
}

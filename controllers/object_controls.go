package controllers

import (
	"context"
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"strings"

	gpuv1alpha1 "github.com/chen-mao/k8s-gpu-operator.git/api/v1alpha1"
	"github.com/davecgh/go-spew/spew"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	// DefaultRuntimeClass represents "nvidia" RuntimeClass
	DefaultRuntimeClass = "xdxct"
	// XdxctAnnotationHashKey indicates annotation name for last applied hash by gpu-operator
	XdxctAnnotationHashKey = "xdxct.com/hasher"
	// PodControllerRevisionHashLabelKey is the annotation key for pod controller revision hash value
	PodControllerRevisionHashLabelKey = "controller-revision-hash"
)

type controlFunc []func(c GPUClusterController) (gpuv1alpha1.State, error)

// create ServiceAccount resource
func ServiceAccount(c GPUClusterController) (gpuv1alpha1.State, error) {
	index := c.index
	saObj := c.resources[index].ServiceAccount.DeepCopy()
	saObj.Namespace = c.namespace

	fmt.Println("ServiceAccount:", saObj.Name)

	// 组件被disabled时，清理掉已经存在资源
	if !c.isStateEnabled(c.componentNames[index]) {
		err := c.client.Delete(c.ctx, saObj)
		if err != nil && !apierrors.IsNotFound(err) {
			fmt.Printf("Failed to delete serviceAccount: %v", err)
			return gpuv1alpha1.NotReady, err
		}
		// fmt.Println("ServiceAccount:", gpuv1alpha1.Disabled)
		return gpuv1alpha1.Disabled, nil
	}

	// 将资源与控制器相关联
	if err := controllerutil.SetControllerReference(c.singleton, saObj, c.schema); err != nil {
		return gpuv1alpha1.NotReady, err
	}

	if err := c.client.Create(c.ctx, saObj); err != nil {
		if apierrors.IsAlreadyExists(err) {
			fmt.Println("Found Resource, Skip update")
			return gpuv1alpha1.Ready, nil
		}

		fmt.Printf("Failed to Create: %v", err)
		return gpuv1alpha1.NotReady, err
	}
	return gpuv1alpha1.Ready, nil
}

// create Role resource
func Role(c GPUClusterController) (gpuv1alpha1.State, error) {
	roleObj := c.resources[0].Role.DeepCopy()
	roleObj.Namespace = c.namespace

	fmt.Println("Role:", roleObj.Name)

	// 组件被disabled时，清理掉已经存在资源
	if !c.isStateEnabled(c.componentNames[0]) {
		err := c.client.Delete(c.ctx, roleObj)
		if err != nil && !apierrors.IsNotFound(err) {
			fmt.Printf("Failed to delete role: %v", err)
			return gpuv1alpha1.NotReady, err
		}
		return gpuv1alpha1.Disabled, nil
	}

	// 将资源与控制器相关联
	if err := controllerutil.SetControllerReference(c.singleton, roleObj, c.schema); err != nil {
		return gpuv1alpha1.NotReady, err
	}

	if err := c.client.Create(c.ctx, roleObj); err != nil {
		if apierrors.IsAlreadyExists(err) {
			fmt.Println("Found Resource, Update")
			err = c.client.Update(c.ctx, roleObj)
			if err != nil {
				fmt.Printf("Failed to Update: %v", err)
			}
			return gpuv1alpha1.Ready, nil
		}

		fmt.Printf("Failed to create: %v", err)
		return gpuv1alpha1.NotReady, err
	}
	return gpuv1alpha1.Ready, nil
}

// create clusterRole resource
func ClusterRole(c GPUClusterController) (gpuv1alpha1.State, error) {
	clusterRoleObj := c.resources[0].ClusterRole.DeepCopy()
	clusterRoleObj.Namespace = c.namespace

	fmt.Println("clusterRole:", clusterRoleObj.Name)

	// 组件被disabled时，清理掉已经存在资源
	if !c.isStateEnabled(c.componentNames[0]) {
		err := c.client.Delete(c.ctx, clusterRoleObj)
		if err != nil && !apierrors.IsNotFound(err) {
			fmt.Printf("Failed to delete clusterRole: %v", err)
			return gpuv1alpha1.NotReady, err
		}
		return gpuv1alpha1.Disabled, nil
	}

	// 将资源与控制器相关联
	if err := controllerutil.SetControllerReference(c.singleton, clusterRoleObj, c.schema); err != nil {
		return gpuv1alpha1.NotReady, err
	}

	if err := c.client.Create(c.ctx, clusterRoleObj); err != nil {
		if apierrors.IsAlreadyExists(err) {
			fmt.Println("Found Resource, Update")
			err = c.client.Update(c.ctx, clusterRoleObj)
			if err != nil {
				fmt.Printf("Failed to Update: %v", err)
			}
			return gpuv1alpha1.Ready, nil
		}

		fmt.Printf("Failed to create: %v", err)
		return gpuv1alpha1.NotReady, err
	}
	return gpuv1alpha1.Ready, nil
}

// create RoleBinding resource
func RoleBinding(c GPUClusterController) (gpuv1alpha1.State, error) {
	RoleBindingObj := c.resources[0].RoleBinding.DeepCopy()
	RoleBindingObj.Namespace = c.namespace

	fmt.Println("RoleBinding:", RoleBindingObj.Name)

	// 组件被disabled时，清理掉已经存在资源
	if !c.isStateEnabled(c.componentNames[0]) {
		err := c.client.Delete(c.ctx, RoleBindingObj)
		if err != nil && !apierrors.IsNotFound(err) {
			fmt.Printf("Failed to delete RoleBinding: %v", err)
			return gpuv1alpha1.NotReady, err
		}
		return gpuv1alpha1.Disabled, nil
	}

	// 将资源与控制器相关联
	if err := controllerutil.SetControllerReference(c.singleton, RoleBindingObj, c.schema); err != nil {
		return gpuv1alpha1.NotReady, err
	}

	if err := c.client.Create(c.ctx, RoleBindingObj); err != nil {
		if apierrors.IsAlreadyExists(err) {
			fmt.Println("Found Resource, Update")
			err = c.client.Update(c.ctx, RoleBindingObj)
			if err != nil {
				fmt.Printf("Failed to Update: %v", err)
			}
			return gpuv1alpha1.Ready, nil
		}

		fmt.Printf("Failed to create: %v", err)
		return gpuv1alpha1.NotReady, err
	}
	return gpuv1alpha1.Ready, nil
}

// create ClusterRoleBinding resource
func ClusterRoleBinding(c GPUClusterController) (gpuv1alpha1.State, error) {
	clusterRoleBindingObj := c.resources[0].ClusterRoleBinding.DeepCopy()
	clusterRoleBindingObj.Namespace = c.namespace

	fmt.Println("clusterRoleBindingObj:", clusterRoleBindingObj.Name)

	// 组件被disabled时，清理掉已经存在资源
	if !c.isStateEnabled(c.componentNames[0]) {
		err := c.client.Delete(c.ctx, clusterRoleBindingObj)
		if err != nil && !apierrors.IsNotFound(err) {
			fmt.Printf("Failed to delete clusterRoleBinding: %v", err)
			return gpuv1alpha1.NotReady, err
		}
		return gpuv1alpha1.Disabled, nil
	}

	// 将资源与控制器相关联
	if err := controllerutil.SetControllerReference(c.singleton, clusterRoleBindingObj, c.schema); err != nil {
		return gpuv1alpha1.NotReady, err
	}

	if err := c.client.Create(c.ctx, clusterRoleBindingObj); err != nil {
		if apierrors.IsAlreadyExists(err) {
			fmt.Println("Found Resource, Update")
			err = c.client.Update(c.ctx, clusterRoleBindingObj)
			if err != nil {
				fmt.Printf("Failed to Update: %v", err)
			}
			return gpuv1alpha1.Ready, nil
		}

		fmt.Printf("Failed to create: %v", err)
		return gpuv1alpha1.NotReady, err
	}
	return gpuv1alpha1.Ready, nil
}

func ConfigMaps(c GPUClusterController) (gpuv1alpha1.State, error) {
	return gpuv1alpha1.Ready, nil
}

// create DaemonSet resource
func DaemonSet(c GPUClusterController) (gpuv1alpha1.State, error) {
	ctx := c.ctx
	daemonSetObj := c.resources[0].Daemonset.DeepCopy()
	daemonSetObj.Namespace = c.namespace

	fmt.Println("daemonSetObj:", daemonSetObj.Name)

	// 组件被disabled时，清理掉已经存在资源
	if !c.isStateEnabled(c.componentNames[0]) {
		err := c.client.Delete(c.ctx, daemonSetObj)
		if err != nil && !apierrors.IsNotFound(err) {
			fmt.Printf("Failed to delete daemonSet: %v", err)
			return gpuv1alpha1.NotReady, err
		}
		return gpuv1alpha1.Disabled, nil
	}

	err := preDeployDaemonSet(c, daemonSetObj)
	if err != nil {
		fmt.Println("failed to pre-config for daemonSet:", err)
		return gpuv1alpha1.NotReady, err
	}

	if err := controllerutil.SetControllerReference(c.singleton, daemonSetObj, c.schema); err != nil {
		fmt.Println("filed to SetControllerReference", err)
		return gpuv1alpha1.NotReady, err
	}

	if daemonSetObj.Labels == nil {
		daemonSetObj.Labels = make(map[string]string)
	}
	for key, value := range c.singleton.Spec.DaemonSets.Labels {
		daemonSetObj.Labels[key] = value
	}

	if daemonSetObj.Annotations == nil {
		daemonSetObj.Annotations = make(map[string]string)
	}
	for key, value := range c.singleton.Spec.DaemonSets.Annotations {
		daemonSetObj.Annotations[key] = value

	}

	foundDs := &appsv1.DaemonSet{}
	err = c.client.Get(ctx, types.NamespacedName{
		Namespace: foundDs.Namespace,
		Name:      foundDs.Name,
	}, foundDs)
	if err != nil && apierrors.IsNotFound(err) {
		fmt.Println("DaemonSet not found:", daemonSetObj.Name, "Start Creating ...")
		// generate hash for the spec.
		hashStr := getDaemonSetHash(daemonSetObj)
		// add annotations to daemonSet with hash
		daemonSetObj.Annotations[XdxctAnnotationHashKey] = hashStr
		err = c.client.Create(ctx, daemonSetObj)
		if err != nil {
			fmt.Printf("failed to create %s daemonSet: %v", daemonSetObj.Name, err)
			return gpuv1alpha1.NotReady, err
		}
		return checkDaemonSetReady(daemonSetObj.Name, c), nil
	} else if err != nil {
		fmt.Printf("failed to get %s daemonSet: %v", daemonSetObj.Name, err)
		return gpuv1alpha1.NotReady, err
	}

	// checkDaemonSetChanged(foundDs, daemonSetObj)
	// checkDaemonSetChanged()
	return checkDaemonSetReady(daemonSetObj.Name, c), nil
}

// pre-config for DaemonSet: fillful daemonset with configuration-info
func preDeployDaemonSet(c GPUClusterController, daemonSetObj *appsv1.DaemonSet) error {
	transformations := map[string]func(*appsv1.DaemonSet, *gpuv1alpha1.GPUClusterSpec, GPUClusterController) error{
		"xdxct-device-plugin-daemonset": TransformDevicePlugin,
	}
	fs, ok := transformations[daemonSetObj.Name]
	if !ok {
		fmt.Printf("No transformation for Daemonset '%s'", daemonSetObj.Name)
		return nil
	}
	// c.singleton.Spec: 用户自定义的config spec
	// daemonSetObj: services中的组件
	err := applyCommonDaemonsetConfig(daemonSetObj, &c.singleton.Spec)
	if err != nil {
		return fmt.Errorf("failed to apply common DaemonSet transformation: %s", daemonSetObj.Name)
	}

	err = fs(daemonSetObj, &c.singleton.Spec, c)
	if err != nil {
		return fmt.Errorf("failed to apply transformation: %s", daemonSetObj.Name)
	}

	applyCommonDaemonsetMetadata(daemonSetObj, &c.singleton.Spec.DaemonSets)
	return nil
}

// the Common Daemonset Config
func applyCommonDaemonsetConfig(daemonSetObj *appsv1.DaemonSet, config *gpuv1alpha1.GPUClusterSpec) error {
	switch config.DaemonSets.UpdateStrategy {
	case "OnDelete":
		daemonSetObj.Spec.UpdateStrategy = appsv1.DaemonSetUpdateStrategy{
			Type: appsv1.OnDeleteDaemonSetStrategyType,
		}
	case "RollingUpdate":
		fallthrough
	default:
		if config.DaemonSets.RollingUpdate == nil || config.DaemonSets.RollingUpdate.MaxUnavilable == "" {
			return nil
		}
		var intOrString intstr.IntOrString
		if strings.HasPrefix(config.DaemonSets.RollingUpdate.MaxUnavilable, "%") {
			intOrString = intstr.IntOrString{
				Type:   intstr.String,
				StrVal: config.DaemonSets.RollingUpdate.MaxUnavilable,
			}
		} else {
			int64Val, err := strconv.ParseInt(config.DaemonSets.RollingUpdate.MaxUnavilable, 10, 32)
			if err != nil {
				return fmt.Errorf("failed to apply rolling update config: %v", err)
			}
			intOrString = intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: int32(int64Val),
			}
		}
		rollingUpdateSpec := appsv1.RollingUpdateDaemonSet{
			MaxUnavailable: &intOrString,
		}
		daemonSetObj.Spec.UpdateStrategy = appsv1.DaemonSetUpdateStrategy{
			Type:          appsv1.RollingUpdateDaemonSetStrategyType,
			RollingUpdate: &rollingUpdateSpec,
		}
	}

	// set PriorityClass
	if config.DaemonSets.PriorityClassName != "" {
		daemonSetObj.Spec.Template.Spec.PriorityClassName = config.DaemonSets.PriorityClassName
	}

	if len(config.DaemonSets.Tolerations) > 0 {
		daemonSetObj.Spec.Template.Spec.Tolerations = config.DaemonSets.Tolerations
	}

	return nil
}

func TransformDevicePlugin(daemonSet *appsv1.DaemonSet, config *gpuv1alpha1.GPUClusterSpec, c GPUClusterController) error {
	// update image
	image, err := gpuv1alpha1.ImagePath(&config.DevicePlugin)
	if err != nil {
		return err
	}
	daemonSet.Spec.Template.Spec.Containers[0].Image = image

	// update image pull policy
	daemonSet.Spec.Template.Spec.Containers[0].ImagePullPolicy = gpuv1alpha1.ImagePullPolicy(config.DevicePlugin.ImagePullPolicy)

	// set image pull secrets
	if len(config.DevicePlugin.ImagePullSecrets) > 0 {
		for _, secret := range config.DevicePlugin.ImagePullSecrets {
			if !containSecret(daemonSet.Spec.Template.Spec.ImagePullSecrets, secret) {
				daemonSet.Spec.Template.Spec.ImagePullSecrets = append(daemonSet.Spec.Template.Spec.ImagePullSecrets, corev1.LocalObjectReference{
					Name: secret,
				})
			}
		}
	}

	// set arguments if specified for device-plugin container
	if len(config.DevicePlugin.Args) > 0 {
		daemonSet.Spec.Template.Spec.Containers[0].Args = config.DevicePlugin.Args
	}

	// set environments if specified for device-plugin container
	if len(config.DevicePlugin.Env) > 0 {
		for _, env := range config.DevicePlugin.Env {
			setContainerEnv(&daemonSet.Spec.Template.Spec.Containers[0], env.Name, env.Value)
		}
	}

	// set resource limits
	if config.DevicePlugin.Resources != nil {
		for i := range daemonSet.Spec.Template.Spec.Containers {
			daemonSet.Spec.Template.Spec.Containers[i].Resources.Requests = config.DevicePlugin.Resources.Requests
			daemonSet.Spec.Template.Spec.Containers[i].Resources.Limits = config.DevicePlugin.Resources.Limits
		}
	}

	// set RuntimeClass for supported runtimes
	setRuntimeClass(&daemonSet.Spec.Template.Spec, c.runtime, config.Operator.RuntimeClass)

	return nil
}

func containSecret(secrets []corev1.LocalObjectReference, secretName string) bool {
	for _, s := range secrets {
		if s.Name == secretName {
			return true
		}
	}
	return false
}

func setContainerEnv(c *corev1.Container, key, value string) {
	for i, v := range c.Env {
		if v.Name != key {
			continue
		}
		c.Env[i].Value = value
		return
	}
	c.Env = append(c.Env, corev1.EnvVar{
		Name:  key,
		Value: value,
	})
}

func setRuntimeClass(podSpec *corev1.PodSpec, runtime gpuv1alpha1.Runtime, runtimeClass string) {
	if runtime == gpuv1alpha1.Containerd {
		if runtimeClass == "" {
			runtimeClass = DefaultRuntimeClass
		}
		podSpec.RuntimeClassName = &runtimeClass
	}
}

func applyCommonDaemonsetMetadata(daemonsetObj *appsv1.DaemonSet, configDsSpec *gpuv1alpha1.DaemonSetsSpec) {
	if len(configDsSpec.Labels) > 0 {
		if daemonsetObj.Spec.Template.ObjectMeta.Labels == nil {
			daemonsetObj.Spec.Template.ObjectMeta.Labels = make(map[string]string)
		}
		for key, value := range configDsSpec.Labels {
			if key == "app" || key == "app.kubernetes.io/part-of" {
				continue
			}
			daemonsetObj.Spec.Template.ObjectMeta.Labels[key] = value
		}
	}

	if len(configDsSpec.Annotations) > 0 {
		if daemonsetObj.Spec.Template.ObjectMeta.Annotations == nil {
			daemonsetObj.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
		}
		for key, value := range configDsSpec.Annotations {
			daemonsetObj.Spec.Template.ObjectMeta.Annotations[key] = value
		}
	}
}

// generate hash for Annotations
func getDaemonSetHash(daemonSet *appsv1.DaemonSet) string {
	hasher := fnv.New32a()
	printer := spew.ConfigState{
		Indent:         " ",
		SortKeys:       true,
		DisableMethods: true,
		SpewKeys:       true,
	}
	printer.Fprintf(hasher, "%#v", daemonSet)
	return fmt.Sprint(hasher.Sum32())
}

func checkDaemonSetReady(name string, c GPUClusterController) gpuv1alpha1.State {
	ctx := c.ctx

	fmt.Println("checking daemonSet for readiness:", name)
	ds := &appsv1.DaemonSet{}
	err := c.client.Get(ctx, types.NamespacedName{
		Namespace: c.namespace,
		Name:      name,
	}, ds)
	if err != nil {
		fmt.Printf("failed to get daemonset: %s", name)
		return gpuv1alpha1.NotReady
	}

	// 检查存在期望pod为0
	if ds.Status.DesiredNumberScheduled == 0 {
		fmt.Printf("DaemonSet has desired pods of 0: %s", name)
		return gpuv1alpha1.Ready
	}
	// 检查是否有不可用的pod
	if ds.Status.NumberUnavailable != 0 {
		fmt.Printf("DaemonSet not ready: %s", name)
		return gpuv1alpha1.NotReady
	}

	// 如果daemonSet的更新策略不是使用"OnDelete",即使存在一些pod更新，ds仍然可以视为就绪状态
	if ds.Spec.UpdateStrategy.Type != appsv1.OnDeleteDaemonSetStrategyType {
		return gpuv1alpha1.NotReady
	}

	// 检查ds中的pod数量
	opts := []client.ListOption{client.MatchingLabels(ds.Spec.Template.ObjectMeta.Labels)}
	list := &corev1.PodList{}
	err = c.client.List(ctx, list, opts...)
	if err != nil {
		fmt.Println("failed to get PodList", err)
		return gpuv1alpha1.NotReady
	}
	if len(list.Items) == 0 {
		return gpuv1alpha1.NotReady
	}

	daemonSetRevision, err := getDaemonSetControllerRevisionHash(ctx, ds, c)
	if err != nil {
		fmt.Println("failed to get revision hash", err)
		return gpuv1alpha1.NotReady
	}

	for _, pod := range list.Items {
		podRevisionHash, err := getPodControllerRevisionHash(&pod)
		if err != nil {
			fmt.Println("Failed to get pod template revision hash", "pod", pod)
			return gpuv1alpha1.NotReady
		}

		if podRevisionHash != daemonSetRevision || pod.Status.Phase != "Running" {
			return gpuv1alpha1.NotReady
		}

		// If the pod generation matches the daemonset generation and the pod is running
		// and it has at least 1 container
		if len(pod.Status.ContainerStatuses) != 0 {
			for i := range pod.Status.ContainerStatuses {
				if !pod.Status.ContainerStatuses[i].Ready {
					// Return false if at least 1 container isn't ready
					return gpuv1alpha1.NotReady
				}
			}
		}
	}

	return gpuv1alpha1.Ready
}

func getDaemonSetControllerRevisionHash(ctx context.Context, daemonSet *appsv1.DaemonSet, c GPUClusterController) (string, error) {
	// get all revisions for the daemonset
	opts := []client.ListOption{
		client.MatchingLabels(daemonSet.Spec.Selector.MatchLabels),
		client.InNamespace(c.namespace),
	}
	list := &appsv1.ControllerRevisionList{}
	err := c.client.List(ctx, list, opts...)
	if err != nil {
		return "", fmt.Errorf("error getting controller revision list for daemonset %s: %v", daemonSet.Name, err)
	}

	var revisions []appsv1.ControllerRevision
	for _, controllerRevision := range list.Items {
		if strings.HasPrefix(controllerRevision.Name, daemonSet.Name) {
			revisions = append(revisions, controllerRevision)
		}
	}

	if len(revisions) == 0 {
		return "", fmt.Errorf("no revision found for daemonset %s", daemonSet.Name)
	}

	// sort the revision list to make sure we obtain latest revision always
	sort.Slice(revisions, func(i, j int) bool { return revisions[i].Revision < revisions[j].Revision })

	currentRevision := revisions[len(revisions)-1]
	hash := strings.TrimPrefix(currentRevision.Name, fmt.Sprintf("%s-", daemonSet.Name))

	return hash, nil
}

func getPodControllerRevisionHash(pod *corev1.Pod) (string, error) {
	if hash, ok := pod.Labels[PodControllerRevisionHashLabelKey]; ok {
		return hash, nil
	}
	return "", fmt.Errorf("controller-revision-hash label not present for pod %s", pod.Name)
}

// func checkDaemonSetChanged(currentDs *appsv1.DaemonSet, newDs *appsv1.DaemonSet) bool {
// 	if currentDs == nil || newDs == nil {
// 		return true
// 	}
// 	fmt.Println(currentDs.Annotations)
// 	return true
// }

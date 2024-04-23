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

package controllers

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	gpuv1alpha1 "github.com/chen-mao/k8s-gpu-operator.git/api/v1alpha1"
)

var gpuClusterCtrl GPUClusterController

// GPUClusterReconciler reconciles a GPUCluster object
type GPUClusterReconciler struct {
	client.Client

	// Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=xdxct.com,resources=gpuclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=xdxct.com,resources=gpuclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=xdxct.com,resources=gpuclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles;clusterrolebindings;roles;rolebindings,verbs=*
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GPUCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *GPUClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	gpuObjects := gpuv1alpha1.GPUCluster{}
	err := r.Client.Get(ctx, req.NamespacedName, &gpuObjects)
	if err != nil {
		err = fmt.Errorf("failed to get gpucluster object: %v", err)
		if client.IgnoreNotFound(err) != nil {
			// cr not found and don't requeue
			return reconcile.Result{}, nil
		}
		// the requeue request
		return reconcile.Result{}, err
	}

	if gpuClusterCtrl.singleton != nil && gpuClusterCtrl.singleton.ObjectMeta.Name != gpuObjects.ObjectMeta.Name {
		gpuObjects.SetStatus(gpuv1alpha1.Ignored, gpuClusterCtrl.namespace)
		return ctrl.Result{}, err
	}

	err = gpuClusterCtrl.init(ctx, r, &gpuObjects)
	if err != nil {
		err = fmt.Errorf("failed to initialize ClusterPolicy controller: %v", err)
		return ctrl.Result{}, err
	}

	// loop: deploy componentes
	for {
		fmt.Println("<---------------->")
		// to remove
		if gpuClusterCtrl.index == 3 {
			break
		}
		status, err := gpuClusterCtrl.step()
		if err != nil {
			return ctrl.Result{
				RequeueAfter: time.Second * 10,
			}, nil
		}
		if status == gpuv1alpha1.NotReady {
			fmt.Println("Components Not Ready")
		} else {
			continue
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GPUClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gpuv1alpha1.GPUCluster{}).
		Complete(r)
}

/*
Copyright 2023.

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

	"github.com/go-logr/logr"
	"github.com/opentracing/opentracing-go/log"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	// rackn ""

	infrastructurev1alpha1 "github.com/adminturneddevops/cluster-api-provider-rackn/api/v1alpha1"
)

// RackNk8sclusterReconciler reconciles a RackNk8scluster object
type RackNk8sclusterReconciler struct {
	client.Client
	WatchFilterValue string
	Scheme           *runtime.Scheme
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=racknk8sclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=racknk8sclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=racknk8sclusters/finalizers,verbs=update

//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
// validate validates if context configuration has all required fields properly populated.

// clusterReconcileContext implements ReconcileContext by reconciling racknCluster object.

type clusterReconcileContext struct {
	ctx            context.Context
	racknCluster   *infrastructurev1alpha1.RackNk8scluster
	patchHelper    *patch.Helper
	cluster        *clusterv1.Cluster
	log            logr.Logger
	client         client.Client
	namespacedName types.NamespacedName
}

func (rcr *RackNk8sclusterReconciler) newReconcileContext(ctx context.Context, namespacedName types.NamespacedName) (*clusterReconcileContext, error) {
	log := ctrl.LoggerFrom(ctx)

	crc := &clusterReconcileContext{
		log:            log.WithValues("RackNk8scluster", namespacedName),
		ctx:            ctx,
		racknCluster:   &infrastructurev1alpha1.RackNk8scluster{},
		client:         rcr.Client,
		namespacedName: namespacedName,
	}

	if err := crc.client.Get(crc.ctx, namespacedName, crc.racknCluster); err != nil {
		if apierrors.IsNotFound(err) {
			crc.log.Info("RackNk8sCluster object not found")

			return nil, nil
		}

		return nil, fmt.Errorf("getting RackNCluster: %w", err)
	}

	patchHelper, err := patch.NewHelper(crc.racknCluster, crc.client)
	if err != nil {
		return nil, fmt.Errorf("initializing patch helper: %w", err)
	}

	crc.patchHelper = patchHelper

	cluster, err := util.GetOwnerCluster(crc.ctx, crc.client, crc.racknCluster.ObjectMeta)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("getting owner cluster: %w", err)
		}
	}

	if cluster == nil {
		crc.log.Info("OwnerCluster is not set yet.")
	}

	crc.cluster = cluster

	return crc, nil
}

const (

	// KubernetesAPIPort is a port used by clusters for Kubernetes API.
	KubernetesAPIPort = 6443
)

// Need to input a function around RackN resource checking.
// Similar to the below example of how TinkerBell is checking hardware resources
// The /pools/{id}/status to return the status of machines in the pool could be a good one here.

// func hardwareIP(hardware *tinkerbell.Hardware) (string, error) {
// 	if hardware == nil {
// 		return "", ErrHardwareIsNil
// 	}

// 	if len(hardware.Spec.Interfaces) == 0 {
// 		return "", ErrHardwareMissingInterfaces
// 	}

// 	if hardware.Spec.Interfaces[0].DHCP == nil {
// 		return "", ErrHardwareFirstInterfaceNotDHCP
// 	}

// 	if hardware.Spec.Interfaces[0].DHCP.IP == nil {
// 		return "", ErrHardwareFirstInterfaceDHCPMissingIP
// 	}

// 	if hardware.Spec.Interfaces[0].DHCP.IP.Address == "" {
// 		return "", ErrHardwareFirstInterfaceDHCPMissingIP
// 	}

// 	return hardware.Spec.Interfaces[0].DHCP.IP.Address, nil
// }

func (crc *clusterReconcileContext) controlPlaneEndpoint() (clusterv1.APIEndpoint, error) {
	switch {
	// If the ControlPlaneEndpoint is already configured, return it.
	case crc.racknCluster.Spec.ControlPlaneEndpoint.IsValid():
		return crc.racknCluster.Spec.ControlPlaneEndpoint, nil
	// If the ControlPlaneEndpoint on the cluster is already configured, return it.
	case crc.cluster.Spec.ControlPlaneEndpoint.IsValid():
		return crc.cluster.Spec.ControlPlaneEndpoint, nil
	// If the cluster isn't ready
	case crc.cluster == nil:
		return clusterv1.APIEndpoint{}, log.Error("Cluster Resource Not Ready")
	}

	endpoint := clusterv1.APIEndpoint{
		Host: crc.cluster.Spec.ControlPlaneEndpoint.Host,
		Port: crc.cluster.Spec.ControlPlaneEndpoint.Port,
	}

	if endpoint.Host == "" {
		endpoint.Host = crc.racknCluster.Spec.ControlPlaneEndpoint.Host
	}

	if endpoint.Port == 0 {
		endpoint.Port = crc.racknCluster.Spec.ControlPlaneEndpoint.Port
	}

	if endpoint.Host == "" {
		return endpoint, log.Error("Control Plane Endpoint Not Set")
	}

	if endpoint.Port == 0 {
		endpoint.Port = KubernetesAPIPort
	}

	return endpoint, nil
}

// Reconcile implements ReconcileContext interface by ensuring that all racknCluster object
// fields are properly populated.
func (crc *clusterReconcileContext) reconcile() error {
	controlPlaneEndpoint, err := crc.controlPlaneEndpoint()
	if err != nil {
		return err
	}

	// Ensure that we are setting the ControlPlaneEndpoint on the racknCluster
	// in the event that it was defined on the Cluster resource instead
	crc.racknCluster.Spec.ControlPlaneEndpoint.Host = controlPlaneEndpoint.Host
	crc.racknCluster.Spec.ControlPlaneEndpoint.Port = controlPlaneEndpoint.Port

	crc.racknCluster.Status.Ready = true

	crc.log.Info("Setting cluster status to ready")

	if err := crc.patchHelper.Patch(crc.ctx, crc.racknCluster); err != nil {
		return fmt.Errorf("patching cluster object: %w", err)
	}

	return nil
}

func (crc *clusterReconcileContext) reconcileDelete() error {
	return nil
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=racknk8sclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=racknk8sclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=racknk8sclusters;clusters/status,verbs=get;list;watch

// Reconcile ensures state of RackN clusters.
func (tcr *RackNk8sclusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	crc, err := tcr.newReconcileContext(ctx, req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("creating reconciliation context: %w", err)
	}

	if crc == nil {
		return ctrl.Result{}, nil
	}

	if !crc.racknCluster.ObjectMeta.DeletionTimestamp.IsZero() {
		if annotations.HasPaused(crc.racknCluster) {
			crc.log.Info("RackNCluster is marked as paused. Won't reconcile deletion")

			return ctrl.Result{}, nil
		}

		crc.log.Info("Removing cluster")

		return ctrl.Result{}, crc.reconcileDelete()
	}

	if crc.cluster == nil {
		return ctrl.Result{}, nil
	}

	if annotations.IsPaused(crc.cluster, crc.racknCluster) {
		crc.log.Info("RackNCluster is marked as paused. Won't reconcile")

		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, crc.reconcile()
}

// SetupWithManager configures reconciler with a given manager.
func (tcr *RackNk8sclusterReconciler) SetupWithManager(
	ctx context.Context,
	mgr ctrl.Manager,
	options controller.Options,
) error {
	log := ctrl.LoggerFrom(ctx)

	mapper := util.ClusterToInfrastructureMapFunc(
		ctx,
		infrastructurev1alpha1.GroupVersion.WithKind("RackNk8sCluster"),
		mgr.GetClient(),
		&infrastructurev1alpha1.RackNk8scluster{},
	)

	builder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrastructurev1alpha1.RackNk8scluster{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, tcr.WatchFilterValue)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(log)).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(mapper),
			builder.WithPredicates(predicates.ClusterUnpaused(log)),
		)

	if err := builder.Complete(tcr); err != nil {
		return fmt.Errorf("failed to configure controller: %w", err)
	}

	return nil
}

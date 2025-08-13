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

package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	jumpstarterdevv1alpha1 "github.com/the78mole/jumpstarter-mono/core/controller/api/v1alpha1"
	"github.com/the78mole/jumpstarter-mono/core/controller/internal/oidc"
)

// ClientReconciler reconciles a Client object
type ClientReconciler struct {
	kclient.Client
	Scheme *runtime.Scheme
	Signer *oidc.Signer
}

// +kubebuilder:rbac:groups=jumpstarter.dev,resources=clients,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=jumpstarter.dev,resources=clients/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=jumpstarter.dev,resources=clients/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *ClientReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var client jumpstarterdevv1alpha1.Client
	if err := r.Get(ctx, req.NamespacedName, &client); err != nil {
		return ctrl.Result{}, kclient.IgnoreNotFound(
			fmt.Errorf("Reconcile: failed to get client: %w", err),
		)
	}

	original := kclient.MergeFrom(client.DeepCopy())

	if err := r.reconcileStatusCredential(ctx, &client); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileStatusEndpoint(ctx, &client); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Status().Patch(ctx, &client, original); err != nil {
		return RequeueConflict(logger, ctrl.Result{}, err)
	}

	return ctrl.Result{}, nil
}

func (r *ClientReconciler) reconcileStatusCredential(
	ctx context.Context,
	client *jumpstarterdevv1alpha1.Client,
) error {
	secret, err := ensureSecret(ctx, kclient.ObjectKey{
		Name:      client.Name + "-client",
		Namespace: client.Namespace,
	}, r.Client, r.Scheme, r.Signer, client.InternalSubject(), client)
	if err != nil {
		return fmt.Errorf("reconcileStatusCredential: failed to prepare credential for client: %w", err)
	}
	client.Status.Credential = &corev1.LocalObjectReference{
		Name: secret.Name,
	}
	return nil
}

// nolint:unparam
func (r *ClientReconciler) reconcileStatusEndpoint(
	ctx context.Context,
	client *jumpstarterdevv1alpha1.Client,
) error {
	logger := log.FromContext(ctx)

	endpoint := controllerEndpoint()
	if client.Status.Endpoint != endpoint {
		logger.Info("reconcileStatusEndpoint: updating controller endpoint")
		client.Status.Endpoint = endpoint
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClientReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jumpstarterdevv1alpha1.Client{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

/*


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

	"github.com/go-logr/logr"
	"github.com/lyft/flytestdlib/logger"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1alpha1 "samza-k8s-operator/api/v1alpha1"
	controllerconfig "samza-k8s-operator/controllers/config"
	"samza-k8s-operator/controllers/kubeutils"
)

// Key is used as a type for context keys
type Key string

const (
	// AppNameKey is the key used to ref Application Name
	AppNameKey Key = "app_name"
	// NamespaceKey is the key used to ref Namespace
	NamespaceKey Key = "ns"
	// PhaseKey is the key used to ref Phase
	PhaseKey Key = "phase"
)

// SamzaApplicationReconciler reconciles a SamzaApplication object
type SamzaApplicationReconciler struct {
	client            client.Client
	cache             cache.Cache
	log               logr.Logger
	scheme            *runtime.Scheme
	samzaStateMachine SamzaApplicationStateMachineInterface
	eventRecorder     record.EventRecorder
}

// NewSamzaApplicationReconciler returns a new instance of SamzaApplicationReconciler
func NewSamzaApplicationReconciler(mgr manager.Manager) (*SamzaApplicationReconciler, error) {
	samzaStateMachine := NewSamzaApplicationStateMachine(mgr)
	return &SamzaApplicationReconciler{
		client:            mgr.GetClient(),
		cache:             mgr.GetCache(),
		log:               ctrl.Log.WithName("controllers").WithName(v1alpha1.SamzaApplicationKind),
		scheme:            mgr.GetScheme(),
		eventRecorder:     mgr.GetEventRecorderFor(v1alpha1.OperatorName),
		samzaStateMachine: samzaStateMachine,
	}, nil
}

func (reconciler *SamzaApplicationReconciler) getResource(ctx context.Context, key types.NamespacedName, obj runtime.Object) error {
	err := reconciler.cache.Get(ctx, key, obj)
	if err != nil && kubeutils.IsKubeObjectNotExist(err) {
		return reconciler.client.Get(ctx, key, obj)
	}
	return err
}

// For failures, we do not want to retry immediately, as we want the underlying resource to recover.
// At the same time, we want to retry faster than the regular success interval.
func (reconciler *SamzaApplicationReconciler) getFailureRetryInterval() time.Duration {
	return controllerconfig.ResyncPeriod / 2
}

// If there is an error we requeue and check again in the next loop else we return sucess
func (reconciler *SamzaApplicationReconciler) getReconcileResultForError(err error) reconcile.Result {
	if err == nil {
		return reconcile.Result{}
	}
	return reconcile.Result{
		RequeueAfter: reconciler.getFailureRetryInterval(),
	}
}

// +kubebuilder:rbac:groups=samzaoperator.samza.apache.org,resources=samzaapplications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=samzaoperator.samza.apache.org,resources=samzaapplications/status,verbs=get;update;patch

// Reconcile the observed state towards the desired state for a SamzaApplication
func (reconciler *SamzaApplicationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, NamespaceKey, req.Namespace)
	ctx = context.WithValue(ctx, AppNameKey, req.Name)
	// var log = reconciler.log.WithValues(v1alpha1.SamzaApplicationKind, req.NamespacedName)

	typeMeta := metaV1.TypeMeta{
		Kind:       v1alpha1.SamzaApplicationKind,
		APIVersion: v1alpha1.GroupVersion.String(),
	}

	instance := &v1alpha1.SamzaApplication{
		TypeMeta: typeMeta,
	}

	err := reconciler.getResource(ctx, req.NamespacedName, instance)
	if err != nil {
		if kubeutils.IsKubeObjectNotExist(err) {
			// Deleted before Reconcile was called - do nothing
			return reconcile.Result{}, nil
		}
		// Error reading the object - we will check again in next loop
		return reconciler.getReconcileResultForError(err), nil
	}
	// We are seeing instances where getResource is removing TypeMeta
	instance.TypeMeta = typeMeta
	reconciler.log.Info(fmt.Sprintf("Starting reconcile loop for SamzaApplication: %s", req.NamespacedName))

	ctx = context.WithValue(ctx, PhaseKey, string(instance.Status.Phase))
	err = reconciler.samzaStateMachine.Handle(ctx, instance)
	if err != nil {
		logger.Warnf(ctx, "Failed to reconcile resource %v: %v", req.NamespacedName, err)
	}
	return reconciler.getReconcileResultForError(err), err
}

// SetupWithManager registers this reconciler with the controller manager and starts watching SamzaApplication resources.
func (reconciler *SamzaApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.SamzaApplication{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&batchv1.Job{}).
		Complete(reconciler)
}

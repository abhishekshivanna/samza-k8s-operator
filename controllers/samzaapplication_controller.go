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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	samzaoperatorv1alpha1 "samza-k8s-operator/api/v1alpha1"
)

// SamzaApplicationReconciler reconciles a SamzaApplication object
type SamzaApplicationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=samzaoperator.samza.apache.org,resources=samzaapplications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=samzaoperator.samza.apache.org,resources=samzaapplications/status,verbs=get;update;patch

// Reconcile the observed state towards the desired state for a SamzaApplication
func (r *SamzaApplicationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("samzaapplication", req.NamespacedName)

	r.Log.Info(fmt.Sprintf("here: %s", req))
	return ctrl.Result{}, nil
}

// SetupWithManager registers this reconciler with the controller manager and starts watching SamzaApplication resources.
func (r *SamzaApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&samzaoperatorv1alpha1.SamzaApplication{}).
		Complete(r)
}

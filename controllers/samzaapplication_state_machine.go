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
	"samza-k8s-operator/api/v1alpha1"
	"samza-k8s-operator/controllers/k8s"
	"samza-k8s-operator/controllers/samza"

	ctrl "sigs.k8s.io/controller-runtime"
)

// SamzaApplicationStateMachineInterface is the interface for the state machine of this operator
type SamzaApplicationStateMachineInterface interface {
	Handle(ctx context.Context, application *v1alpha1.SamzaApplication) error
}

// NewSamzaApplicationStateMachine creates a new SamzaApplicationStateMachine
func NewSamzaApplicationStateMachine(mgr ctrl.Manager) SamzaApplicationStateMachineInterface {
	client := k8s.NewClient(mgr)
	return &SamzaApplicationStateMachine{
		samzaController: samza.NewController(client, mgr),
	}
}

// SamzaApplicationStateMachine holds the context and state to handle
// reconcile request.
type SamzaApplicationStateMachine struct {
	samzaController samza.ControllerInterface
}

// Handle funtion is used apply desired state on the kubernetes cluster based on the status
// of the application deployment
func (handler *SamzaApplicationStateMachine) Handle(ctx context.Context, application *v1alpha1.SamzaApplication) error {
	return nil
}

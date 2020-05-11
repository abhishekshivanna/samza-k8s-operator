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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SamzaApplicationKind = "SamzaApplication"
	OperatorName         = "SamzaOperator"
)

// ImageSpec defines the Samza application image for the JobCoordinator and SamzaContainers
type ImageSpec struct {
	Name string `json:"name"`

	// Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always
	// if :latest tag is specified, or IfNotPresent otherwise.
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
}

// JobCoordinatorPorts defines the port of the JobCoordinator
type JobCoordinatorPorts struct {
	// RPC port
	RPC *int32 `json:"rpc,omitempty"`

	// UI port
	UI *int32 `json:"ui,omitempty"`
}

// JobCoordinatorSpec defines the properties of the JobCoordinator
type JobCoordinatorSpec struct {
	Ports JobCoordinatorPorts `json:"jobCoordinatorPorts,omitempty"`

	// Compute resources required by the JobCoordinator container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Volumes for the JobCoordinator pod.
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// VolumeMounts on the JobCoordinator container.
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
}

// SamzaContainerSpec defines the properties of the SamzaContainers
type SamzaContainerSpec struct {
	// The number of SamzaContainers or replicas requied for processing
	Replicas int32 `json:"replicas"`

	// Compute resources required by the SamzaContainer container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Volumes for the SamzaContainer pod.
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// VolumeMounts on the SamzaContainer container.
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`

	// TODO: We need specs for 1) host-affinity 2) node labels
}

// SamzaApplicationSpec defines the desired state of SamzaApplication
type SamzaApplicationSpec struct {

	// Samza application's image spec
	Image ImageSpec `json:"image"`

	// Job Coordinator spec
	JobCoordinator JobCoordinatorSpec `json:"jobCoordinator"`

	// Samza Container spec
	SamzaContainer SamzaContainerSpec `json:"samzaContainer"`

	// Environment variables shared by all JobCoordinator and SamzaContainer containers.
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`

	// Instance number of the Samza Application
	ApplicationInstance uint32 `json:"applicationInstance"`
}

// SamzaApplicationPhase is used to represent the current phase of the application deployment
type SamzaApplicationPhase string

// SamzaApplicationStatus defines the observed state of SamzaApplication
type SamzaApplicationStatus struct {
	// Phase is the current phase of the state machine the application deployment is in
	Phase SamzaApplicationPhase `json:"phase"`
}

// GetPhase is used to get the phase of the Samza Application deployment
func (status *SamzaApplicationStatus) GetPhase() SamzaApplicationPhase {
	return status.Phase
}

// UpdatePhase is used to update the phase of the Samza Application deployment
func (status *SamzaApplicationStatus) UpdatePhase(phase SamzaApplicationPhase) {
	status.Phase = phase
}

const (
	// SamzaApplicationNew represents a brand new application in the cluster
	SamzaApplicationNew SamzaApplicationPhase = ""
)

// +kubebuilder:object:root=true

// SamzaApplication is the Schema for the samzaapplications API
type SamzaApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SamzaApplicationSpec   `json:"spec,omitempty"`
	Status SamzaApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SamzaApplicationList contains a list of SamzaApplication
type SamzaApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SamzaApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SamzaApplication{}, &SamzaApplicationList{})
}

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

package k8s

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// ClusterInterface is used to interact with the Kuberneters Cluster
type ClusterInterface interface {
	CreateK8sObject(ctx context.Context, object runtime.Object) error
}

// Cluster is an implementation of the ClusterInterface
type Cluster struct {
	cache  cache.Cache
	client client.Client
}

// NewCluster creates a new Cluster instances
func NewCluster(mgr manager.Manager) ClusterInterface {
	return &Cluster{
		cache:  mgr.GetCache(),
		client: mgr.GetClient(),
	}
}

// CreateK8sObject creates a new runtime object on the K8s Cluster
func (k *Cluster) CreateK8sObject(ctx context.Context, object runtime.Object) error {
	objCreate := object.DeepCopyObject()
	err := k.client.Create(ctx, objCreate)
	if err != nil {
		return err
	}
	return nil
}

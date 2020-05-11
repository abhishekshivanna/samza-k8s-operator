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

package kubeutils

import (
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

// IsKubeObjectDoesNotExist returns true if error indicates that the resouce
// is 1) not found, 2) is gone or 3) is expired
func IsKubeObjectNotExist(err error) bool {
	return kerrors.IsNotFound(err) || kerrors.IsGone(err) || kerrors.IsResourceExpired(err)
}

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HTTPRouteBindingSpec defines the desired state of HTTPRouteBinding
type HTTPRouteBindingSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	VirtualServiceBaseRef VirtualServiceBaseRef `json:"virtualServiceBaseRef,omitempty"`
	HTTPRoute             HTTPRoute             `json:"httpRoute,omitempty"`
}

// VirtualServiceBaseRef is reference to VirtualServiceBase resource
type VirtualServiceBaseRef struct {
	metav1.TypeMeta `json:",inline"`

	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// HTTPRouteBindingStatus defines the observed state of HTTPRouteBinding
type HTTPRouteBindingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// HTTPRouteBinding is the Schema for the httproutebindings API
type HTTPRouteBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HTTPRouteBindingSpec   `json:"spec,omitempty"`
	Status HTTPRouteBindingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HTTPRouteBindingList contains a list of HTTPRouteBinding
type HTTPRouteBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HTTPRouteBinding `json:"items"`
}

func (ref VirtualServiceBaseRef) IsReference(base *VirtualServiceBase) bool {
	return ref.APIVersion == base.APIVersion && ref.Kind == base.Kind && ref.Name == base.Name && ref.Namespace == base.Namespace
}

func init() {
	SchemeBuilder.Register(&HTTPRouteBinding{}, &HTTPRouteBindingList{})
}

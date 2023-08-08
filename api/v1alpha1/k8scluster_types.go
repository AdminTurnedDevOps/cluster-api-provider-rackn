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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const ClusterFinalizer = "k8scluster.infrastructure.cluster.x-k8s.io"

// K8sclusterSpec defines the desired state of K8scluster
type K8sclusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of K8scluster. Edit k8scluster_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// K8sclusterStatus defines the observed state of K8scluster
type K8sclusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// K8scluster is the Schema for the k8sclusters API
type K8scluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   K8sclusterSpec   `json:"spec,omitempty"`
	Status K8sclusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// K8sclusterList contains a list of K8scluster
type K8sclusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []K8scluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&K8scluster{}, &K8sclusterList{})
}

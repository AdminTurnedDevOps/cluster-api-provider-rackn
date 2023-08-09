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

//The purpose of the infrastructure cluster is to create any
//required infrastructure for the cluster but not anything related to individual machines/nodes

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// Allows cleanup of resources associated with RackNk8sCluster
// Finalizers are used to mark an object to prevent Kubernetes from deleting it until the finalizer is removed
const ClusterFinalizer = "racknk8scluster.infrastructure.cluster.x-k8s.io"

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RackNk8sclusterSpec defines the desired state of RackNk8scluster
type RackNk8sclusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Omiting output to empty as the output isn't necassary to the end user

	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`
}

// RackNk8sclusterStatus defines the observed state of RackNk8scluster
type RackNk8sclusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Ready bool `json:"ready"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RackNk8scluster is the Schema for the racknk8sclusters API
type RackNk8scluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RackNk8sclusterSpec   `json:"spec,omitempty"`
	Status RackNk8sclusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RackNk8sclusterList contains a list of RackNk8scluster
type RackNk8sclusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RackNk8scluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RackNk8scluster{}, &RackNk8sclusterList{})
}

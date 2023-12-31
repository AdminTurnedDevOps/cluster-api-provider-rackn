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
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RackNk8smachineSpec defines the desired state of RackNk8smachine
type RackNk8smachineSpec struct {
	// Required field for Cluster API
	// Omiting output to empty as the output isn't necassary to the end user

	ProviderID string `json:"providerID,omitempty"`

	ClusterName string `json:"clusterName,omitempty"`

	PubSSHKey string `json:"pubSSHkey,omitempty"`

	VMSize string `json:"vmSize,omitempty"`
}

// RackNk8smachineStatus defines the observed state of RackNk8smachine
type RackNk8smachineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Ready bool `json:"ready"`

	FailureMessage *string `json:"failureMessage,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RackNk8smachine is the Schema for the racknk8smachines API
type RackNk8smachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RackNk8smachineSpec   `json:"spec,omitempty"`
	Status RackNk8smachineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RackNk8smachineList contains a list of RackNk8smachine
type RackNk8smachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RackNk8smachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RackNk8smachine{}, &RackNk8smachineList{})
}

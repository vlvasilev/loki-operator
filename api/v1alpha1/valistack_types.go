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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ValiSpec defines the Vali configuration of the ValiStack
type ValiSpec struct {
	// AuthEnabled turns on Multitenancy
	AuthEnabled bool `json:"authEnabled,omitempty"`
	// Replicas is the number of the Vali replicas
	Replicas int32 `json:"replicas,omitempty"`
}

type HVPASpec struct {
	// Enabled enables the Vali HVPA
	Enabled bool `json:"enabled,omitempty"`
}

// EventLoggerSpec defines the event logger configuration of the ValiStack
type EventLoggerSpec struct {
	//TODO: fill me
}

// ValiStackSpec defines the desired state of ValiStack
type ValiStackSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ValiStack. Edit valistack_types.go to remove/update
	Foo string `json:"foo,omitempty"`
	// ValiSpec defines the Vali configuration of the ValiStack
	Vali *ValiSpec `json:"vali,omitempty"`
	// PriorityClassName defines the the ValiStack PriorityClassName
	PriorityClassName *string `json:"priorityClassName,omitempty"`
}

// ValiStackStatus defines the observed state of ValiStack
type ValiStackStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced,categories=logging,path=valistacks,shortName=vs

// ValiStack is the Schema for the valistacks API
type ValiStack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ValiStackSpec   `json:"spec,omitempty"`
	Status ValiStackStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ValiStackList contains a list of ValiStack
type ValiStackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ValiStack `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ValiStack{}, &ValiStackList{})
}

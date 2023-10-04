/*
Copyright 2023 Red Hat, Inc.

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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenShiftControlPlane is responsible for bootstrapping an OpenShift cluster control plane.
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:subresource:status
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type OpenShiftControlPlane struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of the OpenShiftControlPlane.
	// +kubebuilder:validation:Required
	Spec OpenShiftControlPlaneSpec `json:"spec"`

	// status is the observed state of the OpenShiftControlPlane.
	// +optional
	Status OpenShiftControlPlaneStatus `json:"status,omitempty"`
}

// OpenShiftControlPlaneSpec is the spec of the OpenShift control plane resource.
type OpenShiftControlPlaneSpec struct {
	// machineTemplate defines the machine template used to create the initial bootstrap and control plane machines.
	// Continued management of the control plane machines will be handled by the control plane machine set.
	// The machine template is therefore immutable and only applicable during the bootstrap process.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="machineTemplate is immutable"
	// +kubeubilder:validation:Required
	// + ---
	// + This field, and the fields within the OpenShiftControlPlaneMachineTemplate, are required as part of the Cluster API control plane API contract.
	MachineTemplate OpenShiftControlPlaneMachineTemplate `json:"machineTemplate"`

	// installStateSecretRef is a reference to a secret containing the install state.
	// The install state secret must contain either the install config or the install state, or both.
	// The install state secret must be in the same namespace as the OpenShiftControlPlane.
	// The install config must be under the key `install-config.yaml` and the install state must be under the key `.openshift_install_state.json`.
	// These files will be passed to the installer to generate the ignition configs for the bootstrap node, control plane nodes and worker nodes.
	// +kubebuilder:validation:Required
	InstallStateSecretRef OpenShiftControlPlaneSecretRef `json:"installStateSecretRef"`

	// manifestsSelector is a selector to identify secrets containing manifests to be included in the ignition generation phase.
	// The selector must match the labels on the secrets to be injected.
	// Each key in the secret must be the path to a file to be injected into the ignition.
	// This path should start with either `manifests/` or `openshift/`.
	// When omitted, the default manifests generated by the installer will be used.
	// +optional
	ManifestsSelector metav1.LabelSelector `json:"manifestsSelector,omitempty"`
}

// OpenShiftControlPlaneMachineTemplate is the spec of the OpenShift control plane machines.
type OpenShiftControlPlaneMachineTemplate struct {
	// metadata is the standard object's metadata.
	// This allows for machine labels and annotations to be applied to the control plane machines.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta ObjectMeta `json:"metadata,omitempty"`

	// infrastructureRef is a required reference to a custom resource offered by an infrastructure provider.
	// The infrastructure reference should define a template for the infrastructure provider to create the bootstrap and control plane nodes.
	// +kubebuilder:validation:Required
	InfrastructureRef InfrastructureReference `json:"infrastructureRef"`

	// NodeDrainTimeout is the total amount of time that the controller will spend on draining a controlplane node
	// The default value is 0, meaning that the node can be drained without any time limitations.
	// NOTE: NodeDrainTimeout is different from `kubectl drain --timeout`
	// +optional
	NodeDrainTimeout *metav1.Duration `json:"nodeDrainTimeout,omitempty"`

	// NodeVolumeDetachTimeout is the total amount of time that the controller will spend on waiting for all volumes
	// to be detached. The default value is 0, meaning that the volumes can be detached without any time limitations.
	// +optional
	NodeVolumeDetachTimeout *metav1.Duration `json:"nodeVolumeDetachTimeout,omitempty"`

	// NodeDeletionTimeout defines how long the machine controller will attempt to delete the Node that the Machine
	// hosts after the Machine is marked for deletion. A duration of 0 will retry deletion indefinitely.
	// If no value is provided, the default value for this property of the Machine resource will be used.
	// +optional
	NodeDeletionTimeout *metav1.Duration `json:"nodeDeletionTimeout,omitempty"`
}

// ObjectMeta is a subset of metav1.ObjectMeta.
// We use this to customise the labels and annotations applied to control plane machines.
type ObjectMeta struct {
	// labels is a map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers and services.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// InfrastructureReference is a reference to a custom resource offered by an infrastructure provider.
// This is a subset of corev1.ObjectReference.
// The namespace must be set to the same as the OpenShiftControlPlane, but is required by Cluster API.
// Upstream discussion: https://github.com/kubernetes-sigs/cluster-api/issues/6539
type InfrastructureReference struct {
	// kind of the referent.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`

	// namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`

	// name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// apiVersion of the referent.
	// +kubebuilder:validation:Required
	APIVersion string `json:"apiVersion"`
}

// OpenShiftControlPlaneSecretRef is the reference to a secret in the same namespace as the OpenShiftControlPlane.
type OpenShiftControlPlaneSecretRef struct {
	// name is the name of the secret.
	// It has a maximum length of 253 characters and must be a valid DNS subdomain name.
	// It must consist only of lowercase alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character.
	// +kubebuilder:validation:Pattern="[a-z0-9]([-.a-z0-9]{,251}[a-z0-9])?"
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// OpenShiftControlPlaneStatus contains fields to describe the state of the OpenShiftControlPlane state.
type OpenShiftControlPlaneStatus struct {
	// conditions represents the observations of the OpenShiftControlPlane's current state.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// initialized denotes whether or not the control plane has been initialized.
	// This value will be set true once the first control plane node has joined the bootstrap control plane.
	// +optional
	// + ---
	// + This field is required as part of the Cluster API control plane API contract.
	Initialized bool `json:"initialized"`

	// ready denotes whether or not the control plane has has reached a ready state.
	// This value will be set true once the bootstrap node has completed the cluster bootstrap and the bootstrap node has been shut down.
	// +optional
	// + ---
	// + This field is required as part of the Cluster API control plane API contract.
	Ready bool `json:"ready"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenShiftControlPlaneList contains a list of OpenShiftControlPlane
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
type OpenShiftControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// items contains a list of OpenShiftControlPlanes.
	Items []OpenShiftControlPlane `json:"items"`
}

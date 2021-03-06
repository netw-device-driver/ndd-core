/*
Copyright 2021 Wim Henderickx.

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

package v1

import (
	nddv1 "github.com/netw-device-driver/ndd-runtime/apis/common/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PackageRevisionDesiredState is the desired state of the package revision.
type PackageRevisionDesiredState string

const (
	// PackageRevisionActive is an active package revision.
	PackageRevisionActive PackageRevisionDesiredState = "Active"

	// PackageRevisionInactive is an inactive package revision.
	PackageRevisionInactive PackageRevisionDesiredState = "Inactive"
)

// PackageRevisionSpec defines the desired state of Revision
type PackageRevisionSpec struct {
	// ControllerConfigRef references a ControllerConfig resource that will be
	// used to configure the packaged controller Deployment.
	// +optional
	ControllerConfigReference *nddv1.Reference `json:"controllerConfigRef,omitempty"`

	// DesiredState of the PackageRevision. Can be either Active or Inactive.
	DesiredState PackageRevisionDesiredState `json:"desiredState"`

	// Package image used by install Pod to extract package contents.
	PackageImage string `json:"packageImage"`

	// PackagePullSecrets are named secrets in the same namespace that can be
	// used to fetch packages from private registries. They are also applied to
	// any images pulled for the package, such as a provider's controller image.
	// +optional
	PackagePullSecrets []corev1.LocalObjectReference `json:"packagePullSecrets,omitempty"`

	// PackagePullPolicy defines the pull policy for the package. It is also
	// applied to any images pulled for the package, such as a provider's
	// controller image.
	// Default is IfNotPresent.
	// +optional
	// +kubebuilder:default=IfNotPresent
	PackagePullPolicy *corev1.PullPolicy `json:"packagePullPolicy,omitempty"`

	// Revision number. Indicates when the revision will be garbage collected
	// based on the parent's RevisionHistoryLimit.
	Revision int64 `json:"revision"`

	// AutoPilot specifies how the provider operates
	// When set to true the provider applies delta/diff changes to the device
	// manged resources automatically, if set to false the provider will report
	// the delta and the operator should intervene what to do with the delta/diffs
	// Defaults to true. Can be disabled by explicitly setting to flase.
	// +optional
	// +kubebuilder:default=true
	AutoPilot *bool `json:"autoPilot,omitempty"`

	// SkipDependencyResolution indicates to the package manager whether to skip
	// resolving dependencies for a package. Setting this value to true may have
	// unintended consequences.
	// Default is false.
	// +optional
	// +kubebuilder:default=false
	SkipDependencyResolution *bool `json:"skipDependencyResolution,omitempty"`
}

// PackageRevisionStatus defines the observed state of a PackageRevision
type PackageRevisionStatus struct {
	nddv1.ConditionedStatus `json:",inline"`
	ControllerRef           nddv1.Reference `json:"controllerRef,omitempty"`

	// References to objects owned by PackageRevision.
	ObjectRefs []nddv1.TypedReference `json:"objectRefs,omitempty"`

	// Dependency information.
	FoundDependencies     int64 `json:"foundDependencies,omitempty"`
	InstalledDependencies int64 `json:"installedDependencies,omitempty"`
	InvalidDependencies   int64 `json:"invalidDependencies,omitempty"`

	// PermissionRequests made by this package. The package declares that its
	// controller needs these permissions to run. The RBAC manager is
	// responsible for granting them.
	PermissionRequests []rbacv1.PolicyRule `json:"permissionRequests,omitempty"`
}

// +kubebuilder:object:root=true
// +genclient
// +genclient:nonNamespaced

// A ProviderRevision that has been added to Network device Driver.
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="HEALTHY",type="string",JSONPath=".status.conditions[?(@.kind=='PackageHealthy')].status"
// +kubebuilder:printcolumn:name="REVISION",type="string",JSONPath=".spec.revision"
// +kubebuilder:printcolumn:name="PKGIMAGE",type="string",JSONPath=".spec.packageImage"
// +kubebuilder:printcolumn:name="STATE",type="string",JSONPath=".spec.desiredState"
// +kubebuilder:printcolumn:name="DEP-FOUND",type="string",JSONPath=".status.foundDependencies"
// +kubebuilder:printcolumn:name="DEP-INSTALLED",type="string",JSONPath=".status.installedDependencies"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={ndd,pkg},shortName=pvr
type ProviderRevision struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageRevisionSpec   `json:"spec,omitempty"`
	Status PackageRevisionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProviderRevisionList contains a list of ProviderRevision.
type ProviderRevisionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderRevision `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProviderRevision{}, &ProviderRevisionList{})
}

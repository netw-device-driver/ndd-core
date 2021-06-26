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

// Package v1 contains API Schema definitions for the pkg v1 API group
//+kubebuilder:object:generate=true
//+groupName=pkg.ndd.henderiw.be
package v1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// Package type metadata.
const (
	Group   = "pkg.ndd.henderiw.be"
	Version = "v1"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: Group, Version: Version}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// Provider type metadata.
var (
	ProviderKind             = reflect.TypeOf(Provider{}).Name()
	ProviderGroupKind        = schema.GroupKind{Group: Group, Kind: ProviderKind}.String()
	ProviderKindAPIVersion   = ProviderKind + "." + GroupVersion.String()
	ProviderGroupVersionKind = GroupVersion.WithKind(ProviderKind)
)

// ProviderRevision type metadata.
var (
	ProviderRevisionKind             = reflect.TypeOf(ProviderRevision{}).Name()
	ProviderRevisionGroupKind        = schema.GroupKind{Group: Group, Kind: ProviderRevisionKind}.String()
	ProviderRevisionKindAPIVersion   = ProviderRevisionKind + "." + GroupVersion.String()
	ProviderRevisionGroupVersionKind = GroupVersion.WithKind(ProviderRevisionKind)
)

// Lock type metadata.
var (
	LockKind             = reflect.TypeOf(Lock{}).Name()
	LockGroupKind        = schema.GroupKind{Group: Group, Kind: LockKind}.String()
	LockKindAPIVersion   = LockKind + "." + GroupVersion.String()
	LockGroupVersionKind = GroupVersion.WithKind(LockKind)
)

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
	"github.com/netw-device-driver/ndd-runtime/pkg/resource"
)

var _ Nn = &NetworkNode{}

// +k8s:deepcopy-gen=false
type Nn interface {
	resource.Object
	resource.Conditioned

	GetGrpcServerPort() int
	SetGrpcServerPort(p *int)

	GetDeviceDriverKind() DeviceDriverKind
	SetDeviceDriverKind(k *DeviceDriverKind)

	GetTargetAddress() string
	SetTargetAddress(a *string)

	GetTargetProxy() string
	SetTargetProxy(p *string)

	GetTargetCredentialsName() string
	SetTargetCredentialsName(c *string)

	GetTargetTLSCredentialsName() string
	SetTargetTLSCredentialsName(s *string)

	GetTargetSkipVerify() bool
	SetTargetSkipVerify(s *bool)

	GetTargetInsecure() bool
	SetTargetInsecure(s *bool)

	GetTargetEncoding() string
	SetTargetEncoding(e *string)
}

// GetCondition of this Network Node.
func (nn *NetworkNode) GetCondition(ct nddv1.ConditionKind) nddv1.Condition {
	return nn.Status.GetCondition(ct)
}

// SetConditions of this Network Node.
func (nn *NetworkNode) SetConditions(c ...nddv1.Condition) {
	nn.Status.SetConditions(c...)
}

func (nn *NetworkNode) GetGrpcServerPort() int {
	return *nn.Spec.GrpcServerPort
}

func (nn *NetworkNode) SetGrpcServerPort(p *int) {
	nn.Spec.GrpcServerPort = p
}

func (nn *NetworkNode) GetDeviceDriverKind() DeviceDriverKind {
	return *nn.Spec.DeviceDriverKind
}

func (nn *NetworkNode) SetDeviceDriverKind(k *DeviceDriverKind) {
	nn.Spec.DeviceDriverKind = k
}

func (nn *NetworkNode) GetTargetAddress() string {
	return *nn.Spec.Target.Address
}

func (nn *NetworkNode) SetTargetAddress(a *string) {
	nn.Spec.Target.Address = a
}

func (nn *NetworkNode) GetTargetProxy() string {
	return *nn.Spec.Target.Proxy
}

func (nn *NetworkNode) SetTargetProxy(p *string) {
	nn.Spec.Target.Proxy = p
}

func (nn *NetworkNode) GetTargetCredentialsName() string {
	return *nn.Spec.Target.CredentialsName
}

func (nn *NetworkNode) SetTargetCredentialsName(c *string) {
	nn.Spec.Target.CredentialsName = c
}

func (nn *NetworkNode) GetTargetTLSCredentialsName() string {
	return *nn.Spec.Target.TLSCredentialsName
}

func (nn *NetworkNode) SetTargetTLSCredentialsName(c *string) {
	nn.Spec.Target.TLSCredentialsName = c
}

func (nn *NetworkNode) GetTargetSkipVerify() bool {
	return *nn.Spec.Target.SkipVerify
}

func (nn *NetworkNode) SetTargetSkipVerify(s *bool) {
	nn.Spec.Target.SkipVerify = s
}

func (nn *NetworkNode) GetTargetInsecure() bool {
	return *nn.Spec.Target.Insecure
}

func (nn *NetworkNode) SetTargetInsecure(s *bool) {
	nn.Spec.Target.Insecure = s
}

func (nn *NetworkNode) GetTargetEncoding() string {
	return *nn.Spec.Target.Encoding
}

func (nn *NetworkNode) SetTargetEncoding(s *string) {
	nn.Spec.Target.Encoding = s
}
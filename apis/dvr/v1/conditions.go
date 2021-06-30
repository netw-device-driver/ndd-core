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
//+groupName=dvr.ndd.henderiw.be
package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nddv1 "github.com/netw-device-driver/ndd-runtime/apis/common/v1"
)

// Condition Kinds.
const (
	// A DeviceDriverInstalled indicates whether the device driver has been installed.
	ConditionKindDeviceDriverInstalled nddv1.ConditionKind = "DeviceDriverInstalled"

	// A DeviceDriverHealthy indicates whether the device driver is healthy.
	ConditionKindDeviceDriverHealthy nddv1.ConditionKind = "DeviceDriverHealthy"

	// A ConditionKindDeviceDriverReady indicates whether the device driver is discovered
	// and connected to the network device.
	ConditionKindDeviceDriverReady nddv1.ConditionKind = "DeviceDriverReady"
)

// ConditionReasons a package is or is not installed.
const (
	ConditionReasonInactive         nddv1.ConditionReason = "InactiveDeviceDriver"
	ConditionReasonActive           nddv1.ConditionReason = "ActiveDeviceDriver"
	ConditionReasonUnhealthy        nddv1.ConditionReason = "UnhealthyDeviceDriver"
	ConditionReasonHealthy          nddv1.ConditionReason = "HealthyDeviceDriver"
	ConditionReasonUnknownHealth    nddv1.ConditionReason = "UnknownDeviceDriverHealth"
	ConditionReasonDiscoveredReady  nddv1.ConditionReason = "DeviceDriverReady"
	ConditionReasonNotDiscovered    nddv1.ConditionReason = "UndiscoveredDeviceDriver"
	ConditionReasonUnknownDiscovery nddv1.ConditionReason = "UnknownDeviceDriverDiscovery"
)

// Inactive indicates that the device driver is waiting to be
// transitioned to an active state.
func Inactive() nddv1.Condition {
	return nddv1.Condition{
		Kind:               ConditionKindDeviceDriverInstalled,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonInactive,
	}
}

// Active indicates that the device driver is is installed and activated
func Active() nddv1.Condition {
	return nddv1.Condition{
		Kind:               ConditionKindDeviceDriverInstalled,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonActive,
	}
}

// Unhealthy indicates that the device driver is unhealthy.
func Unhealthy() nddv1.Condition {
	return nddv1.Condition{
		Kind:               ConditionKindDeviceDriverHealthy,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonUnhealthy,
	}
}

// Healthy indicates that the device driver is healthy.
func Healthy() nddv1.Condition {
	return nddv1.Condition{
		Kind:               ConditionKindDeviceDriverHealthy,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonHealthy,
	}
}

// UnknownHealth indicates that the health of the device driver is unknown.
func UnknownHealth() nddv1.Condition {
	return nddv1.Condition{
		Kind:               ConditionKindDeviceDriverHealthy,
		Status:             corev1.ConditionUnknown,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonUnknownHealth,
	}
}

// Discovered indicates that the device driver is discovered.
func Discovered() nddv1.Condition {
	return nddv1.Condition{
		Kind:               ConditionKindDeviceDriverReady,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonDiscoveredReady,
	}
}

// NotDiscovered indicates that the device driver is not discovered.
func NotDiscovered() nddv1.Condition {
	return nddv1.Condition{
		Kind:               ConditionKindDeviceDriverReady,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonNotDiscovered,
	}
}

// UnknownDiscovery indicates that the device driver discovery is unknown.
func UnknownDiscovery() nddv1.Condition {
	return nddv1.Condition{
		Kind:               ConditionKindDeviceDriverReady,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonUnknownDiscovery,
	}
}

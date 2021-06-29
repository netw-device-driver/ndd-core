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

package networknode

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	// Errors
	errCreateService = "create service error"
	errUpdateService = "update service error"
	errDeleteService = "delete service error"
	errGetService    = "get service error"
)

func (e *APIEstablisher) createService(ctx context.Context, name string, port int) error {
	log := e.log.WithValues("createServiceName", name)
	log.Debug("creating service...")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ndd-service-" + name,
			Namespace: e.namespace,
			Labels: map[string]string{
				"netwDDriver": "ndd-" + name,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "ndd-" + name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "proxy",
					Port:       int32(port),
					TargetPort: intstr.FromInt(port),
					Protocol:   "TCP",
				},
			},
		},
	}

	err := e.client.Create(ctx, s)
	if err != nil {
		return errors.Wrap(err, errCreateService)
	}
	log.Debug("created service...")

	return nil
}

func (e *APIEstablisher) updateService(ctx context.Context, name string, port int) error {
	log := e.log.WithValues("updateServiceName", name)
	log.Debug("updating service...")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ndd-service-" + name,
			Namespace: e.namespace,
			Labels: map[string]string{
				"netwDDriver": "ndd-" + name,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "ndd-" + name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "proxy",
					Port:       int32(port),
					TargetPort: intstr.FromInt(port),
					Protocol:   "TCP",
				},
			},
		},
	}
	err := e.client.Update(ctx, s)
	if err != nil {
		return errors.Wrap(err, errUpdateService)
	}
	log.Debug("updated service...")

	return nil
}

func (e *APIEstablisher) deleteService(ctx context.Context, name string) error {
	log := e.log.WithValues("deleteServiceName", name)
	log.Debug("deleting service...")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ndd-service-" + name,
			Namespace: e.namespace,
			Labels: map[string]string{
				"netwDDriver": "ndd-" + name,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "ndd-" + name,
			},
		},
	}
	err := e.client.Delete(ctx, s)
	if err != nil {
		return errors.Wrap(err, errDeleteService)
	}
	log.Debug("deleted service...")
	return nil
}

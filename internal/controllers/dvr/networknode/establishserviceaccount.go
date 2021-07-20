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
)

const (
	// Errors
	errCreateServiceAccount = "create ServiceAccount error"
	errUpdateServiceAccount = "update ServiceAccount error"
	errDeleteServiceAccount = "delete ServiceAccount error"
	errGetServiceAccount    = "get ServiceAccount error"
)

func (e *APIEstablisher) createServiceAccount(ctx context.Context, name string) error {
	log := e.log.WithValues("createServiceAccountName", name)
	log.Debug("creating ServiceAccount...")

	s := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: e.namespace,
		},
	}
	err := e.client.Create(ctx, s)
	if err != nil {
		return errors.Wrap(err, errCreateServiceAccount)
	}
	log.Debug("created ServiceAccount...")
	return nil
}

func (e *APIEstablisher) deleteServiceAccount(ctx context.Context, name string) error {
	log := e.log.WithValues("deleteServiceAccountName", name)
	log.Debug("deleting ServiceAccount...")

	s := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: e.namespace,
		},
	}
	err := e.client.Delete(ctx, s)
	if err != nil {
		return errors.Wrap(err, errDeleteServiceAccount)
	}
	log.Debug("deleted ServiceAccount...")
	return nil
}

func (e *APIEstablisher) updateServiceAccount(ctx context.Context, name string) error {
	log := e.log.WithValues("updateServiceAccountName", name)
	log.Debug("updating ServiceAccount...")

	s := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: e.namespace,
		},
	}

	err := e.client.Update(ctx, s)
	if err != nil {
		return errors.Wrap(err, errUpdateServiceAccount)
	}
	log.Debug("updated ServiceAccount...")

	return nil
}

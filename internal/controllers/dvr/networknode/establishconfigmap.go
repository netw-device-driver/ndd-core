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
	errCreateConfigMap = "create configmap error"
	errUpdateConfigMap = "update configmap error"
	errDeleteConfigMap = "delete configmap error"
	errGetConfigMap    = "get configmap error"
)

func (e *APIEstablisher) createConfigMap(ctx context.Context, name string) error {
	log := e.log.WithValues("createConfigMapName", name)
	log.Debug("creating configmap...")
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ndd-cm-" + name,
			Namespace: e.namespace,
			Labels: map[string]string{
				"netwDDriver": "ndd-" + name,
			},
		},
	}
	err := e.client.Create(ctx, cm)
	if err != nil {
		return errors.Wrap(err, errCreateConfigMap)
	}
	log.Debug("created configmap...")
	return nil
}

func (e *APIEstablisher) deleteConfigMap(ctx context.Context, name string) error {
	log := e.log.WithValues("deleteConfigMapName", name)
	log.Debug("deleting configmap...")

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ndd-cm-" + name,
			Namespace: e.namespace,
			Labels: map[string]string{
				"netwDDriver": "ndd-" + name,
			},
		},
	}
	err := e.client.Delete(ctx, cm)
	if err != nil {
		return errors.Wrap(err, errDeleteConfigMap)
	}
	log.Debug("deleted configmap...")
	return nil
}

func (e *APIEstablisher) updateConfigMap(ctx context.Context, name string) error {
	log := e.log.WithValues("updateConfigMapName", name)
	log.Debug("updating configmap...")

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ndd-cm-" + name,
			Namespace: e.namespace,
			Labels: map[string]string{
				"netwDDriver": "ndd-" + name,
			},
		},
	}

	err := e.client.Update(ctx, cm)
	if err != nil {
		return errors.Wrap(err, errUpdateConfigMap)
	}
	log.Debug("updated configmap...")

	return nil
}

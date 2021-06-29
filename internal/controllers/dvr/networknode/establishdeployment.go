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
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Errors
	errCreateDeployment = "create deployment error"
	errUpdateDeployment = "update deployment error"
	errDeleteDeployment = "delete deployment error"
	errGetDeployment    = "get deployment error"
)

func (e *APIEstablisher) createDeployment(ctx context.Context, name string, c *corev1.Container) error {
	log := e.log.WithValues("createDeploymentName", name, "container", c)
	log.Debug("creating deployment...")
	deployment := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ndd-deployment-" + name,
			Namespace: e.namespace,
			Labels: map[string]string{
				"netwDDriver": "ndd-" + name,
			},
		},
		Spec: appv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "ndd-" + name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "ndd-" + name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						*c,
					},
				},
			},
		},
	}

	if err := e.client.Create(ctx, deployment); err != nil {
		return errors.Wrap(err, errCreateDeployment)
	}
	log.Debug("created deployment...")
	return nil
}

func (e *APIEstablisher) updateDeployment(ctx context.Context, name string, c *corev1.Container) error {
	log := e.log.WithValues("updateDeploymentName", name, "container", c)
	log.Debug("updating deployment...")
	deployment := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nddr-deployment-" + name,
			Namespace: e.namespace,
			Labels: map[string]string{
				"netwDDriver": "ndd-" + name,
			},
		},
		Spec: appv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "ndd-" + name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "ndd-" + name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						*c,
					},
				},
			},
		},
	}

	err := e.client.Update(ctx, deployment)
	if err != nil {
		return errors.Wrap(err, errUpdateDeployment)
	}
	log.Debug("updating deployment...")
	return nil
}

func (e *APIEstablisher) deleteDeployment(ctx context.Context, name string) error {
	log := e.log.WithValues("deleteDeploymentName", name)
	log.Debug("deleting deployment...")

	deployment := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ndd-deployment-" + name,
			Namespace: e.namespace,
			Labels: map[string]string{
				"netwDDriver": "ndd-" + name,
			},
		},
	}
	err := e.client.Delete(ctx, deployment)
	if err != nil {
		return errors.Wrap(err, errDeleteDeployment)

	}
	log.Debug("delete deployment...")
	return nil
}

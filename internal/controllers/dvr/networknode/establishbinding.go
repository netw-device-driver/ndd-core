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
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Errors
	errCreateClusterRoleBinding = "create clusterrolebinding error"
	errUpdateClusterRoleBinding = "update clusterrolebinding error"
	errDeleteClusterRoleBinding = "delete clusterrolebinding error"
	errGetClusterRoleBinding    = "get clusterrolebinding error"

	// items
	kindClusterRole = "ClusterRole"
	clusterRoleName = "ndd-manager-role"
)

func (e *APIEstablisher) createClusterRoleBinding(ctx context.Context, name string) error {
	log := e.log.WithValues("createClusterRoleBindingName", name)
	log.Debug("creating ClusterRoleBinding...")

	subjects := make([]rbacv1.Subject, 0)
	subjects = append(subjects, rbacv1.Subject{
		Kind:      rbacv1.ServiceAccountKind,
		Namespace: e.namespace,
		Name:      name,
	})

	b := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: e.namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     kindClusterRole,
			Name:     clusterRoleName,
		},
		Subjects: subjects,
	}
	err := e.client.Create(ctx, b)
	if err != nil {
		return errors.Wrap(err, errCreateClusterRoleBinding)
	}
	log.Debug("created ClusterRoleBinding...")
	return nil
}

func (e *APIEstablisher) deleteClusterRoleBinding(ctx context.Context, name string) error {
	log := e.log.WithValues("deleteClusterRoleBindingName", name)
	log.Debug("deleting ClusterRoleBinding...")

	b := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: e.namespace,
		},
	}
	err := e.client.Delete(ctx, b)
	if err != nil {
		return errors.Wrap(err, errDeleteClusterRoleBinding)
	}
	log.Debug("deleted ClusterRoleBinding...")
	return nil
}

func (e *APIEstablisher) updateClusterRoleBinding(ctx context.Context, name string) error {
	log := e.log.WithValues("updateClusterRoleBindingName", name)
	log.Debug("updating ClusterRoleBinding...")

	subjects := make([]rbacv1.Subject, 0)
	subjects = append(subjects, rbacv1.Subject{
		Kind:      rbacv1.ServiceAccountKind,
		Namespace: e.namespace,
		Name:      name,
	})

	b := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: e.namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     kindClusterRole,
			Name:     clusterRoleName,
		},
		Subjects: subjects,
	}

	err := e.client.Update(ctx, b)
	if err != nil {
		return errors.Wrap(err, errUpdateClusterRoleBinding)
	}
	log.Debug("updated ClusterRoleBinding...")

	return nil
}

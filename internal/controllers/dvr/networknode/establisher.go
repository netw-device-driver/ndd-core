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

	corev1 "k8s.io/api/core/v1"

	"github.com/netw-device-driver/ndd-runtime/pkg/logging"
	"github.com/netw-device-driver/ndd-runtime/pkg/resource"
)

// An Establisher brings up or down a set of resources in the API server
type Establisher interface {
	// Create brings up resources
	Create(ctx context.Context, name string, port int, c *corev1.Container) error

	// Update brings up resources
	Update(ctx context.Context, name string, port int, c *corev1.Container) error

	// Delete brings down resources
	Delete(ctx context.Context, name string) error
}

// APIEstablisher establishes control or ownership of resources in the API
// server for a parent.
type APIEstablisher struct {
	client    resource.ClientApplicator
	log       logging.Logger
	namespace string
}

// NewAPIEstablisher creates a new APIEstablisher.
func NewAPIEstablisher(client resource.ClientApplicator, log logging.Logger, namespace string) *APIEstablisher {
	return &APIEstablisher{
		client:    client,
		log:       log,
		namespace: namespace,
	}
}

// Create creates a set of resources
func (e *APIEstablisher) Create(ctx context.Context, name string, port int, c *corev1.Container) error { // nolint:gocyclo
	if err := e.createConfigMap(ctx, name); err != nil {
		return err
	}
	if err := e.createService(ctx, name, port); err != nil {
		return err
	}
	if err := e.createDeployment(ctx, name, c); err != nil {
		return err
	}
	return nil
}

// Update updates a set of resources
func (e *APIEstablisher) Update(ctx context.Context, name string, port int, c *corev1.Container) error { // nolint:gocyclo
	if err := e.updateConfigMap(ctx, name); err != nil {
		return err
	}
	if err := e.updateService(ctx, name, port); err != nil {
		return err
	}
	if err := e.updateDeployment(ctx, name, c); err != nil {
		return err
	}
	return nil
}

// Delete deletes a set of resources
func (e *APIEstablisher) Delete(ctx context.Context, name string) error { // nolint:gocyclo
	if err := e.deleteConfigMap(ctx, name); err != nil {
		return err
	}
	if err := e.deleteService(ctx, name); err != nil {
		return err
	}
	if err := e.deleteDeployment(ctx, name); err != nil {
		return err
	}
	return nil
}

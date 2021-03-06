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

package nn

import (
	"context"

	"github.com/netw-device-driver/ndd-runtime/pkg/logging"
	"github.com/netw-device-driver/ndd-runtime/pkg/resource"
	corev1 "k8s.io/api/core/v1"
)

// A Creds validates the validaty of various resources.
type Validator interface {
	// Validates the credentials
	ValidateCredentials(ctx context.Context, namespace, credentialsName, targetAddress string) (*Credentials, error)

	// Validates the device driver
	ValidateDeviceDriver(ctx context.Context, namespace, name, kind string, port int) (*corev1.Container, error)
}

type NnValidator struct {
	client resource.ClientApplicator
	log    logging.Logger
}

func NewNnValidator(client resource.ClientApplicator, log logging.Logger) *NnValidator {
	return &NnValidator{
		client: client,
		log:    log,
	}
}

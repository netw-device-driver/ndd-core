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

package pkg

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/netw-device-driver/ndd-core/internal/controllers/pkg/manager"
	"github.com/netw-device-driver/ndd-core/internal/controllers/pkg/resolver"
	"github.com/netw-device-driver/ndd-core/internal/controllers/pkg/revision"
	"github.com/netw-device-driver/ndd-core/internal/nddpkg"
	"github.com/netw-device-driver/ndd-runtime/pkg/logging"
)

// Setup package controllers.
func Setup(mgr ctrl.Manager, l logging.Logger, c nddpkg.Cache, namespace string) error {
	for _, setup := range []func(ctrl.Manager, logging.Logger, string) error{
		manager.SetupProvider,
		resolver.Setup,
	} {
		if err := setup(mgr, l, namespace); err != nil {
			return err
		}
	}
	for _, setup := range []func(ctrl.Manager, logging.Logger, nddpkg.Cache, string) error{
		revision.SetupProviderRevision,
	} {
		if err := setup(mgr, l, c, namespace); err != nil {
			return err
		}
	}
	return nil
}

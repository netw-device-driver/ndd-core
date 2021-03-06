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

package initializer

import (
	"context"

	v1 "github.com/netw-device-driver/ndd-core/apis/pkg/v1"
	"github.com/netw-device-driver/ndd-runtime/pkg/resource"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	errApplyLock = "cannot apply lock object"
)

// NewLockObject returns a new *LockObject initializer.
func NewLockObject() *LockObject {
	return &LockObject{}
}

// LockObject has the initializer for creating the Lock object.
type LockObject struct{}

// Run makes sure Lock object exists.
func (lo *LockObject) Run(ctx context.Context, kube client.Client) error {
	l := &v1.Lock{
		ObjectMeta: metav1.ObjectMeta{
			Name: "lock",
		},
	}
	return errors.Wrap(resource.NewAPIPatchingApplicator(kube).Apply(ctx, l), errApplyLock)
}

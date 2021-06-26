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

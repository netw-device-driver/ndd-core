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

package manager

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1 "github.com/netw-device-driver/ndd-core/apis/pkg/v1"
	"github.com/netw-device-driver/ndd-runtime/pkg/event"
	"github.com/netw-device-driver/ndd-runtime/pkg/logging"
	"github.com/netw-device-driver/ndd-runtime/pkg/resource"
)

const (
	parentLabel      = "pkg.ndd.henderiw.be/pakage"
	reconcileTimeout = 1 * time.Minute

	shortWait     = 30 * time.Second
	veryShortWait = 5 * time.Second
	pullWait      = 1 * time.Minute
)

func pullBasedRequeue(p *corev1.PullPolicy) reconcile.Result {
	r := reconcile.Result{}
	if p != nil && *p == corev1.PullAlways {
		r.RequeueAfter = pullWait
	}
	return r
}

const (
	errGetPackage           = "cannot get package"
	errListRevisions        = "cannot list revisions for package"
	errUnpack               = "cannot unpack package"
	errApplyPackageRevision = "cannot apply package revision"
	errGCPackageRevision    = "cannot garbage collect old package revision"

	errUpdateStatus                  = "cannot update package status"
	errUpdateInactivePackageRevision = "cannot update inactive package revision"

	errUnhealthyPackageRevision     = "current package revision is unhealthy"
	errUnknownPackageRevisionHealth = "current package revision health is unknown"
)

// Reconciler reconciles packages.
type Reconciler struct {
	client resource.ClientApplicator
	//pkg    Revisioner
	log    logging.Logger
	record event.Recorder

	newPackage func() v1.Package
	//newPackageRevision     func() v1.PackageRevision
	//newPackageRevisionList func() v1.PackageRevisionList
}

// ReconcilerOption is used to configure the Reconciler.
type ReconcilerOption func(*Reconciler)

// WithNewPackageFn determines the type of package being reconciled.
func WithNewPackageFn(f func() v1.Package) ReconcilerOption {
	return func(r *Reconciler) {
		r.newPackage = f
	}
}

// WithLogger specifies how the Reconciler should log messages.
func WithLogger(log logging.Logger) ReconcilerOption {
	return func(r *Reconciler) {
		r.log = log
	}
}

// WithRecorder specifies how the Reconciler should record Kubernetes events.
func WithRecorder(er event.Recorder) ReconcilerOption {
	return func(r *Reconciler) {
		r.record = er
	}
}

// SetupProvider adds a controller that reconciles Providers.
func SetupProvider(mgr ctrl.Manager, l logging.Logger, namespace string) error {
	name := "packages/" + strings.ToLower(v1.ProviderGroupKind)
	p := func() v1.Package { return &v1.Provider{} }
	//r := func() v1.PackageRevision { return &v1.ProviderRevision{} }
	//rl := func() v1.PackageRevisionList { return &v1.ProviderRevisionList{} }

	//clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	//if err != nil {
	//	return errors.Wrap(err, "failed to initialize clientset")
	//}

	r := NewReconciler(mgr,
		WithNewPackageFn(p),
		//WithNewPackageRevisionFn(r),
		//WithNewPackageRevisionListFn(rl),
		//WithRevisioner(NewPackageRevisioner(nddpkg.NewK8sFetcher(clientset, namespace))),
		WithLogger(l.WithValues("controller", name)),
		WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
	)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1.Provider{}).
		//Owns(&v1.ProviderRevision{}).
		Complete(r)
}

// NewReconciler creates a new package reconciler.
func NewReconciler(mgr ctrl.Manager, opts ...ReconcilerOption) *Reconciler {
	r := &Reconciler{
		client: resource.ClientApplicator{
			Client:     mgr.GetClient(),
			Applicator: resource.NewAPIPatchingApplicator(mgr.GetClient()),
		},
		//pkg:    NewNopRevisioner(),
		log:    logging.NewNopLogger(),
		record: event.NewNopRecorder(),
	}

	for _, f := range opts {
		f(r)
	}

	return r
}

// Reconcile package.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) { // nolint:gocyclo
	log := r.log.WithValues("request", req)
	log.Debug("Reconciling")

	ctx, cancel := context.WithTimeout(ctx, reconcileTimeout)
	defer cancel()

	p := r.newPackage()
	if err := r.client.Get(ctx, req.NamespacedName, p); err != nil {
		// There's no need to requeue if we no longer exist. Otherwise we'll be
		// requeued implicitly because we return an error.
		log.Debug(errGetPackage, "error", err)
		return reconcile.Result{}, errors.Wrap(resource.IgnoreNotFound(err), errGetPackage)
	}

	log = log.WithValues(
		"uid", p.GetUID(),
		"version", p.GetResourceVersion(),
		"name", p.GetName(),
	)

	// NOTE(hasheddan): when the first package revision is created for a
	// package, the health of the package is not set until the revision reports
	// its health. If updating from an existing revision, the package health
	// will match the health of the old revision until the next reconcile.
	return pullBasedRequeue(p.GetPackagePullPolicy()), errors.Wrap(r.client.Status().Update(ctx, p), errUpdateStatus)
}

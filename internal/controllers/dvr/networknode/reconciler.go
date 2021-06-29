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
	"fmt"
	"strings"
	"time"

	dvrv1 "github.com/netw-device-driver/ndd-core/apis/dvr/v1"
	"github.com/netw-device-driver/ndd-runtime/pkg/event"
	"github.com/netw-device-driver/ndd-runtime/pkg/logging"
	"github.com/netw-device-driver/ndd-runtime/pkg/meta"
	"github.com/netw-device-driver/ndd-runtime/pkg/resource"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	// Finalizer
	finalizer = "networknode.dvr.ndd.henderiw.be"

	// default
	defaultGrpcPort = 9999

	// Timers
	reconcileTimeout = 1 * time.Minute
	shortWait        = 30 * time.Second
	veryShortWait    = 5 * time.Second

	// Errors
	errGetNetworkNode = "cannot get network node resource"
	errUpdateStatus   = "cannot update network node status"

	errAddFinalizer    = "cannot add network node finalizer"
	errRemoveFinalizer = "cannot remove network node finalizer"

	errCredentials = "invalid credentials"

	errDeleteObjects = "cannot delete configmap, servide or deployment"
	errCreateObjects = "cannot create configmap, servide or deployment"

	// Event reasons
	reasonSync event.Reason = "SyncNetworkNode"
)

// ReconcilerOption is used to configure the Reconciler.
type ReconcilerOption func(*Reconciler)

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

// WithEstablisher specifies how the reconciler should create/delete/update
// resources through the k8s api
func WithEstablisher(er *APIEstablisher) ReconcilerOption {
	return func(r *Reconciler) {
		r.objects = er
	}
}

// WithValidator specifies how the Reconciler should perform object
// validation.
func WithValidator(v Validator) ReconcilerOption {
	return func(r *Reconciler) {
		r.validator = v
	}
}

// Reconciler reconciles packages.
type Reconciler struct {
	client      client.Client
	nnFinalizer resource.Finalizer
	log         logging.Logger
	record      event.Recorder
	objects     Establisher
	validator   Validator
}

// Setup adds a controller that reconciles the Lock.
func Setup(mgr ctrl.Manager, l logging.Logger, namespace string) error {
	name := "dvr/" + strings.ToLower(dvrv1.NetworkNodeKind)

	r := NewReconciler(mgr,
		WithLogger(l.WithValues("controller", name)),
		WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		WithEstablisher(NewAPIEstablisher(resource.ClientApplicator{
			Client:     mgr.GetClient(),
			Applicator: resource.NewAPIPatchingApplicator(mgr.GetClient()),
		}, l, namespace)),
		WithValidator(NewNnValidator(resource.ClientApplicator{
			Client:     mgr.GetClient(),
			Applicator: resource.NewAPIPatchingApplicator(mgr.GetClient()),
		}, l)),
	)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&dvrv1.NetworkNode{}).
		Watches(
			&source.Kind{Type: &dvrv1.DeviceDriver{}},
			handler.EnqueueRequestsFromMapFunc(r.DeviceDriverMapFunc),
		).
		WithEventFilter(resource.IgnoreUpdateWithoutGenerationChangePredicate()).
		Complete(r)
}

// NewReconciler creates a new package revision reconciler.
func NewReconciler(mgr manager.Manager, opts ...ReconcilerOption) *Reconciler {
	r := &Reconciler{
		client:      mgr.GetClient(),
		nnFinalizer: resource.NewAPIFinalizer(mgr.GetClient(), finalizer),
		log:         logging.NewNopLogger(),
		record:      event.NewNopRecorder(),
	}

	for _, f := range opts {
		f(r)
	}

	return r
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;delete
// +kubebuilder:rbac:groups="",resources=pods,verbs=list;watch;get;patch;create;update;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=list;watch;get
// +kubebuilder:rbac:groups="",resources=events,verbs=list;watch;get;patch;create;update;delete
// +kubebuilder:rbac:groups=dvr.ndd.henderiw.be,resources=devicedrivers,verbs=get;list;watch
// +kubebuilder:rbac:groups=dvr.ndd.henderiw.be,resources=networknodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dvr.ndd.henderiw.be,resources=networknodes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dvr.ndd.henderiw.be,resources=networknodes/finalizers,verbs=update

// Reconcile network node.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) { // nolint:gocyclo
	log := r.log.WithValues("request", req)
	log.Debug("Network Node", "NameSpace", req.NamespacedName)

	nn := &dvrv1.NetworkNode{}
	if err := r.client.Get(ctx, req.NamespacedName, nn); err != nil {
		// There's no need to requeue if we no longer exist. Otherwise we'll be
		// requeued implicitly because we return an error.
		log.Debug(errGetNetworkNode, "error", err)
		return reconcile.Result{}, errors.Wrap(resource.IgnoreNotFound(err), errGetNetworkNode)
	}
	log.Debug("network Node", "NN", nn)
	log.Debug("Health status", "status", nn.GetCondition(dvrv1.ConditionKindDeviceDriverHealthy).Status)

	if meta.WasDeleted(nn) {
		// remove delete the configmap, service, deployment when the service was healthy
		if nn.GetCondition(dvrv1.ConditionKindDeviceDriverHealthy).Status == corev1.ConditionTrue {
			if err := r.objects.Delete(ctx, nn.Name); err != nil {
				log.Debug(errDeleteObjects, "error", err)
				r.record.Event(nn, event.Warning(reasonSync, errors.Wrap(err, errDeleteObjects)))
				nn.SetConditions(dvrv1.Unhealthy(), dvrv1.Inactive(), dvrv1.NotDiscovered())
				return reconcile.Result{RequeueAfter: shortWait}, errors.Wrap(r.client.Status().Update(ctx, nn), errUpdateStatus)
			}
		}
		// Delete finalizer after the object is deleted
		if err := r.nnFinalizer.RemoveFinalizer(ctx, nn); err != nil {
			log.Debug(errRemoveFinalizer, "error", err)
			r.record.Event(nn, event.Warning(reasonSync, errors.Wrap(err, errRemoveFinalizer)))
			return reconcile.Result{RequeueAfter: shortWait}, nil
		}
		return reconcile.Result{Requeue: false}, nil
	}

	// Add a finalizer to newly created objects and update the conditions
	if err := r.nnFinalizer.AddFinalizer(ctx, nn); err != nil {
		log.Debug(errAddFinalizer, "error", err)
		r.record.Event(nn, event.Warning(reasonSync, errors.Wrap(err, errAddFinalizer)))
		return reconcile.Result{RequeueAfter: shortWait}, nil
	}

	// Retrieve the Login details from the network node spec and validate
	// the network node details and build the credentials for communicating
	// to the network node.
	if nn.Spec.Target.CredentialsName != nil && nn.Spec.Target.Address != nil {

	}
	creds, err := r.validator.ValidateCredentials(ctx, nn.Namespace, *nn.Spec.Target.CredentialsName, *nn.Spec.Target.Address)
	log.Debug("Network node creds", "creds", creds, "err", err)
	if err != nil || creds == nil {
		// remove delete the configmap, service, deployment when the service was healthy
		if nn.GetCondition(dvrv1.ConditionKindDeviceDriverHealthy).Status == corev1.ConditionTrue {
			if err := r.objects.Delete(ctx, nn.Name); err != nil {
				log.Debug(errDeleteObjects, "error", err)
				r.record.Event(nn, event.Warning(reasonSync, errors.Wrap(err, errDeleteObjects)))
				nn.SetConditions(dvrv1.Unhealthy(), dvrv1.Inactive(), dvrv1.NotDiscovered())
				return reconcile.Result{RequeueAfter: shortWait}, errors.Wrap(r.client.Status().Update(ctx, nn), errUpdateStatus)
			}
		}
		log.Debug(errCredentials, "error", err)
		r.record.Event(nn, event.Warning(reasonSync, errors.Wrap(err, errCredentials)))
		nn.SetConditions(dvrv1.Unhealthy(), dvrv1.Inactive(), dvrv1.NotDiscovered())
		return reconcile.Result{RequeueAfter: shortWait}, errors.Wrap(r.client.Status().Update(ctx, nn), errUpdateStatus)
	}

	// validate device driver information
	// NOTE: the parameters are required in the api and will get defaults if not specified, so we dont have to add validation
	// if they exist in the api or not
	c, err := r.validator.ValidateDeviceDriver(ctx, nn.Namespace, nn.Name, string(*nn.Spec.DeviceDriverKind), *nn.Spec.GrpcServerPort)
	log.Debug("Validate device driver", "containerInfo", c, "err", err)
	if err != nil || c == nil {
		if nn.GetCondition(dvrv1.ConditionKindDeviceDriverHealthy).Status == corev1.ConditionTrue {
			if err := r.objects.Delete(ctx, nn.Name); err != nil {
				log.Debug(errDeleteObjects, "error", err)
				r.record.Event(nn, event.Warning(reasonSync, errors.Wrap(err, errDeleteObjects)))
				nn.SetConditions(dvrv1.Unhealthy(), dvrv1.Inactive(), dvrv1.NotDiscovered())
				return reconcile.Result{RequeueAfter: shortWait}, errors.Wrap(r.client.Status().Update(ctx, nn), errUpdateStatus)
			}
		}
	}

	log.Debug("Health status", "status", nn.GetCondition(dvrv1.ConditionKindDeviceDriverHealthy).Status)
	if nn.GetCondition(dvrv1.ConditionKindDeviceDriverHealthy).Status == corev1.ConditionFalse ||
		nn.GetCondition(dvrv1.ConditionKindDeviceDriverHealthy).Status == corev1.ConditionUnknown {
		if err := r.objects.Create(ctx, nn.Name, *nn.Spec.GrpcServerPort, c); err != nil {
			log.Debug(errCreateObjects, "error", err)
			r.record.Event(nn, event.Warning(reasonSync, errors.Wrap(err, errCreateObjects)))
			nn.SetConditions(dvrv1.Unhealthy(), dvrv1.Inactive(), dvrv1.NotDiscovered())
			return reconcile.Result{RequeueAfter: shortWait}, errors.Wrap(r.client.Status().Update(ctx, nn), errUpdateStatus)
		}
		nn.SetConditions(dvrv1.Healthy(), dvrv1.Active(), dvrv1.NotDiscovered())
		return reconcile.Result{}, errors.Wrap(r.client.Status().Update(ctx, nn), errUpdateStatus)
	}
	return reconcile.Result{}, nil
}

// DeviceDriverMapFunc is a handler.ToRequestsFunc to be used to enqeue
// request for reconciliation of NetworkNode when the device driver object changes
func (r *Reconciler) DeviceDriverMapFunc(o client.Object) []ctrl.Request {
	log := r.log.WithValues("deviceDriver Object", o)
	result := []ctrl.Request{}

	dd, ok := o.(*dvrv1.DeviceDriver)
	if !ok {
		panic(fmt.Sprintf("Expected a DeviceDriver but got a %T", o))
	}
	log.WithValues(dd.GetName(), dd.GetNamespace()).Info("DeviceDriver MapFunction")

	selectors := []client.ListOption{
		client.InNamespace(dd.Namespace),
		client.MatchingLabels{},
	}
	nn := &dvrv1.NetworkNodeList{}
	if err := r.client.List(context.TODO(), nn, selectors...); err != nil {
		return result
	}

	for _, n := range nn.Items {
		// only enqueue if the network node device driver kind matches with the device driver label
		for k, v := range dd.GetLabels() {
			if k == "ddriver-kind" {
				if string(*n.Spec.DeviceDriverKind) == v {
					name := client.ObjectKey{
						Namespace: n.GetNamespace(),
						Name:      n.GetName(),
					}
					log.WithValues(n.GetName(), n.GetNamespace()).Info("DeviceDriverMapFunc networknode ReQueue")
					result = append(result, ctrl.Request{NamespacedName: name})
				}
			}
		}
	}
	return result
}

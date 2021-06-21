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

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/netw-device-driver/ndd-core/internal/controllers/pkg"
	"github.com/netw-device-driver/ndd-core/internal/nddpkg"
	"github.com/netw-device-driver/ndd-runtime/pkg/logging"
	//+kubebuilder:scaffold:imports
)

var (
	metricsAddr          string
	probeAddr            string
	enableLeaderElection bool
	concurrency          int
	namespace            string
	cacheDir             string
)

// startCmd represents the start command for the network device driver
var startCmd = &cobra.Command{
	Use:          "start",
	Short:        "start the network device driver core",
	Long:         "start the network device driver core",
	Aliases:      []string{"start"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		zlog := zap.New(zap.UseDevMode(debug), zap.JSONEncoder())
		if debug {
			// Only use a logr.Logger when debug is on
			ctrl.SetLogger(zlog)
		}
		zlog.Info("create manager")
		mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
			Scheme: scheme,
			//MetricsBindAddress:     metricsAddr,
			//Port:                   9443,
			HealthProbeBindAddress: probeAddr,
			LeaderElection:         false,
			//LeaderElection:         enableLeaderElection,
			LeaderElectionID: "c66ce353.ndd.henderiw.be",
		})
		if err != nil {
			return errors.Wrap(err, "Cannot create manager")
		}

		//log.Info("setup reconcilers/controllers")
		//ctx := ctrl.SetupSignalHandler()
		//setupReconcilers(ctx, mgr)

		pkgCache := nddpkg.NewImageCache(cacheDir, afero.NewOsFs())
		zlog.Info("Cache Directory", "cacheDir", cacheDir)

		if err := pkg.Setup(mgr, logging.NewLogrLogger(zlog.WithName("nddcore-pkg")), pkgCache, namespace); err != nil {
			return errors.Wrap(err, "Cannot add packages controllers to manager")
		}

		// +kubebuilder:scaffold:builder

		if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
			return errors.Wrap(err, "unable to set up health check")
		}
		if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
			return errors.Wrap(err, "unable to set up ready check")
		}

		zlog.Info("starting manager")
		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			return errors.Wrap(err, "problem running manager")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&metricsAddr, "metrics-bind-address", "m", ":8080", "The address the metric endpoint binds to.")
	startCmd.Flags().StringVarP(&probeAddr, "health-probe-bind-address", "p", ":8081", "The address the probe endpoint binds to.")
	startCmd.Flags().BoolVarP(&enableLeaderElection, "leader-elect", "l", false, "Enable leader election for controller manager. "+
		"Enabling this will ensure there is only one active controller manager.")
	startCmd.Flags().IntVarP(&concurrency, "concurrency", "", 1, "Number of items to process simultaneously")
	startCmd.Flags().StringVarP(&namespace, "namespace", "n", viper.GetString("POD_NAMESPACE"), "Namespace used to unpack and run packages.")
	startCmd.Flags().StringVarP(&cacheDir, "cache-dir", "c", "/cache", "Directory used for caching package images.")

}

/*
func setupReconcilers(ctx context.Context, mgr ctrl.Manager) {
	if err := (&drivercontrollers.NetworkNodeReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		errors.Wrap(err, "unable to create controller NetworkNode")
	}
	if err := (&drivercontrollers.DeviceDriverReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		errors.Wrap(err, "unable to create controller DeviceDriver")
	}
	if err := (&pkgcontrollers.ProviderReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		errors.Wrap(err, "unable to create controller Provider")
		os.Exit(1)
	}
}
*/

func nddConcurrency(c int) controller.Options {
	return controller.Options{MaxConcurrentReconciles: c}
}

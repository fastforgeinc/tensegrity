/*
Copyright 2024 FastForge Inc. support@fastforge.com.

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

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"os"
	"time"

	"reconciler.io/runtime/reconcilers"
	"reconciler.io/runtime/tracker"

	rolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"

	apiargov1alpha1 "github.com/fastforgeinc/tensegrity/api/argo/v1alpha1"
	apik8sv1alpha1 "github.com/fastforgeinc/tensegrity/api/k8s/v1alpha1"
	apiv1alpha1 "github.com/fastforgeinc/tensegrity/api/v1alpha1"
	controllerargov1alpha1 "github.com/fastforgeinc/tensegrity/internal/controller/argo/v1alpha1"
	controllerk8sv1alpha1 "github.com/fastforgeinc/tensegrity/internal/controller/k8s/v1alpha1"
	controllerv1alpha1 "github.com/fastforgeinc/tensegrity/internal/controller/v1alpha1"

	"github.com/fastforgeinc/tensegrity/pkg/client"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(apiv1alpha1.AddToScheme(scheme))
	utilruntime.Must(apik8sv1alpha1.AddToScheme(scheme))
	utilruntime.Must(apiargov1alpha1.AddToScheme(scheme))
	utilruntime.Must(rolloutsv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var secureMetrics bool
	var enableHTTP2 bool
	var certDir string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.StringVar(&certDir, "cert-dir", "", "The directory that contains the server key and certificate.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&secureMetrics, "metrics-secure", false,
		"If set the metrics endpoint is served securely")
	flag.BoolVar(&enableHTTP2, "enable-http2", false,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// if the enable-http2 flag is false (the default), http/2 should be disabled
	// due to its vulnerabilities. More specifically, disabling http/2 will
	// prevent from being vulnerable to the HTTP/2 Stream Cancellation and
	// Rapid Reset CVEs. For more information see:
	// - https://github.com/advisories/GHSA-qppj-fm5r-hxr3
	// - https://github.com/advisories/GHSA-4374-p667-p6c8
	disableHTTP2 := func(c *tls.Config) {
		setupLog.Info("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}

	var tlsOpts []func(*tls.Config)
	if !enableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}

	webhookOptions := webhook.Options{
		TLSOpts: tlsOpts,
	}
	if len(certDir) > 0 {
		webhookOptions.CertDir = certDir
	}

	webhookServer := webhook.NewServer(webhookOptions)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:    scheme,
		NewClient: client.New,
		Metrics: metricsserver.Options{
			BindAddress:   metricsAddr,
			SecureServing: secureMetrics,
			TLSOpts:       tlsOpts,
		},
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "83f10920.tensegrity.fastforge.io",
		//LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	reconcilerConfig := &reconcilers.Config{
		Client:    mgr.GetClient(),
		APIReader: mgr.GetAPIReader(),
		Recorder:  mgr.GetEventRecorderFor("tensegrity"),
		Tracker:   tracker.New(scheme, 1*time.Hour),
	}

	ctx := context.Background()
	consumerReconciler := controllerv1alpha1.NewConsumerReconciler()
	consumerSecretReconciler := controllerv1alpha1.NewConsumerSecretReconciler()
	consumerConfigMapReconciler := controllerv1alpha1.NewConsumerConfigMapReconciler()
	producerReconciler := controllerv1alpha1.NewProducerReconciler()
	producerSecretReconciler := controllerv1alpha1.NewProducerSecretReconciler()
	producerConfigMapReconciler := controllerv1alpha1.NewProducerConfigMapReconciler()

	if err = controllerk8sv1alpha1.NewDeploymentReconciler(
		reconcilerConfig,
		consumerReconciler,
		consumerSecretReconciler,
		consumerConfigMapReconciler,
		producerReconciler,
		producerSecretReconciler,
		producerConfigMapReconciler).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Deployment", "version", "k8s/v1alpha1")
		os.Exit(1)
	}
	if err = controllerk8sv1alpha1.NewStatefulSetReconciler(
		reconcilerConfig,
		consumerReconciler,
		consumerSecretReconciler,
		consumerConfigMapReconciler,
		producerReconciler,
		producerSecretReconciler,
		producerConfigMapReconciler).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "StatefulSet", "version", "k8s/v1alpha1")
		os.Exit(1)
	}
	if err = controllerk8sv1alpha1.NewDaemonSetReconciler(
		reconcilerConfig,
		consumerReconciler,
		consumerSecretReconciler,
		consumerConfigMapReconciler,
		producerReconciler,
		producerSecretReconciler,
		producerConfigMapReconciler).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DaemonSet", "version", "k8s/v1alpha1")
		os.Exit(1)
	}
	if err = controllerargov1alpha1.NewRolloutReconciler(
		reconcilerConfig,
		consumerReconciler,
		consumerSecretReconciler,
		consumerConfigMapReconciler,
		producerReconciler,
		producerSecretReconciler,
		producerConfigMapReconciler).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Rollout", "version", "argo/v1alpha1")
		os.Exit(1)
	}
	if err = controllerv1alpha1.NewStaticReconciler(
		reconcilerConfig,
		producerReconciler,
		producerSecretReconciler,
		producerConfigMapReconciler).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Static", "version", "v1alpha1")
		os.Exit(1)
	}
	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		if err = new(apik8sv1alpha1.Deployment).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Deployment")
			os.Exit(1)
		}
		if err = new(apik8sv1alpha1.StatefulSet).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "StatefulSet")
			os.Exit(1)
		}
		if err = new(apik8sv1alpha1.DaemonSet).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "DaemonSet")
			os.Exit(1)
		}
		if err = new(apiargov1alpha1.Rollout).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Rollout")
			os.Exit(1)
		}
		if err = new(apiv1alpha1.Static).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Static")
			os.Exit(1)
		}
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

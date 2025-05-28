/*
This file is part of the Tensegrity distribution (https://github.com/fastforgeinc/tensegrity)
Copyright (C) 2024 FastForge, Inc.

Tensegrity is free software: you can redistribute it and/or modify it
under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License,
or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with
this program. If not, see http://www.gnu.org/licenses/.
*/

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"os"
	"time"

	"reconciler.io/runtime/reconcilers"

	apik8sv1alpha1 "github.com/fastforgeinc/tensegrity/api/k8s/v1alpha1"
	apiv1alpha1 "github.com/fastforgeinc/tensegrity/api/v1alpha1"
	controllerk8sv1alpha1 "github.com/fastforgeinc/tensegrity/internal/controller/k8s/v1alpha1"
	controllerv1alpha1 "github.com/fastforgeinc/tensegrity/internal/controller/v1alpha1"

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
	scheme     = runtime.NewScheme()
	setupLog   = ctrl.Log.WithName("setup")
	syncPeriod = 1 * time.Hour
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(apiv1alpha1.AddToScheme(scheme))
	utilruntime.Must(apik8sv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var secureMetrics bool
	var enableHTTP2 bool
	var enableWebhooks bool
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
	flag.BoolVar(&enableWebhooks, "enable-webhooks", false, "If set, webhook validation will be enabled")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	disableHTTP2 := func(c *tls.Config) {
		setupLog.Info("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}

	var tlsOpts []func(*tls.Config)
	if !enableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}

	mgrOptions := ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress:   metricsAddr,
			SecureServing: secureMetrics,
			TLSOpts:       tlsOpts,
		},
		HealthProbeBindAddress:        probeAddr,
		LeaderElection:                enableLeaderElection,
		LeaderElectionID:              "83f10920.tensegrity.fastforge.io",
		LeaderElectionReleaseOnCancel: true,
	}

	if enableWebhooks {
		webhookOptions := webhook.Options{
			TLSOpts: tlsOpts,
		}
		if len(certDir) > 0 {
			webhookOptions.CertDir = certDir
		}

		webhookServer := webhook.NewServer(webhookOptions)
		mgrOptions.WebhookServer = webhookServer
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), mgrOptions)
	if err != nil {
		setupLog.Error(err, "unable to start controller manager")
		os.Exit(1)
	}

	ctx := context.Background()
	config := reconcilers.NewConfig(mgr, nil, syncPeriod)
	validationReconciler := controllerv1alpha1.NewValidationReconciler()
	consumerReconciler := controllerv1alpha1.NewConsumerReconciler()
	consumerSecretReconciler := controllerv1alpha1.NewConsumerSecretReconciler()
	consumerConfigMapReconciler := controllerv1alpha1.NewConsumerConfigMapReconciler()
	producerReconciler := controllerv1alpha1.NewProducerReconciler()
	producerSecretReconciler := controllerv1alpha1.NewProducerSecretReconciler()
	producerConfigMapReconciler := controllerv1alpha1.NewProducerConfigMapReconciler()

	if err = controllerk8sv1alpha1.NewDeploymentReconciler(
		&config,
		validationReconciler,
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
		&config,
		validationReconciler,
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
		&config,
		validationReconciler,
		consumerReconciler,
		consumerSecretReconciler,
		consumerConfigMapReconciler,
		producerReconciler,
		producerSecretReconciler,
		producerConfigMapReconciler).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DaemonSet", "version", "k8s/v1alpha1")
		os.Exit(1)
	}
	if err = controllerv1alpha1.NewStaticReconciler(
		&config,
		validationReconciler,
		producerReconciler,
		producerSecretReconciler,
		producerConfigMapReconciler).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Static", "version", "v1alpha1")
		os.Exit(1)
	}
	if enableWebhooks {
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

	setupLog.Info("starting controller manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running controller manager")
		os.Exit(1)
	}
}

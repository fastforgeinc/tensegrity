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

package v1alpha1

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"reconciler.io/runtime/reconcilers"
	"reconciler.io/runtime/tracker"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	k8sv1alpha1 "github.com/fastforgeinc/tensegrity/api/k8s/v1alpha1"
	controllerv1alpha1 "github.com/fastforgeinc/tensegrity/internal/controller/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var mgr manager.Manager
var k8sClient client.Client
var testEnv *envtest.Environment
var reconcilerConfig *reconcilers.Config
var validationReconciler *controllerv1alpha1.ValidationReconciler
var consumerReconciler *controllerv1alpha1.ConsumerReconciler
var consumerSecretReconciler *controllerv1alpha1.ConsumerSecretReconciler
var consumerConfigMapReconciler *controllerv1alpha1.ConsumerConfigMapReconciler
var producerReconcilerInstance *controllerv1alpha1.ProducerReconciler
var producerSecretReconcilerInstance *controllerv1alpha1.ProducerSecretReconciler
var producerConfigMapReconcilerInstance *controllerv1alpha1.ProducerConfigMapReconciler

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,

		// The BinaryAssetsDirectory is only required if you want to run the tests directly
		// without call the makefile target test. If not informed it will look for the
		// default path defined in controller-runtime which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform
		// the tests directly. When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join("..", "..", "..", "..", "bin", "k8s",
			fmt.Sprintf("1.29.0-%s-%s", runtime.GOOS, runtime.GOARCH)),
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = k8sv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	mgr, err = controllerruntime.NewManager(cfg, controllerruntime.Options{})
	Expect(err).NotTo(HaveOccurred())
	Expect(mgr).NotTo(BeNil())

	reconcilerConfig = &reconcilers.Config{
		Client:    k8sClient,
		APIReader: mgr.GetAPIReader(),
		Recorder:  mgr.GetEventRecorderFor("tensegrity"),
		Tracker:   tracker.New(scheme.Scheme, 1*time.Hour),
	}

	validationReconciler = controllerv1alpha1.NewValidationReconciler()
	consumerReconciler = controllerv1alpha1.NewConsumerReconciler()
	consumerSecretReconciler = controllerv1alpha1.NewConsumerSecretReconciler()
	consumerConfigMapReconciler = controllerv1alpha1.NewConsumerConfigMapReconciler()
	producerReconcilerInstance = controllerv1alpha1.NewProducerReconciler()
	producerSecretReconcilerInstance = controllerv1alpha1.NewProducerSecretReconciler()
	producerConfigMapReconcilerInstance = controllerv1alpha1.NewProducerConfigMapReconciler()
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

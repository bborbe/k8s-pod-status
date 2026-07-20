// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package factory

import (
	"context"
	"net/http"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	libk8s "github.com/bborbe/k8s"
	libtime "github.com/bborbe/time"
	"k8s.io/client-go/kubernetes"

	"github.com/bborbe/k8s-pod-status/pkg"
	"github.com/bborbe/k8s-pod-status/pkg/handler"
)

func CreateListStatusHandler(kubeconfig string, namespace libk8s.Namespace) http.Handler {
	return libhttp.NewErrorHandler(handler.NewListStatusHandler(kubeconfig, namespace))
}

func CreateListFailedHandler(kubeconfig string, namespace libk8s.Namespace) http.Handler {
	return libhttp.NewErrorHandler(handler.NewListFailedHandler(kubeconfig, namespace))
}

func CreatePodWatcher(kubeconfig string, namespace libk8s.Namespace) func(context.Context) error {
	return func(ctx context.Context) error {
		clientset, err := libk8s.CreateClientset(kubeconfig)
		if err != nil {
			return errors.Wrap(ctx, err, "create clientset failed")
		}
		return createPodWatcher(clientset, namespace).Watch(ctx)
	}
}

func createPodWatcher(
	clientset kubernetes.Interface,
	namespace libk8s.Namespace,
) libk8s.PodWatcher {
	return libk8s.NewPodWatcherRetry(
		libk8s.NewPodWatcher(
			clientset,
			libk8s.NewPodEventProcessorSkipError(
				pkg.NewPodEventProcessor(
					pkg.NewPodStatusMetrics(),
				),
			),
			namespace,
		),
		libtime.NewWaiterDuration(),
		1*libtime.Second,
	)
}

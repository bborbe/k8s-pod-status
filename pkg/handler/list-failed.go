// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	libk8s "github.com/bborbe/k8s"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/bborbe/k8s-pod-status/pkg"
)

func NewListFailedHandler(kubeconfig string, namespace libk8s.Namespace) libhttp.WithError {
	return libhttp.WithErrorFunc(
		func(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
			k8sClientset, err := libk8s.CreateClientset(kubeconfig)
			if err != nil {
				return errors.Wrap(ctx, err, "create k8sClientset failed")
			}
			list, err := k8sClientset.CoreV1().
				Pods(namespace.String()).
				List(ctx, metav1.ListOptions{})
			if err != nil {
				return errors.Wrap(ctx, err, "list failed")
			}
			fmt.Fprintf(resp, "Failed:\n")
			for _, pod := range list.Items {
				status := pkg.PodStatus(pod)
				if status == pkg.RunningStatus || status == pkg.SucceededStatus {
					continue
				}
				fmt.Fprintf(resp, "- Name: %s Status: %s\n", pod.Name, status)
			}
			glog.V(2).Infof("list completed")
			return nil
		},
	)
}

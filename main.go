// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"

	libhttp "github.com/bborbe/http"
	libk8s "github.com/bborbe/k8s"
	"github.com/bborbe/run"
	libsentry "github.com/bborbe/sentry"
	"github.com/bborbe/service"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bborbe/k8s-pod-status/pkg/factory"
	libfactory "github.com/bborbe/k8s-pod-status/pkg/libfactory"
	libmetrics "github.com/bborbe/k8s-pod-status/pkg/libmetrics"
)

func main() {
	app := &application{}
	os.Exit(service.Main(context.Background(), app, &app.SentryDSN, &app.SentryProxy))
}

type application struct {
	SentryDSN      string            `required:"true"  arg:"sentry-dsn"       env:"SENTRY_DSN"       usage:"SentryDSN"                 display:"length"`
	SentryProxy    string            `required:"false" arg:"sentry-proxy"     env:"SENTRY_PROXY"     usage:"Sentry Proxy"`
	Listen         string            `required:"true"  arg:"listen"           env:"LISTEN"           usage:"address to listen to"`
	Kubeconfig     string            `required:"false" arg:"kubeconfig"       env:"KUBECONFIG"       usage:"Path to k8s config"`
	Namespace      libk8s.Namespace  `required:"true"  arg:"namespace"        env:"NAMESPACE"        usage:"Kubernetes namespace"`
	BuildGitCommit string            `required:"false" arg:"build-git-commit" env:"BUILD_GIT_COMMIT" usage:"Build Git commit hash"                      default:"none"`
	BuildDate      *libtime.DateTime `required:"false" arg:"build-date"       env:"BUILD_DATE"       usage:"Build timestamp (RFC3339)"`
}

func (a *application) Run(ctx context.Context, sentryClient libsentry.Client) error {
	libmetrics.NewBuildInfoMetrics().SetBuildInfo(a.BuildDate)

	return service.Run(
		ctx,
		a.watchStatus(),
		a.createHTTPServer(),
	)
}

func (a *application) watchStatus() run.Func {
	return factory.CreatePodWatcher(a.Kubeconfig, a.Namespace)
}

func (a *application) createHTTPServer() run.Func {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		router := mux.NewRouter()
		router.Path("/healthz").Handler(libhttp.NewPrintHandler("OK"))
		router.Path("/readiness").Handler(libhttp.NewPrintHandler("OK"))
		router.Path("/metrics").Handler(promhttp.Handler())
		router.Path("/setloglevel/{level}").Handler(libfactory.CreateSetLoglevelHandler(ctx))
		router.Path("/status").Handler(factory.CreateListStatusHandler(a.Kubeconfig, a.Namespace))

		router.Path("/failed").Handler(factory.CreateListFailedHandler(a.Kubeconfig, a.Namespace))

		glog.V(2).Infof("starting http server listen on %s", a.Listen)
		return libhttp.NewServer(
			a.Listen,
			router,
		).Run(ctx)
	}
}

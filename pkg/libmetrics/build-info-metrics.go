// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metrics

import (
	libtime "github.com/bborbe/time"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	buildInfo = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "build_info",
			Help: "Build timestamp as Unix time. Service identified by Prometheus job label.",
		},
	)
)

func init() {
	prometheus.MustRegister(buildInfo)
}

//counterfeiter:generate -o ../mocks/build-info-metrics.go --fake-name BuildInfoMetrics . BuildInfoMetrics
type BuildInfoMetrics interface {
	SetBuildInfo(buildDate *libtime.DateTime)
}

func NewBuildInfoMetrics() BuildInfoMetrics {
	return &buildInfoMetrics{}
}

type buildInfoMetrics struct{}

func (m *buildInfoMetrics) SetBuildInfo(buildDate *libtime.DateTime) {
	if buildDate == nil {
		return
	}
	buildInfo.Set(float64(buildDate.Unix()))
}

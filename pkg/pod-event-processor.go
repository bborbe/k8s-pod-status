// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"context"

	libk8s "github.com/bborbe/k8s"
	corev1 "k8s.io/api/core/v1"
)

// NewPodEventProcessor creates a libk8s.PodEventProcessor that updates
// Prometheus metrics based on pod events.
func NewPodEventProcessor(metrics PodStatusMetrics) libk8s.PodEventProcessor {
	return &podEventProcessor{
		metrics: metrics,
	}
}

type podEventProcessor struct {
	metrics PodStatusMetrics
}

func (p *podEventProcessor) OnUpdate(ctx context.Context, pod corev1.Pod) error {
	p.metrics.Set(pod)
	return nil
}

func (p *podEventProcessor) OnDelete(ctx context.Context, pod corev1.Pod) error {
	p.metrics.Remove(pod)
	return nil
}

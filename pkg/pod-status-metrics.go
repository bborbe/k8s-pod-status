// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"time"

	libk8s "github.com/bborbe/k8s"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
)

//counterfeiter:generate -o ../mocks/pod-status-metrics.go --fake-name PodStatusMetrics . PodStatusMetrics
type PodStatusMetrics interface {
	Set(pod corev1.Pod)
	Remove(pod corev1.Pod)
}

var (
	counter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "k8s",
		Subsystem: "pod_status",
		Name:      "counter",
		Help:      "Counts processed messages",
	}, []string{"namespace", "name", "app", "status"})

	restartCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "k8s",
		Subsystem: "pod_status",
		Name:      "restart_count",
		Help:      "Current restart count of containers in pod",
	}, []string{"namespace", "name", "app", "container"})

	lastRestartTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "k8s",
		Subsystem: "pod_status",
		Name:      "last_restart_time",
		Help:      "Unix timestamp of last container restart",
	}, []string{"namespace", "name", "app", "container"})
)

func init() {
	prometheus.MustRegister(
		counter,
		restartCount,
		lastRestartTime,
	)
}

func NewPodStatusMetrics() PodStatusMetrics {
	return &podStatusMetrics{}
}

type podStatusMetrics struct {
}

func (m *podStatusMetrics) Set(pod corev1.Pod) {
	name := libk8s.NameFromPod(pod)
	status := PodStatus(pod)
	glog.V(2).Infof("add metrics for '%s' with phase '%s'", name, pod.Status.Phase)

	// Set status counter metrics
	for _, s := range AvailableStatus {
		if s == status {
			counter.With(prometheus.Labels{
				"name":      name.String(),
				"status":    s.String(),
				"app":       pod.Labels["app"],
				"namespace": pod.Namespace,
			}).Set(1)
		} else {
			counter.With(prometheus.Labels{
				"name":      name.String(),
				"status":    s.String(),
				"app":       pod.Labels["app"],
				"namespace": pod.Namespace,
			}).Set(0)
		}
	}

	// Set restart count metrics for each container
	for _, containerStatus := range pod.Status.ContainerStatuses {
		labels := prometheus.Labels{
			"name":      name.String(),
			"app":       pod.Labels["app"],
			"namespace": pod.Namespace,
			"container": containerStatus.Name,
		}

		restartCount.With(labels).Set(float64(containerStatus.RestartCount))

		// Track last restart time
		var lastRestart time.Time
		if containerStatus.LastTerminationState.Terminated != nil &&
			!containerStatus.LastTerminationState.Terminated.FinishedAt.IsZero() {
			// Use last termination time if available
			lastRestart = containerStatus.LastTerminationState.Terminated.FinishedAt.Time
		} else if containerStatus.RestartCount > 0 && containerStatus.State.Running != nil && !containerStatus.State.Running.StartedAt.IsZero() {
			// If no termination info but has restarts, use container start time as approximation
			lastRestart = containerStatus.State.Running.StartedAt.Time
		} else {
			// No restart info available, use pod creation time
			lastRestart = pod.CreationTimestamp.Time
		}

		lastRestartTime.With(labels).Set(float64(lastRestart.Unix()))

		glog.V(3).Infof(
			"container '%s' in pod '%s' has %d restarts, last restart: %s",
			containerStatus.Name,
			name,
			containerStatus.RestartCount,
			lastRestart.Format(time.RFC3339),
		)
	}
}

func (m *podStatusMetrics) Remove(pod corev1.Pod) {
	name := libk8s.NameFromPod(pod)
	glog.V(2).Infof("remove metrics for '%s'", name)

	// Delete status counter metrics completely to prevent stale data
	for _, p := range AvailableStatus {
		deleted := counter.Delete(prometheus.Labels{
			"name":      name.String(),
			"status":    p.String(),
			"app":       pod.Labels["app"],
			"namespace": pod.Namespace,
		})
		if deleted {
			glog.V(3).Infof("deleted metric for pod '%s' with status '%s'", name, p.String())
		}
	}

	// Delete restart count and last restart time metrics for each container
	for _, containerStatus := range pod.Status.ContainerStatuses {
		labels := prometheus.Labels{
			"name":      name.String(),
			"app":       pod.Labels["app"],
			"namespace": pod.Namespace,
			"container": containerStatus.Name,
		}

		restartDeleted := restartCount.Delete(labels)
		timeDeleted := lastRestartTime.Delete(labels)

		if restartDeleted || timeDeleted {
			glog.V(3).
				Infof("deleted restart metrics for container '%s' in pod '%s'", containerStatus.Name, name)
		}
	}
}

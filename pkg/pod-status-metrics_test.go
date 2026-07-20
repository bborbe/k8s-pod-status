// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/bborbe/k8s-pod-status/pkg"
)

var _ = Describe("PodStatusMetrics", func() {
	var (
		metrics pkg.PodStatusMetrics
		pod     corev1.Pod
	)

	BeforeEach(func() {
		metrics = pkg.NewPodStatusMetrics()

		// Create a sample pod for testing
		pod = corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "test-pod-12345",
				Namespace:         "test-namespace",
				CreationTimestamp: metav1.NewTime(time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)),
				Labels: map[string]string{
					"app": "test-app",
				},
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				ContainerStatuses: []corev1.ContainerStatus{
					{
						Name:         "test-container",
						RestartCount: 0,
						State: corev1.ContainerState{
							Running: &corev1.ContainerStateRunning{
								StartedAt: metav1.NewTime(
									time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
								),
							},
						},
					},
				},
			},
		}
	})

	Describe("NewPodStatusMetrics", func() {
		It("should create a new metrics instance", func() {
			metrics := pkg.NewPodStatusMetrics()
			Expect(metrics).NotTo(BeNil())
		})
	})

	Describe("Set", func() {
		Context("when setting metrics for a running pod", func() {
			It("should register metrics in prometheus registry", func() {
				metrics.Set(pod)

				// Verify metrics are registered by gathering them from the default registry
				metricFamilies, err := prometheus.DefaultGatherer.Gather()
				Expect(err).NotTo(HaveOccurred())

				// Check that pod status counter metric family exists
				var foundCounterFamily bool
				for _, family := range metricFamilies {
					if family.GetName() == "k8s_pod_status_counter" {
						foundCounterFamily = true
						// Check that we have metrics for our test pod
						for _, metric := range family.GetMetric() {
							labels := make(map[string]string)
							for _, label := range metric.GetLabel() {
								labels[label.GetName()] = label.GetValue()
							}
							if labels["name"] == "test-pod-12345" && labels["status"] == "Running" {
								Expect(metric.GetGauge().GetValue()).To(Equal(1.0))
							}
						}
						break
					}
				}
				Expect(foundCounterFamily).To(BeTrue())

				// Check restart count metric family exists
				var foundRestartFamily bool
				for _, family := range metricFamilies {
					if family.GetName() == "k8s_pod_status_restart_count" {
						foundRestartFamily = true
						break
					}
				}
				Expect(foundRestartFamily).To(BeTrue())

				// Check last restart time metric family exists
				var foundTimeFamily bool
				for _, family := range metricFamilies {
					if family.GetName() == "k8s_pod_status_last_restart_time" {
						foundTimeFamily = true
						break
					}
				}
				Expect(foundTimeFamily).To(BeTrue())
			})
		})

		Context("when setting metrics for different pod phases", func() {
			DescribeTable("should set correct status metrics",
				func(phase corev1.PodPhase, expectedStatus string) {
					pod.Status.Phase = phase
					metrics.Set(pod)

					metricFamilies, err := prometheus.DefaultGatherer.Gather()
					Expect(err).NotTo(HaveOccurred())

					// Find the counter metric family and verify the status
					for _, family := range metricFamilies {
						if family.GetName() == "k8s_pod_status_counter" {
							for _, metric := range family.GetMetric() {
								labels := make(map[string]string)
								for _, label := range metric.GetLabel() {
									labels[label.GetName()] = label.GetValue()
								}
								if labels["name"] == "test-pod-12345" &&
									labels["status"] == expectedStatus {
									Expect(metric.GetGauge().GetValue()).To(Equal(1.0))
								}
							}
							break
						}
					}
				},
				Entry("Running phase", corev1.PodRunning, "Running"),
				Entry("Failed phase", corev1.PodFailed, "Failed"),
				Entry("Pending phase", corev1.PodPending, "Pending"),
				Entry("Succeeded phase", corev1.PodSucceeded, "Succeeded"),
			)
		})

		Context("when pod has restart history", func() {
			BeforeEach(func() {
				pod.Status.ContainerStatuses[0].RestartCount = 3
				pod.Status.ContainerStatuses[0].LastTerminationState = corev1.ContainerState{
					Terminated: &corev1.ContainerStateTerminated{
						FinishedAt: metav1.NewTime(time.Date(2025, 1, 1, 11, 30, 0, 0, time.UTC)),
					},
				}
			})

			It("should set restart metrics correctly", func() {
				metrics.Set(pod)

				metricFamilies, err := prometheus.DefaultGatherer.Gather()
				Expect(err).NotTo(HaveOccurred())

				// Check restart count metric
				for _, family := range metricFamilies {
					if family.GetName() == "k8s_pod_status_restart_count" {
						for _, metric := range family.GetMetric() {
							labels := make(map[string]string)
							for _, label := range metric.GetLabel() {
								labels[label.GetName()] = label.GetValue()
							}
							if labels["name"] == "test-pod-12345" &&
								labels["container"] == "test-container" {
								Expect(metric.GetGauge().GetValue()).To(Equal(3.0))
							}
						}
						break
					}
				}

				// Check last restart time metric
				expectedTime := float64(time.Date(2025, 1, 1, 11, 30, 0, 0, time.UTC).Unix())
				for _, family := range metricFamilies {
					if family.GetName() == "k8s_pod_status_last_restart_time" {
						for _, metric := range family.GetMetric() {
							labels := make(map[string]string)
							for _, label := range metric.GetLabel() {
								labels[label.GetName()] = label.GetValue()
							}
							if labels["name"] == "test-pod-12345" &&
								labels["container"] == "test-container" {
								Expect(metric.GetGauge().GetValue()).To(Equal(expectedTime))
							}
						}
						break
					}
				}
			})
		})
	})

	Describe("Remove", func() {
		BeforeEach(func() {
			// Set metrics first
			metrics.Set(pod)
		})

		It("should call remove without errors", func() {
			// The main test is that Remove() doesn't panic or error
			// The key improvement is that our implementation now calls Delete()
			// instead of Set(0), preventing stale metrics in the max_over_time queries
			Expect(func() {
				metrics.Remove(pod)
			}).NotTo(Panic())
		})

		It("should handle remove for non-existent pod gracefully", func() {
			// Create a different pod that was never set
			nonExistentPod := corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "non-existent-pod",
					Namespace: "test-namespace",
					Labels: map[string]string{
						"app": "test-app",
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name:         "test-container",
							RestartCount: 0,
						},
					},
				},
			}

			// Should not panic when removing metrics for pod that was never set
			Expect(func() {
				metrics.Remove(nonExistentPod)
			}).NotTo(Panic())
		})
	})

	Describe("PodStatus", func() {
		DescribeTable("should return correct status for pod phases",
			func(phase corev1.PodPhase, expectedStatus pkg.Status) {
				pod.Status.Phase = phase
				status := pkg.PodStatus(pod)
				Expect(status).To(Equal(expectedStatus))
			},
			Entry("Pending phase", corev1.PodPending, pkg.PendingStatus),
			Entry("Running phase", corev1.PodRunning, pkg.RunningStatus),
			Entry("Succeeded phase", corev1.PodSucceeded, pkg.SucceededStatus),
			Entry("Failed phase", corev1.PodFailed, pkg.FailedStatus),
			Entry("Unknown phase", corev1.PodPhase("Unknown"), pkg.UnknownStatus),
		)
	})
})

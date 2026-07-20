// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import corev1 "k8s.io/api/core/v1"

const (
	SucceededStatus Status = "Succeeded"
	RunningStatus   Status = "Running"
	PendingStatus   Status = "Pending"
	FailedStatus    Status = "Failed"
	UnknownStatus   Status = "Unknown"
)

var AvailableStatus = Statuses{
	SucceededStatus,
	RunningStatus,
	PendingStatus,
	FailedStatus,
	UnknownStatus,
}

func PodStatus(pod corev1.Pod) Status {
	switch pod.Status.Phase {
	case corev1.PodPending:
		return PendingStatus
	case corev1.PodRunning:
		return RunningStatus
	case corev1.PodSucceeded:
		return SucceededStatus
	case corev1.PodFailed:
		return FailedStatus
	default:
		return UnknownStatus
	}
}

type Status string

func (f Status) String() string {
	return string(f)
}

type Statuses []Status

func (p Statuses) Contains(status Status) bool {
	for _, ph := range p {
		if ph == status {
			return true
		}
	}
	return false
}

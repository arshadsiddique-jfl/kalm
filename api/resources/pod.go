package resources

import (
	"github.com/kalmhq/kalm/controller/api/v1alpha1"
	coreV1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getPodStatusPhase returns one of four pod status phases (Pending, Running, Succeeded, Failed)
func getPodStatusPhase(pod coreV1.Pod, warnings []coreV1.Event, workloadType v1alpha1.WorkloadType) coreV1.PodPhase {

	// For Completed Cronjob
	isTerminatedCompleted := false
	for i := len(pod.Status.ContainerStatuses) - 1; i >= 0; i-- {
		container := pod.Status.ContainerStatuses[i]
		if container.State.Terminated != nil && container.State.Terminated.Reason == "Completed" {
			isTerminatedCompleted = true
		}
	}
	if workloadType == v1alpha1.WorkloadTypeCronjob && isTerminatedCompleted {
		return coreV1.PodSucceeded
	}

	// For terminated pods that failed
	if pod.Status.Phase == coreV1.PodFailed {
		return coreV1.PodFailed
	}

	// For successfully terminated pods
	if pod.Status.Phase == coreV1.PodSucceeded {
		return coreV1.PodSucceeded
	}

	ready := false
	initialized := false
	for _, c := range pod.Status.Conditions {
		if c.Type == coreV1.PodReady {
			ready = c.Status == coreV1.ConditionTrue
		}
		if c.Type == coreV1.PodInitialized {
			initialized = c.Status == coreV1.ConditionTrue
		}
	}

	if initialized && ready && pod.Status.Phase == coreV1.PodRunning {
		return coreV1.PodRunning
	}

	// If the pod would otherwise be pending but has warning then label it as
	// failed and show and error to the user.
	if len(warnings) > 0 {
		return coreV1.PodFailed
	}

	// pending
	return coreV1.PodPending
}

type PodListChannel struct {
	List  chan *coreV1.PodList
	Error chan error
}

func (resourceManager *ResourceManager) GetPodListChannel(opts ...client.ListOption) *PodListChannel {
	channel := &PodListChannel{
		List:  make(chan *coreV1.PodList, 1),
		Error: make(chan error, 1),
	}

	go func() {
		var list coreV1.PodList
		err := resourceManager.List(&list, opts...)
		channel.List <- &list
		channel.Error <- err
	}()

	return channel
}

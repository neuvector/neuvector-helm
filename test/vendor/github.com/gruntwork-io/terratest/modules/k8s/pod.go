package k8s

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
)

// ListPods will look for pods in the given namespace that match the given filters and return them. This will fail the
// test if there is an error.
func ListPods(t *testing.T, options *KubectlOptions, filters metav1.ListOptions) []corev1.Pod {
	pods, err := ListPodsE(t, options, filters)
	require.NoError(t, err)
	return pods
}

// ListPodsE will look for pods in the given namespace that match the given filters and return them.
func ListPodsE(t *testing.T, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.Pod, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.CoreV1().Pods(options.Namespace).List(filters)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetPod returns a Kubernetes pod resource in the provided namespace with the given name. This will
// fail the test if there is an error.
func GetPod(t *testing.T, options *KubectlOptions, podName string) *corev1.Pod {
	pod, err := GetPodE(t, options, podName)
	require.NoError(t, err)
	return pod
}

// GetPodE returns a Kubernetes pod resource in the provided namespace with the given name.
func GetPodE(t *testing.T, options *KubectlOptions, podName string) (*corev1.Pod, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().Pods(options.Namespace).Get(podName, metav1.GetOptions{})
}

// WaitUntilNumPodsCreated waits until the desired number of pods are created that match the provided filter. This will
// retry the check for the specified amount of times, sleeping for the provided duration between each try. This will
// fail the test if the retry times out.
func WaitUntilNumPodsCreated(
	t *testing.T,
	options *KubectlOptions,
	filters metav1.ListOptions,
	desiredCount int,
	retries int,
	sleepBetweenRetries time.Duration,
) {
	require.NoError(t, WaitUntilNumPodsCreatedE(t, options, filters, desiredCount, retries, sleepBetweenRetries))
}

// WaitUntilNumPodsCreatedE waits until the desired number of pods are created that match the provided filter. This will
// retry the check for the specified amount of times, sleeping for the provided duration between each try.
func WaitUntilNumPodsCreatedE(
	t *testing.T,
	options *KubectlOptions,
	filters metav1.ListOptions,
	desiredCount int,
	retries int,
	sleepBetweenRetries time.Duration,
) error {
	statusMsg := fmt.Sprintf("Wait for num pods created to match desired count %d.", desiredCount)
	message, err := retry.DoWithRetryE(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			pods, err := ListPodsE(t, options, filters)
			if err != nil {
				return "", err
			}
			if len(pods) != desiredCount {
				return "", DesiredNumberOfPodsNotCreated{Filter: filters, DesiredCount: desiredCount}
			}
			return "Desired number of Pods created", nil
		},
	)
	if err != nil {
		logger.Logf(t, "Timedout waiting for the desired number of Pods to be created: %s", err)
		return err
	}
	logger.Logf(t, message)
	return nil
}

// WaitUntilPodAvailable waits until the pod is running, retrying the check for the specified amount of times, sleeping
// for the provided duration between each try. This will fail the test if there is an error or if the check times out.
func WaitUntilPodAvailable(t *testing.T, options *KubectlOptions, podName string, retries int, sleepBetweenRetries time.Duration) {
	require.NoError(t, WaitUntilPodAvailableE(t, options, podName, retries, sleepBetweenRetries))
}

// WaitUntilPodAvailableE waits until the pod is running, retrying the check for the specified amount of times, sleeping
// for the provided duration between each try.
func WaitUntilPodAvailableE(t *testing.T, options *KubectlOptions, podName string, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for pod %s to be provisioned.", podName)
	message, err := retry.DoWithRetryE(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			pod, err := GetPodE(t, options, podName)
			if err != nil {
				return "", err
			}
			if !IsPodAvailable(pod) {
				return "", NewPodNotAvailableError(pod)
			}
			return "Pod is now available", nil
		},
	)
	if err != nil {
		logger.Logf(t, "Timedout waiting for Pod to be provisioned: %s", err)
		return err
	}
	logger.Logf(t, message)
	return nil
}

// IsPodAvailable returns true if the pod is running.
func IsPodAvailable(pod *corev1.Pod) bool {
	return pod.Status.Phase == corev1.PodRunning
}

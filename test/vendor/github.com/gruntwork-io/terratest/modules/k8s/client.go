package k8s

import (
	"testing"

	"k8s.io/client-go/kubernetes"

	// The following line loads the gcp plugin which is required to authenticate against GKE clusters.
	// See: https://github.com/kubernetes/client-go/issues/242
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/gruntwork-io/terratest/modules/logger"
)

// GetKubernetesClientE returns a Kubernetes API client that can be used to make requests.
func GetKubernetesClientE(t *testing.T) (*kubernetes.Clientset, error) {
	kubeConfigPath, err := GetKubeConfigPathE(t)
	if err != nil {
		return nil, err
	}

	options := NewKubectlOptions("", kubeConfigPath)
	return GetKubernetesClientFromOptionsE(t, options)
}

// GetKubernetesClientFromOptionsE returns a Kubernetes API client given a configured KubectlOptions object.
func GetKubernetesClientFromOptionsE(t *testing.T, options *KubectlOptions) (*kubernetes.Clientset, error) {
	var err error

	kubeConfigPath, err := options.GetConfigPath(t)
	if err != nil {
		return nil, err
	}
	logger.Logf(t, "Configuring kubectl using config file %s with context %s", kubeConfigPath, options.ContextName)
	// Load API config (instead of more low level ClientConfig)
	config, err := LoadApiClientConfigE(kubeConfigPath, options.ContextName)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

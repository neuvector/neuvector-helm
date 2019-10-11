package k8s

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngressNotAvailable is returned when a Kubernetes service is not yet available to accept traffic.
type IngressNotAvailable struct {
	ingress *extensionsv1beta1.Ingress
}

func (err IngressNotAvailable) Error() string {
	return fmt.Sprintf("Ingress %s is not available", err.ingress.Name)
}

// UnknownKubeResourceType is returned if the given resource type does not match the list of known resource types.
type UnknownKubeResourceType struct {
	ResourceType KubeResourceType
}

func (err UnknownKubeResourceType) Error() string {
	return fmt.Sprintf("ResourceType ID %d is unknown", err.ResourceType)
}

// DesiredNumberOfPodsNotCreated is returned when the number of pods matching a filter condition does not match the
// desired number of Pods.
type DesiredNumberOfPodsNotCreated struct {
	Filter       metav1.ListOptions
	DesiredCount int
}

func (err DesiredNumberOfPodsNotCreated) Error() string {
	return fmt.Sprintf("Desired number of pods (%d) matching filter %v not yet created", err.DesiredCount, err.Filter)
}

// ServiceAccountTokenNotAvailable is returned when a Kubernetes ServiceAccount does not have a token provisioned yet.
type ServiceAccountTokenNotAvailable struct {
	Name string
}

func (err ServiceAccountTokenNotAvailable) Error() string {
	return fmt.Sprintf("ServiceAccount %s does not have a token yet.", err.Name)
}

// PodNotAvailable is returned when a Kubernetes service is not yet available to accept traffic.
type PodNotAvailable struct {
	pod *corev1.Pod
}

func (err PodNotAvailable) Error() string {
	return fmt.Sprintf("Pod %s is not available", err.pod.Name)
}

func NewPodNotAvailableError(pod *corev1.Pod) PodNotAvailable {
	return PodNotAvailable{pod}
}

// ServiceNotAvailable is returned when a Kubernetes service is not yet available to accept traffic.
type ServiceNotAvailable struct {
	service *corev1.Service
}

func (err ServiceNotAvailable) Error() string {
	return fmt.Sprintf("Service %s is not available", err.service.Name)
}

func NewServiceNotAvailableError(service *corev1.Service) ServiceNotAvailable {
	return ServiceNotAvailable{service}
}

// UnknownServiceType is returned when a Kubernetes service has a type that is not yet handled by the test functions.
type UnknownServiceType struct {
	service *corev1.Service
}

func (err UnknownServiceType) Error() string {
	return fmt.Sprintf("Service %s has an unknown service type", err.service.Name)
}

func NewUnknownServiceTypeError(service *corev1.Service) UnknownServiceType {
	return UnknownServiceType{service}
}

// UnknownServicePort is returned when the given service port is not an exported port of the service.
type UnknownServicePort struct {
	service *corev1.Service
	port    int32
}

func (err UnknownServicePort) Error() string {
	return fmt.Sprintf("Port %d is not a part of the service %s", err.port, err.service.Name)
}

func NewUnknownServicePortError(service *corev1.Service, port int32) UnknownServicePort {
	return UnknownServicePort{service, port}
}

// NoNodesInKubernetes is returned when the Kubernetes cluster has no nodes registered.
type NoNodesInKubernetes struct{}

func (err NoNodesInKubernetes) Error() string {
	return "There are no nodes in the Kubernetes cluster"
}

func NewNoNodesInKubernetesError() NoNodesInKubernetes {
	return NoNodesInKubernetes{}
}

// NodeHasNoHostname is returned when a Kubernetes node has no discernible hostname
type NodeHasNoHostname struct {
	node *corev1.Node
}

func (err NodeHasNoHostname) Error() string {
	return fmt.Sprintf("Node %s has no hostname", err.node.Name)
}

func NewNodeHasNoHostnameError(node *corev1.Node) NodeHasNoHostname {
	return NodeHasNoHostname{node}
}

// MalformedNodeID is returned when a Kubernetes node has a malformed node id scheme
type MalformedNodeID struct {
	node *corev1.Node
}

func (err MalformedNodeID) Error() string {
	return fmt.Sprintf("Node %s has malformed ID %s", err.node.Name, err.node.Spec.ProviderID)
}

func NewMalformedNodeIDError(node *corev1.Node) MalformedNodeID {
	return MalformedNodeID{node}
}

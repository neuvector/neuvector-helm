package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	corev1 "k8s.io/api/core/v1"
)

func checkControllerServiceDefault(t *testing.T, svc corev1.Service) {
	if svc.Name != "neuvector-svc-controller" {
		t.Errorf("Service name is wrong. name=%v\n", svc.Name)
	}
	if svc.Spec.Type != "" || svc.Spec.ClusterIP != "None" {
		t.Errorf("Service type is wrong. type=%v clusterIP=%v\n", svc.Spec.Type, svc.Spec.ClusterIP)
	}
	if len(svc.Spec.Ports) != 3 {
		t.Errorf("Service port is wrong. ports=%+v\n", svc.Spec.Ports)
	}
	if app, ok := svc.Spec.Selector["app"]; !ok || app != "neuvector-controller-pod" {
		t.Errorf("Service selector is invalid. selector=%+v\n", svc.Spec.Selector)
	}
}

func checkControllerServiceAPI(t *testing.T, svc corev1.Service, svcType string) {
	if svc.Name != "neuvector-svc-controller-api" {
		t.Errorf("Service name is wrong. name=%v\n", svc.Name)
	}
	if string(svc.Spec.Type) != svcType {
		t.Errorf("Service type is wrong. type=%v\n", svc.Spec.Type)
	}
	if len(svc.Spec.Ports) != 1 || svc.Spec.Ports[0].Port != 10443 {
		t.Errorf("Service port is wrong. ports=%+v\n", svc.Spec.Ports)
	}
	if app, ok := svc.Spec.Selector["app"]; !ok || app != "neuvector-controller-pod" {
		t.Errorf("Service selector is invalid. selector=%+v\n", svc.Spec.Selector)
	}
}

func checkControllerServiceFedMaster(t *testing.T, svc corev1.Service, svcType string) {
	if svc.Name != "neuvector-svc-controller-fed-master" {
		t.Errorf("Service name is wrong. name=%v\n", svc.Name)
	}
	if string(svc.Spec.Type) != svcType {
		t.Errorf("Service type is wrong. type=%v\n", svc.Spec.Type)
	}
	if len(svc.Spec.Ports) != 1 || svc.Spec.Ports[0].Port != 11443 {
		t.Errorf("Service port is wrong. ports=%+v\n", svc.Spec.Ports)
	}
	if app, ok := svc.Spec.Selector["app"]; !ok || app != "neuvector-controller-pod" {
		t.Errorf("Service selector is invalid. selector=%+v\n", svc.Spec.Selector)
	}
}

func checkControllerServiceFedManaged(t *testing.T, svc corev1.Service, svcType string) {
	if svc.Name != "neuvector-svc-controller-fed-managed" {
		t.Errorf("Service name is wrong. name=%v\n", svc.Name)
	}
	if string(svc.Spec.Type) != svcType {
		t.Errorf("Service type is wrong. type=%v\n", svc.Spec.Type)
	}
	if len(svc.Spec.Ports) != 1 || svc.Spec.Ports[0].Port != 10443 {
		t.Errorf("Service port is wrong. ports=%+v\n", svc.Spec.Ports)
	}
	if app, ok := svc.Spec.Selector["app"]; !ok || app != "neuvector-controller-pod" {
		t.Errorf("Service selector is invalid. selector=%+v\n", svc.Spec.Selector)
	}
}

func TestControllerService(t *testing.T) {
	helmChartPath := ".."

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Test controller service
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/controller-service.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	var svc corev1.Service
	helm.UnmarshalK8SYaml(t, outs[0], &svc)

	checkControllerServiceDefault(t, svc)
}

func TestControllerServiceAPI(t *testing.T) {
	helmChartPath := ".."

	options := &helm.Options{
		SetValues: map[string]string{
			"controller.apisvc.type": "NodePort",
		},
	}

	// Test controller service
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/controller-service.yaml"})
	outs := splitYaml(out)

	if len(outs) != 2 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var svc corev1.Service
		helm.UnmarshalK8SYaml(t, output, &svc)

		switch i {
		case 0:
			checkControllerServiceDefault(t, svc)
		case 1:
			checkControllerServiceAPI(t, svc, "NodePort")
		}
	}
}

func TestControllerFedMaster(t *testing.T) {
	helmChartPath := ".."

	options := &helm.Options{
		SetValues: map[string]string{
			"controller.federation.mastersvc.type": "NodePort",
		},
	}

	// Test controller service
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/controller-service.yaml"})
	outs := splitYaml(out)

	if len(outs) != 2 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var svc corev1.Service
		helm.UnmarshalK8SYaml(t, output, &svc)

		switch i {
		case 0:
			checkControllerServiceDefault(t, svc)
		case 2:
			checkControllerServiceFedMaster(t, svc, "NodePort")
		}
	}
}

func TestControllerServiceAPIandFedMaster(t *testing.T) {
	helmChartPath := ".."

	options := &helm.Options{
		SetValues: map[string]string{
			"controller.apisvc.type":               "NodePort",
			"controller.federation.mastersvc.type": "NodePort",
		},
	}

	// Test controller service
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/controller-service.yaml"})
	outs := splitYaml(out)

	if len(outs) != 3 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var svc corev1.Service
		helm.UnmarshalK8SYaml(t, output, &svc)

		switch i {
		case 0:
			checkControllerServiceDefault(t, svc)
		case 1:
			checkControllerServiceAPI(t, svc, "NodePort")
		case 2:
			checkControllerServiceFedMaster(t, svc, "NodePort")
		}
	}
}

func TestControllerServiceAPIandFedManaged(t *testing.T) {
	helmChartPath := ".."

	options := &helm.Options{
		SetValues: map[string]string{
			"controller.apisvc.type":                "NodePort",
			"controller.federation.managedsvc.type": "NodePort",
		},
	}

	// Test controller service
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/controller-service.yaml"})
	outs := splitYaml(out)

	if len(outs) != 3 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var svc corev1.Service
		helm.UnmarshalK8SYaml(t, output, &svc)

		switch i {
		case 0:
			checkControllerServiceDefault(t, svc)
		case 1:
			checkControllerServiceAPI(t, svc, "NodePort")
		case 2:
			checkControllerServiceFedManaged(t, svc, "NodePort")
		}
	}
}

func TestControllerServiceIngress(t *testing.T) {
	helmChartPath := ".."

	options := &helm.Options{
		SetValues: map[string]string{
			"controller.ingress.enabled": "true",
		},
	}

	// Test controller service
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/controller-service.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	var svc corev1.Service
	helm.UnmarshalK8SYaml(t, outs[0], &svc)

	checkControllerServiceDefault(t, svc)
}

// -- manager

func checkManagerService(t *testing.T, svc corev1.Service, svcType string) {
	if svc.Name != "neuvector-service-webui" {
		t.Errorf("Service name is wrong. name=%v\n", svc.Name)
	}
	if string(svc.Spec.Type) != svcType {
		t.Errorf("Service type is wrong. type=%v\n", svc.Spec.Type)
	}
	if len(svc.Spec.Ports) != 1 || svc.Spec.Ports[0].Port != 8443 {
		t.Errorf("Service port is wrong. ports=%+v\n", svc.Spec.Ports)
	}
	if app, ok := svc.Spec.Selector["app"]; !ok || app != "neuvector-manager-pod" {
		t.Errorf("Service selector is invalid. selector=%+v\n", svc.Spec.Selector)
	}
}

func TestManagerService(t *testing.T) {
	helmChartPath := ".."

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Test controller service
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/manager-service.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	var svc corev1.Service
	helm.UnmarshalK8SYaml(t, outs[0], &svc)

	checkManagerService(t, svc, "NodePort")
}

func TestManagerServiceLB(t *testing.T) {
	helmChartPath := ".."

	svcType := "LoadBalancer"
	options := &helm.Options{
		SetValues: map[string]string{
			"manager.svc.type": svcType,
		},
	}

	// Test controller service
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/manager-service.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	var svc corev1.Service
	helm.UnmarshalK8SYaml(t, outs[0], &svc)

	checkManagerService(t, svc, svcType)
}

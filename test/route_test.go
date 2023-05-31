package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	corev1 "k8s.io/api/core/v1"
)

func TestManagerRouter(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"openshift":             "true",
			"manager.route.enabled": "true",
		},
	}

	// Test controller service
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/manager-route.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	var svc corev1.Service
	helm.UnmarshalK8SYaml(t, outs[0], &svc)
}

func TestControllerAPISVCRouter(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"openshift":                       "true",
			"controller.apisvc.route.enabled": "true",
		},
	}

	// Test controller service
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/controller-route.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	var svc corev1.Service
	helm.UnmarshalK8SYaml(t, outs[0], &svc)
}

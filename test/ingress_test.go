package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
)

func TestIngressController(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"controller.ingress.enabled":                       "true",
			"controller.federation.mastersvc.ingress.enabled":  "true",
			"controller.federation.managedsvc.ingress.enabled": "true",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/controller-ingress.yaml"})
	outs := splitYaml(out)

	if len(outs) != 3 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var ing extv1beta1.Ingress
		helm.UnmarshalK8SYaml(t, output, &ing)

		switch i {
		case 0:
			if ing.Name != "neuvector-restapi-ingress" {
				t.Errorf("Ingress name is wrong. name=%v\n", ing.Name)
			}
		case 1:
			if ing.Name != "neuvector-mastersvc-ingress" {
				t.Errorf("Ingress name is wrong. name=%v\n", ing.Name)
			}
		case 2:
			if ing.Name != "neuvector-managedsvc-ingress" {
				t.Errorf("Ingress name is wrong. name=%v\n", ing.Name)
			}
		}
	}
}

func TestIngressManager(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"manager.ingress.enabled": "true",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/manager-ingress.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var ing extv1beta1.Ingress
		helm.UnmarshalK8SYaml(t, output, &ing)

		switch i {
		case 0:
			if ing.Name != "neuvector-webui-ingress" {
				t.Errorf("Ingress name is wrong. name=%v\n", ing.Name)
			}
		}
	}
}

package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func TestCRD(t *testing.T) {
	helmChartPath := "../charts/crd"

	options := &helm.Options{
		SetValues: map[string]string{
			"controller.ranchersso.enabled": "true",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/crd.yaml"})
	outs := splitYaml(out)

	if len(outs) != 24 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

func TestCoreCRD(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"controller.ranchersso.enabled": "true",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/crd.yaml"})
	outs := splitYaml(out)

	if len(outs) != 24 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

func TestCoreCRDDisabled(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"crdwebhook.enabled": "false",
			"controller.ranchersso.enabled": "false",
		},
	}

	// Test ingress
	out, _ := helm.RenderTemplateE(t, options, helmChartPath, nvRel, []string{"templates/crd.yaml"})
	outs := splitYaml(out)

	if len(outs) != 0 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

func TestCoreCRDDisabledRancherSSOEnabled(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"crdwebhook.enabled": "false",
			"controller.ranchersso.enabled": "true",
		},
	}

	// Test ingress
	out, _ := helm.RenderTemplateE(t, options, helmChartPath, nvRel, []string{"templates/crd.yaml"})
	outs := splitYaml(out)

	if len(outs) != 16 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

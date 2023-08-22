package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func TestRegistryAdapter(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"leastPrivilege":      "true",
			"cve.adapter.enabled": "true",
			"generateSecret":      "false",
		},
	}

	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/registry-adapter.yaml"})
	outs := splitYaml(out)

	if len(outs) != 2 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

func TestRegistryAdapterIngress(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"leastPrivilege":              "true",
			"cve.adapter.enabled":         "true",
			"cve.adapter.ingress.enabled": "true",
			"openshift":                   "true",
			"generateSecret":              "false",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/registry-adapter-ingress.yaml"})
	outs := splitYaml(out)

	if len(outs) != 2 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

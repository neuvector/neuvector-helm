package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func TestServiceAccountLeastPrivilege(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"leastPrivilege": "true",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/serviceaccount-least.yaml"})
	outs := splitYaml(out)

	if len(outs) != 4 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

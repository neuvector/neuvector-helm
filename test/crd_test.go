package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func TestCRD(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/crd.yaml"})
	outs := splitYaml(out)

	if len(outs) != 7 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

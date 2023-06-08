package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func TestExporter(t *testing.T) {
	helmChartPath := "../charts/monitor"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/exporter-deployment.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

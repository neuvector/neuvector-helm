package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
)

func TestRoleBinding(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/rolebinding.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

func TestClusterRole(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/clusterrole.yaml"})
	outs := splitYaml(out)

	if len(outs) != 3 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

func TestClusterRoleBinding(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/clusterrolebinding.yaml"})
	outs := splitYaml(out)

	if len(outs) != 4 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

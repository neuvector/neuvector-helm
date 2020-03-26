package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	appsv1 "k8s.io/api/apps/v1"
)

func TestScanner(t *testing.T) {
	helmChartPath := ".."

	options := &helm.Options{
		SetValues: map[string]string{
			"cve.scanner.enabled": "true",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/scanner-deployment.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var scanner appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &scanner)

		switch i {
		case 0:
			if *scanner.Spec.Replicas != int32(3) {
				t.Errorf("Incorrect scanner replicas. replicas=%v\n", scanner.Spec.Replicas)
			}
		}
	}
}

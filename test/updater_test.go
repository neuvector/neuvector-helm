package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	batv1beta1 "k8s.io/api/batch/v1beta1"
)

func TestUpdater(t *testing.T) {
	helmChartPath := ".."

	options := &helm.Options{
		SetValues: map[string]string{
			"cve.updater.enabled": "true",
			"cve.scanner.enabled": "false",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/updater-cronjob.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var job batv1beta1.CronJob
		helm.UnmarshalK8SYaml(t, output, &job)

		switch i {
		case 0:
			if job.Name != "neuvector-updater-pod" {
				t.Errorf("Incorrect cronjob name. name=%v\n", job.Name)
			}
			if job.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Lifecycle != nil {
				t.Errorf("No need to update scanner.\n")
			}
		}
	}
}

func TestUpdaterWithScanner(t *testing.T) {
	helmChartPath := ".."

	options := &helm.Options{
		SetValues: map[string]string{
			"cve.updater.enabled": "true",
			"cve.scanner.enabled": "true",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/updater-cronjob.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var job batv1beta1.CronJob
		helm.UnmarshalK8SYaml(t, output, &job)

		switch i {
		case 0:
			if job.Name != "neuvector-updater-pod" {
				t.Errorf("Incorrect cronjob name. name=%v\n", job.Name)
			}
			if job.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Lifecycle == nil {
				t.Errorf("Missing update scanner.\n")
			}
		}
	}
}

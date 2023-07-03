package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	appsv1 "k8s.io/api/apps/v1"
)

func TestEnforcerDaemonset(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/enforcer-daemonset.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

func TestEnforcerDaemonsetRuntime(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"crio.enabled": "true",
			"crio.path":    "/var/run/crio.sock",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/enforcer-daemonset.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	var ds appsv1.DaemonSet
	helm.UnmarshalK8SYaml(t, outs[0], &ds)

	var found bool
	for _, m := range ds.Spec.Template.Spec.Containers[0].VolumeMounts {
		if m.Name == "runtime-sock" && m.MountPath == "/var/run/crio/crio.sock" {
			found = true
		}
	}
	if !found {
		t.Errorf("Volume mount for the runtime socket is not found. Mounts=%+v\n",
			ds.Spec.Template.Spec.Containers[0].VolumeMounts)
	}

	for _, v := range ds.Spec.Template.Spec.Volumes {
		if v.Name == "runtime-sock" && v.HostPath.Path == "/var/run/crio.sock" {
			found = true
		}
	}
	if !found {
		t.Errorf("Volume for the runtime socket is not found. Volumes=%+v\n",
			ds.Spec.Template.Spec.Volumes)
	}
}

func TestEnforcerDaemonsetLeastPrivilege(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"leastPrivilege": "true",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/enforcer-daemonset.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var dep appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &dep)

		switch i {
		case 0:
			if dep.Spec.Template.Spec.ServiceAccountName != "enforcer" {
				t.Errorf("Incorrect service account. sa=%+v\n", dep.Spec.Template.Spec.ServiceAccountName)
			}
		}
	}
}

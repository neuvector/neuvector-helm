package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	appsv1 "k8s.io/api/apps/v1"
)

func TestControllerDeployment(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/controller-deployment.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

func TestControllerDeploymentRegistry(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"registry": "registry.neuvector.com",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/controller-deployment.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	var dep appsv1.Deployment
	helm.UnmarshalK8SYaml(t, outs[0], &dep)
	if dep.Spec.Template.Spec.Containers[0].Image != "registry.neuvector.com/controller:latest" {
		t.Errorf("Image location is wrong, %v\n", dep.Spec.Template.Spec.Containers[0].Image)
	}
}

func TestControllerDeploymentOEM(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"registry": "registry.neuvector.com",
			"oem":      "oem",
			"tag":      "0.9",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/controller-deployment.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	var dep appsv1.Deployment
	helm.UnmarshalK8SYaml(t, outs[0], &dep)
	if dep.Spec.Template.Spec.Containers[0].Image != "registry.neuvector.com/oem/controller:0.9" {
		t.Errorf("Image location is wrong, %v\n", dep.Spec.Template.Spec.Containers[0].Image)
	}
}

func TestControllerDeploymentCert(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"controller.certificate.secret": "https-cert",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/controller-deployment.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

func TestControllerDeploymentDisrupt(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"controller.disruptionbudget": "2",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/controller-deployment.yaml"})
	outs := splitYaml(out)

	if len(outs) != 2 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

// --

func checkManagerDeployment(t *testing.T, dep appsv1.Deployment, ssl bool) {
	if dep.Name != "neuvector-manager-pod" {
		t.Errorf("Deployment name is wrong. name=%v\n", dep.Name)
	}

	if len(dep.Spec.Template.Spec.Containers) != 1 {
		t.Errorf("Container spec count is wrong. count=%v\n", len(dep.Spec.Template.Spec.Containers))
	}

	if ssl {
		for _, kv := range dep.Spec.Template.Spec.Containers[0].Env {
			if kv.Name == "MANAGER_SSL" {
				t.Errorf("MANAGER_SSL env should not exist.\n")
				break
			}
		}
	} else {
		found := false
		for _, kv := range dep.Spec.Template.Spec.Containers[0].Env {
			if kv.Name == "MANAGER_SSL" && kv.Value == "off" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("MANAGER_SSL env does not exist or value is not 'off'. env=%+v\n", dep.Spec.Template.Spec.Containers[0].Env)
		}
	}
}

func TestManagerDeployment(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/manager-deployment.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var dep appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &dep)

		switch i {
		case 0:
			checkManagerDeployment(t, dep, true)
		}
	}
}

func TestManagerDeploymentCert(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"manager.certificate.secret": "https-cert",
		},
	}

	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/manager-deployment.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

func TestManagerDeploymentNonSSL(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"manager.env.ssl": "false",
		},
	}

	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/manager-deployment.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var dep appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &dep)

		switch i {
		case 0:
			checkManagerDeployment(t, dep, false)
		}
	}
}

// --

func TestControllerDeploymentLeastPrivilege(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"leastPrivilege": "true",
		},
	}

	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/controller-deployment.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var dep appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &dep)

		switch i {
		case 0:
			if dep.Spec.Template.Spec.ServiceAccountName != "controller" {
				t.Errorf("Incorrect service account. sa=%+v\n", dep.Spec.Template.Spec.ServiceAccountName)
			}
		}
	}
}

func TestManagerDeploymentLeastPrivilege(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"leastPrivilege": "true",
		},
	}

	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/manager-deployment.yaml"})
	outs := splitYaml(out)

	if len(outs) != 1 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		var dep appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &dep)

		switch i {
		case 0:
			if dep.Spec.Template.Spec.ServiceAccountName != "basic" {
				t.Errorf("Incorrect service account. sa=%+v\n", dep.Spec.Template.Spec.ServiceAccountName)
			}
		}
	}
}

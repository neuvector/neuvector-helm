package test

import (
	"testing"

	"strings"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/assert"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestControllerDeployment(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"autoGenerateCert": "false",
		},
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
			"registry":         "registry.neuvector.com",
			"autoGenerateCert": "false",
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
	if strings.HasPrefix(dep.Spec.Template.Spec.Containers[0].Image, "registry.neuvector.com/controler:") {
		t.Errorf("Image location is wrong, %v\n", dep.Spec.Template.Spec.Containers[0].Image)
	}
}

func TestControllerDeploymentOEM(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"registry":         "registry.neuvector.com",
			"oem":              "oem",
			"tag":              "0.9",
			"autoGenerateCert": "false",
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
			"autoGenerateCert":              "false",
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
			"autoGenerateCert":            "false",
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
	if dep.Name != "neuvector-manager" {
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
		SetValues: map[string]string{
			"autoGenerateCert": "false",
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
			checkManagerDeployment(t, dep, true)
		}
	}
}

func TestManagerDeploymentCert(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"manager.certificate.secret": "https-cert",
			"autoGenerateCert":           "false",
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
			"manager.env.ssl":  "false",
			"autoGenerateCert": "false",
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
			"leastPrivilege":   "true",
			"autoGenerateCert": "false",
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
			"leastPrivilege":   "true",
			"autoGenerateCert": "false",
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

func TestControllerSecrets(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{
		"templates/controller-deployment.yaml",
		"templates/controller-secret.yaml",
	})
	outs := splitYaml(out)

	// Secret will be created
	for _, output := range outs {
		var secret corev1.Secret
		helm.UnmarshalK8SYaml(t, output, &secret)
		if secret.Name == "neuvector-controller" {
			assert.NotNil(t, secret.Data)
			assert.NotEmpty(t, secret.Data["ssl-cert.key"])
			assert.NotEmpty(t, secret.Data["ssl-cert.pem"])
		}
	}

	out = helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{
		"templates/controller-deployment.yaml",
		"templates/controller-secret.yaml",
	})
	outs = splitYaml(out)

	// Secret will be created and mounted
	for _, output := range outs {
		var dep appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &dep)
		if dep.Name == "neuvector-controller" {

			assert.Contains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "cert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "neuvector-controller",
					},
				},
			})
		}
	}

	for _, output := range outs {
		var dep appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &dep)
		if dep.Name == "neuvector-controller" {

			// cert and usercert will be mounted.
			assert.Contains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "cert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "neuvector-controller",
					},
				},
			})

			assert.NotContains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "usercert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "nv-ssl",
					},
				},
			})

			for _, container := range dep.Spec.Template.Spec.Containers {
				if container.Name == "neuvector-controller" {

					assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "cert",
						MountPath: "/etc/neuvector/certs/ssl-cert.key",
						SubPath:   "ssl-cert.key",
						ReadOnly:  true,
					})

					assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "cert",
						MountPath: "/etc/neuvector/certs/ssl-cert.pem",
						SubPath:   "ssl-cert.pem",
						ReadOnly:  true,
					})

				}

			}

		}
	}
}

func TestControllerNoSecrets(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"autoGenerateCert": "false",
		},
	}

	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{
		"templates/controller-deployment.yaml",
		//"templates/controller-secret.yaml",
	})
	outs := splitYaml(out)

	// Secret will not be created
	for _, output := range outs {
		var secret corev1.Secret
		helm.UnmarshalK8SYaml(t, output, &secret)
		assert.NotEqual(t, "neuvector-controller", secret.Name)
	}

	for _, output := range outs {
		var dep appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &dep)
		if dep.Name == "neuvector-controller" {

			// cert and usercert will be mounted.
			assert.NotContains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "cert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "neuvector-controller",
					},
				},
			})

			assert.NotContains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "usercert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "nv-ssl",
					},
				},
			})

			for _, container := range dep.Spec.Template.Spec.Containers {
				if container.Name == "neuvector-controller" {

					assert.NotContains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "cert",
						MountPath: "/etc/neuvector/certs/ssl-cert.key",
						SubPath:   "ssl-cert.key",
						ReadOnly:  true,
					})

					assert.NotContains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "cert",
						MountPath: "/etc/neuvector/certs/ssl-cert.pem",
						SubPath:   "ssl-cert.pem",
						ReadOnly:  true,
					})
				}

			}

		}
	}
}

func TestControllerWithOnlySSLKeys(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"controller.certificate.secret":  "nv-ssl",
			"controller.certificate.keyFile": "key3.pem",
			"controller.certificate.pemFile": "cert3.pem",
		},
	}

	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{
		"templates/controller-deployment.yaml",
		"templates/controller-secret.yaml",
	})
	outs := splitYaml(out)

	// Secret will be created and mounted
	for _, output := range outs {
		var dep appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &dep)
		if dep.Name == "neuvector-controller" {

			// cert and usercert will be mounted.
			assert.Contains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "cert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "neuvector-controller",
					},
				},
			})

			assert.Contains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "usercert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "nv-ssl",
					},
				},
			})

			for _, container := range dep.Spec.Template.Spec.Containers {
				if container.Name == "neuvector-controller" {

					assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "usercert",
						MountPath: "/etc/neuvector/certs/ssl-cert.key",
						SubPath:   "key3.pem",
						ReadOnly:  true,
					})

					assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "usercert",
						MountPath: "/etc/neuvector/certs/ssl-cert.pem",
						SubPath:   "cert3.pem",
						ReadOnly:  true,
					})
				}

			}

		}
	}
}

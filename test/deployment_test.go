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
			"generateSecret": "false",
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
			"registry":       "registry.neuvector.com",
			"generateSecret": "false",
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
			"registry":       "registry.neuvector.com",
			"oem":            "oem",
			"tag":            "0.9",
			"generateSecret": "false",
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
			"generateSecret":                "false",
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
			"generateSecret":              "false",
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
		SetValues: map[string]string{
			"generateSecret": "false",
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
			"generateSecret":             "false",
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
			"generateSecret":  "false",
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
			"generateSecret": "false",
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
			"generateSecret": "false",
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
		if secret.Name == "neuvector-controller-secret" {
			assert.NotNil(t, secret.Data)
			assert.NotEmpty(t, secret.Data["ssl-cert.key"])
			assert.NotEmpty(t, secret.Data["ssl-cert.pem"])
			assert.NotEmpty(t, secret.Data["jwt-signing.key"])
			assert.NotEmpty(t, secret.Data["jwt-signing.pem"])
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
		if dep.Name == "neuvector-controller-pod" {

			assert.Contains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "cert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "neuvector-controller-secret",
					},
				},
			})
		}
	}

	for _, output := range outs {
		var dep appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &dep)
		if dep.Name == "neuvector-controller-pod" {

			// cert, usercert and userjwtcert will be mounted.
			assert.Contains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "cert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "neuvector-controller-secret",
					},
				},
			})

			assert.NotContains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "usercert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "nv-ssl-secret",
					},
				},
			})

			assert.NotContains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "userjwtcert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "nv-jwt-secret",
					},
				},
			})
			for _, container := range dep.Spec.Template.Spec.Containers {
				if container.Name == "neuvector-controller-pod" {

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

					assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "cert",
						MountPath: "/etc/neuvector/certs/jwt-signing.key",
						SubPath:   "jwt-signing.key",
						ReadOnly:  true,
					})

					assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "cert",
						MountPath: "/etc/neuvector/certs/jwt-signing.pem",
						SubPath:   "jwt-signing.pem",
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
			"generateSecret": "false",
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
		assert.NotEqual(t, "neuvector-controller-secret", secret.Name)
	}

	for _, output := range outs {
		var dep appsv1.Deployment
		helm.UnmarshalK8SYaml(t, output, &dep)
		if dep.Name == "neuvector-controller-pod" {

			// cert, usercert and userjwtcert will be mounted.
			assert.NotContains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "cert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "neuvector-controller-secret",
					},
				},
			})

			assert.NotContains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "usercert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "nv-ssl-secret",
					},
				},
			})

			assert.NotContains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "userjwtcert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "nv-jwt-secret",
					},
				},
			})
			for _, container := range dep.Spec.Template.Spec.Containers {
				if container.Name == "neuvector-controller-pod" {

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

					assert.NotContains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "cert",
						MountPath: "/etc/neuvector/certs/jwt-signing.key",
						SubPath:   "jwt-signing.key",
						ReadOnly:  true,
					})

					assert.NotContains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "cert",
						MountPath: "/etc/neuvector/certs/jwt-signing.pem",
						SubPath:   "jwt-signing.pem",
						ReadOnly:  true,
					})
				}

			}

		}
	}
}

func TestControllerWithSSLAndJWTKeys(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"controller.certificate.secret":     "nv-ssl-secret",
			"controller.certificate.keyFile":    "key2.pem",
			"controller.certificate.pemFile":    "cert2.pem",
			"controller.jwtCertificate.secret":  "nv-jwt-secret",
			"controller.jwtCertificate.keyFile": "key2.pem",
			"controller.jwtCertificate.pemFile": "cert2.pem",
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
		if dep.Name == "neuvector-controller-pod" {

			// cert, usercert and userjwtcert will be mounted.
			assert.Contains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "cert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "neuvector-controller-secret",
					},
				},
			})

			assert.Contains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "usercert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "nv-ssl-secret",
					},
				},
			})

			assert.Contains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "userjwtcert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "nv-jwt-secret",
					},
				},
			})
			for _, container := range dep.Spec.Template.Spec.Containers {
				if container.Name == "neuvector-controller-pod" {

					assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "usercert",
						MountPath: "/etc/neuvector/certs/ssl-cert.key",
						SubPath:   "key2.pem",
						ReadOnly:  true,
					})

					assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "usercert",
						MountPath: "/etc/neuvector/certs/ssl-cert.pem",
						SubPath:   "cert2.pem",
						ReadOnly:  true,
					})

					assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "userjwtcert",
						MountPath: "/etc/neuvector/certs/jwt-signing.key",
						SubPath:   "key2.pem",
						ReadOnly:  true,
					})

					assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "userjwtcert",
						MountPath: "/etc/neuvector/certs/jwt-signing.pem",
						SubPath:   "cert2.pem",
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
			"controller.certificate.secret":  "nv-ssl-secret",
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
		if dep.Name == "neuvector-controller-pod" {

			// cert, usercert will be mounted but not userjwtcert.
			assert.Contains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "cert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "neuvector-controller-secret",
					},
				},
			})

			assert.Contains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "usercert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "nv-ssl-secret",
					},
				},
			})

			assert.NotContains(t, dep.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: "userjwtcert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "nv-jwt-secret",
					},
				},
			})

			for _, container := range dep.Spec.Template.Spec.Containers {
				if container.Name == "neuvector-controller-pod" {

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

					assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "cert",
						MountPath: "/etc/neuvector/certs/jwt-signing.key",
						SubPath:   "jwt-signing.key",
						ReadOnly:  true,
					})

					assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
						Name:      "cert",
						MountPath: "/etc/neuvector/certs/jwt-signing.pem",
						SubPath:   "jwt-signing.pem",
						ReadOnly:  true,
					})
				}

			}

		}
	}
}

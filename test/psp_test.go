package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
)

func TestPSP(t *testing.T) {
	helmChartPath := ".."

	options := &helm.Options{
		SetValues: map[string]string{
			"psp": "true",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/psp.yaml"})
	outs := splitYaml(out)

	if len(outs) != 3 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}

	for i, output := range outs {
		switch i {
		case 0:
			var psp extv1beta1.PodSecurityPolicy
			helm.UnmarshalK8SYaml(t, output, &psp)
			if psp.Name != "neuvector-binding-psp" {
				t.Errorf("PSP policy name is wrong. name=%v\n", psp.Name)
			}
		case 1:
			var role rbacv1.Role
			helm.UnmarshalK8SYaml(t, output, &role)
			if role.Name != "neuvector-binding-psp" {
				t.Errorf("PSP role name is wrong. name=%v\n", role.Name)
			}
			if len(role.Rules) != 1 || len(role.Rules[0].Resources) != 1 {
				t.Errorf("Unexpected PSP role. %v\n", role)
			}
		case 2:
			var rb rbacv1.RoleBinding
			helm.UnmarshalK8SYaml(t, output, &rb)
			if rb.Name != "neuvector-binding-psp" {
				t.Errorf("PSP rolebinding name is wrong. name=%v\n", rb.Name)
			}
			if len(rb.Subjects) != 1 || rb.Subjects[0].Kind != "ServiceAccount" || rb.Subjects[0].Name != "neuvector" {
				t.Errorf("Unexpected PSP rolebinding. %v\n", rb)
			}
		}
	}
}

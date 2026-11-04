package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	coordinationv1 "k8s.io/api/coordination/v1"
)

// Fields mutated by the cert upgrader must not be rendered, or Helm 4 server-side apply conflicts on upgrade.
func TestCertUpgraderHelm4Fields(t *testing.T) {
	helmChartPath := "../charts/core"

	cases := []struct {
		name        string
		values      map[string]string
		expectLease bool
	}{
		{
			name:        "default values",
			values:      map[string]string{},
			expectLease: false,
		},
		{
			name:        "lease enabled",
			values:      map[string]string{"lease.enabled": "true"},
			expectLease: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			options := &helm.Options{SetValues: c.values}

			out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/upgrader-cronjob.yaml"})
			var cronjob batchv1.CronJob
			helm.UnmarshalK8SYaml(t, out, &cronjob)
			require.Equal(t, "neuvector-cert-upgrader-pod", cronjob.Name)
			assert.NotContains(t, cronjob.Annotations, "cert-upgrader-uid")

			for _, tc := range []struct {
				template string
				name     string
			}{
				{template: "templates/controller-lease.yaml", name: "neuvector-controller"},
				{template: "templates/upgrader-lease.yaml", name: "neuvector-cert-upgrader"},
			} {
				out, err := helm.RenderTemplateE(t, options, helmChartPath, nvRel, []string{tc.template})
				if !c.expectLease {
					// helm template --show-only fails when the template renders nothing.
					require.Error(t, err, tc.template)
					assert.Contains(t, err.Error(), "could not find template")
					continue
				}
				require.NoError(t, err, tc.template)
				var lease coordinationv1.Lease
				helm.UnmarshalK8SYaml(t, out, &lease)
				assert.Equal(t, tc.name, lease.Name)
			}
		})
	}
}

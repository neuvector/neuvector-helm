package test

import (
	"maps"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
)

func TestRegistryAdapter(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"leastPrivilege":      "true",
			"cve.adapter.enabled": "true",
			"autoGenerateCert":    "false",
		},
	}

	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/registry-adapter.yaml"})
	outs := splitYaml(out)

	if len(outs) != 2 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

func TestRegistryAdapterDeployment(t *testing.T) {
	helmChartPath := "../charts/core"

	cases := []struct {
		name                   string
		values                 map[string]string
		expectedReplicas       int32
		expectedStrategyType   appsv1.DeploymentStrategyType
		expectedMaxSurge       string
		expectedMaxUnavailable string
	}{
		{
			name:                   "defaults",
			values:                 map[string]string{},
			expectedReplicas:       1,
			expectedStrategyType:   appsv1.RollingUpdateDeploymentStrategyType,
			expectedMaxSurge:       "1",
			expectedMaxUnavailable: "0",
		},
		{
			name: "custom replicas and rolling update",
			values: map[string]string{
				"cve.adapter.replicas":                              "3",
				"cve.adapter.strategy.rollingUpdate.maxSurge":       "2",
				"cve.adapter.strategy.rollingUpdate.maxUnavailable": "1",
			},
			expectedReplicas:       3,
			expectedStrategyType:   appsv1.RollingUpdateDeploymentStrategyType,
			expectedMaxSurge:       "2",
			expectedMaxUnavailable: "1",
		},
		{
			name: "zero replicas",
			values: map[string]string{
				"cve.adapter.replicas": "0",
			},
			expectedReplicas:       0,
			expectedStrategyType:   appsv1.RollingUpdateDeploymentStrategyType,
			expectedMaxSurge:       "1",
			expectedMaxUnavailable: "0",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			values := map[string]string{
				"cve.adapter.enabled": "true",
				"autoGenerateCert":    "false",
			}
			maps.Copy(values, c.values)
			options := &helm.Options{SetValues: values}

			out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/registry-adapter.yaml"})
			outs := splitYaml(out)
			require.Len(t, outs, 2)

			var deploy appsv1.Deployment
			for _, o := range outs {
				var obj appsv1.Deployment
				helm.UnmarshalK8SYaml(t, o, &obj)
				if obj.Kind == "Deployment" {
					deploy = obj
				}
			}
			require.Equal(t, "Deployment", deploy.Kind)

			require.NotNil(t, deploy.Spec.Replicas)
			assert.Equal(t, c.expectedReplicas, *deploy.Spec.Replicas)
			assert.Equal(t, c.expectedStrategyType, deploy.Spec.Strategy.Type)
			require.NotNil(t, deploy.Spec.Strategy.RollingUpdate)
			require.NotNil(t, deploy.Spec.Strategy.RollingUpdate.MaxSurge)
			require.NotNil(t, deploy.Spec.Strategy.RollingUpdate.MaxUnavailable)
			assert.Equal(t, c.expectedMaxSurge, deploy.Spec.Strategy.RollingUpdate.MaxSurge.String())
			assert.Equal(t, c.expectedMaxUnavailable, deploy.Spec.Strategy.RollingUpdate.MaxUnavailable.String())
		})
	}
}

func TestRegistryAdapterIngress(t *testing.T) {
	helmChartPath := "../charts/core"

	options := &helm.Options{
		SetValues: map[string]string{
			"leastPrivilege":              "true",
			"cve.adapter.enabled":         "true",
			"cve.adapter.ingress.enabled": "true",
			"openshift":                   "true",
			"autoGenerateCert":            "false",
		},
	}

	// Test ingress
	out := helm.RenderTemplate(t, options, helmChartPath, nvRel, []string{"templates/registry-adapter-ingress.yaml"})
	outs := splitYaml(out)

	if len(outs) != 2 {
		t.Errorf("Resource count is wrong. count=%v\n", len(outs))
	}
}

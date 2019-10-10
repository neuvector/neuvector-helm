package helm

import (
	"github.com/gruntwork-io/terratest/modules/k8s"
)

type Options struct {
	ValuesFiles    []string            // List of values files to render.
	SetValues      map[string]string   // Values that should be set via the command line.
	SetStrValues   map[string]string   // Values that should be set via the command line explicitly as `string` types.
	SetFiles       map[string]string   // Values that should be set from a file. These should be file paths. Use to avoid logging secrets.
	KubectlOptions *k8s.KubectlOptions // KubectlOptions to control how to authenticate to kubernetes cluster. `nil` => use defaults.
	HomePath       string              // The path to the helm home to use when calling out to helm. Empty string means use default ($HOME/.helm).
	EnvVars        map[string]string   // Environment variables to set when running helm
}

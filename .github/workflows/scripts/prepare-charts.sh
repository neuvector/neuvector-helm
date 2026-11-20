#!/usr/bin/env bash
# Prepares the NeuVector charts of a release for the ship-family-staged action of
# rancher/charts. It takes the charts that NeuVector releases, charts/core and charts/crd
# of a release branch, and makes from them the two charts of the Rancher catalog, in
# <output dir>:
#
#   neuvector          the chart, as Rancher installs it
#   neuvector-crd      the CRDs, as a chart of their own
#   neuvector-monitor  the add-on that exports the metrics for Prometheus
#
# The action stamps, checks, packages and pushes the copies, so the source stays as it is.
#
# What makes the Rancher charts is in .github/rancher-chart, for each of core and monitor:
#   patches/  the changes to the templates that Rancher needs (the image registry,
#             the service account and the pull secrets), as patches that do not depend
#             on the version of the release
#   overlay/  files that are only in the Rancher chart: the questions of the UI, its
#             readme, and, for core, the check of the pod security policies
# and in this script, which edits Chart.yaml and values.yaml with yq, because those two
# change with every release: the name, the annotations, the images, and what Rancher
# switches on.
#
# When a patch no longer applies to a release, the script stops and names it: the
# template changed upstream, and the patch has to be made again against it.
#
# Usage: prepare-charts.sh <release checkout> <output dir>
#   e.g. prepare-charts.sh release "$RUNNER_TEMP/charts"
#
# Environment (both required):
#   RANCHER_VERSION  the Rancher minor the charts are for, such as 2.16
#   KUBE_VERSION     the Kubernetes versions the release supports, with both bounds, such
#                    as ">= 1.34.0-0 < 1.37.0-0"
#
# Requires yq (mikefarah, v4) and patch on PATH.
set -euo pipefail

REL=${1:?usage: $0 <release checkout> <output dir>}
OUT=${2:?usage: $0 <release checkout> <output dir>}
RANCHER_VERSION=${RANCHER_VERSION:?RANCHER_VERSION is required, such as 2.16}
KUBE_VERSION=${KUBE_VERSION:?KUBE_VERSION is required, such as \">= 1.34.0-0 < 1.37.0-0\"}
HERE=$(cd "$(dirname "$0")" && pwd)
RANCHER_CHART="$HERE/../../rancher-chart"

if ! [[ "$RANCHER_VERSION" =~ ^2\.([0-9]+)$ ]]; then
  echo "Error: RANCHER_VERSION $RANCHER_VERSION must look like 2.16" >&2
  exit 1
fi
RANCHER_MINOR="${BASH_REMATCH[1]}"
RANCHER_RANGE=">= 2.${RANCHER_MINOR}.0-0 < 2.$((RANCHER_MINOR + 1)).0-0"

if ! [[ "$KUBE_VERSION" =~ ^\>=\ [0-9]+\.[0-9]+\.[0-9]+-0\ \<\ [0-9]+\.[0-9]+\.[0-9]+-0$ ]]; then
  echo "Error: KUBE_VERSION '$KUBE_VERSION' must look like \">= 1.34.0-0 < 1.37.0-0\"" >&2
  exit 1
fi

for file in charts/core/Chart.yaml charts/core/values.yaml charts/crd/Chart.yaml charts/crd/templates/crd.yaml charts/crd/templates/_helpers.tpl charts/monitor/Chart.yaml charts/monitor/values.yaml; do
  if [ ! -f "$REL/$file" ]; then
    echo "Error: there is no $file in $REL" >&2
    exit 1
  fi
done
CHART_VERSION=$(yq '.version' "$REL/charts/core/Chart.yaml")
APP_VERSION=$(yq '.appVersion' "$REL/charts/core/Chart.yaml")

rm -rf "$OUT/neuvector" "$OUT/neuvector-crd" "$OUT/neuvector-monitor"
mkdir -p "$OUT"
cp -r "$REL/charts/core" "$OUT/neuvector"

# What the Rancher chart does not have: the CRDs, which are a chart of their own, the
# templates of the cloud marketplace (CSP) installs, and the schema of the values.
rm -f "$OUT"/neuvector/templates/{crd,csp-clusterrole,csp-clusterrolebinding,csp-crd,csp-deployment,csp-role,csp-rolebinding,csp-serviceaccount}.yaml \
  "$OUT/neuvector/values.schema.json"

# apply_rancher_files <core|monitor> <chart dir>: the patches of the templates, and the
# files that only the Rancher chart has.
apply_rancher_files() {
  local patch
  while IFS= read -r patch; do
    if ! patch -p1 -s -d "$2" <"$patch"; then
      echo "Error: ${patch#"$HERE"/../../} does not apply to the release: its template changed, and the patch has to be made again" >&2
      exit 1
    fi
  done < <(find "$RANCHER_CHART/$1/patches" -name '*.patch' | sort)
  cp -r "$RANCHER_CHART/$1/overlay/." "$2/"
}

apply_rancher_files core "$OUT/neuvector"

# Chart.yaml: the name in the catalog and the annotations that it reads.
RANGE="$RANCHER_RANGE" KUBE="$KUBE_VERSION" UPSTREAM="$CHART_VERSION" yq -i '
  .name = "neuvector" |
  .description = "Helm feature chart for NeuVector container security platform." |
  .keywords = ["security"] |
  del(.sources[] | select(. == "https://github.com/neuvector/neuvector-helm")) |
  .annotations."catalog.cattle.io/auto-install" = "neuvector-crd=match" |
  .annotations."catalog.cattle.io/certified" = "rancher" |
  .annotations."catalog.cattle.io/display-name" = "NeuVector" |
  .annotations."catalog.cattle.io/kube-version" = strenv(KUBE) |
  .annotations."catalog.cattle.io/namespace" = "cattle-neuvector-system" |
  .annotations."catalog.cattle.io/os" = "linux" |
  .annotations."catalog.cattle.io/permits-os" = "linux" |
  .annotations."catalog.cattle.io/provides-gvr" = "neuvector.com/v1" |
  .annotations."catalog.cattle.io/rancher-version" = strenv(RANGE) |
  .annotations."catalog.cattle.io/release-name" = "neuvector" |
  .annotations."catalog.cattle.io/type" = "cluster-tool" |
  .annotations."catalog.cattle.io/upstream-version" = strenv(UPSTREAM)
' "$OUT/neuvector/Chart.yaml"

# values.yaml: what Rancher sets for its own installs, the images that SUSE publishes as
# rancher/neuvector-<image>, and no marketplace values. The tag of the NeuVector images
# is the one of the release.
yq -i '
  del(.imagePullSecrets) |
  del(.psp) |
  .serviceAccount = "neuvector" |
  .global.cattle.imagePullSecrets = [] |
  .global.cattle.systemDefaultRegistry = "" |
  .global.cattle.psp.enabled = false |
  del(.global.azure) |
  del(.global.aws) |
  .controller.image.repository = "rancher/neuvector-controller" |
  .controller.image.tag = .tag |
  .controller.ranchersso.enabled = true |
  .controller.prime.enabled = true |
  .controller.prime.image.repository = "rancher/neuvector-compliance-config" |
  .enforcer.image.repository = "rancher/neuvector-enforcer" |
  .enforcer.image.tag = .tag |
  .manager.image.repository = "rancher/neuvector-manager" |
  .manager.image.tag = .tag |
  .cve.adapter.image.repository = "rancher/neuvector-registry-adapter" |
  .cve.updater.image.repository = "rancher/neuvector-updater" |
  .cve.scanner.image.repository = "rancher/neuvector-scanner"
' "$OUT/neuvector/values.yaml"

# The CRD chart: the CRDs of the release, without the label that Helm adds to every
# object, behind the switch that Rancher uses for them, and the helpers of their chart.
mkdir -p "$OUT/neuvector-crd/templates"
{
  echo "{{- if .Values.crdwebhook.enabled -}}"
  grep -v 'heritage: {{ .Release.Service }}' "$REL/charts/crd/templates/crd.yaml"
  echo "{{- end }}"
} >"$OUT/neuvector-crd/templates/crd.yaml"
cp "$REL/charts/crd/templates/_helpers.tpl" "$OUT/neuvector-crd/templates/_helpers.tpl"
cat >"$OUT/neuvector-crd/values.yaml" <<'YAML'
# Default values for neuvector.
# This is a YAML-formatted file.
# Declare variables to be passed into the templates.

openshift: false

crdwebhook:
  type: ClusterIP
  enabled: true
YAML
cat >"$OUT/neuvector-crd/Chart.yaml" <<YAML
apiVersion: v2
name: neuvector-crd
description: Helm chart for NeuVector's CRD services
version: ${CHART_VERSION}
appVersion: ${APP_VERSION}
type: application
home: https://neuvector.com
icon: https://avatars2.githubusercontent.com/u/19367275?s=200&v=4
maintainers:
  - name: becitsthere
    email: support@neuvector.com
annotations:
  catalog.cattle.io/certified: rancher
  catalog.cattle.io/hidden: "true"
  catalog.cattle.io/namespace: cattle-neuvector-system
  catalog.cattle.io/release-name: neuvector-crd
YAML

# The monitor: the exporter of the metrics. It takes the images and the pull secrets the
# way the main chart does, from the patches.
cp -r "$REL/charts/monitor" "$OUT/neuvector-monitor"
apply_rancher_files monitor "$OUT/neuvector-monitor"
RANGE="$RANCHER_RANGE" KUBE="$KUBE_VERSION" UPSTREAM="$(yq '.version' "$REL/charts/monitor/Chart.yaml")" yq -i '
  .name = "neuvector-monitor" |
  .description = "Helm feature chart (optional) add-on to NeuVector for monitoring with Prometheus/Grafana." |
  .keywords = ["security"] |
  .sources = ["https://github.com/neuvector/neuvector"] |
  .annotations."catalog.cattle.io/certified" = "rancher" |
  .annotations."catalog.cattle.io/display-name" = "NeuVector Monitor" |
  .annotations."catalog.cattle.io/kube-version" = strenv(KUBE) |
  .annotations."catalog.cattle.io/namespace" = "cattle-neuvector-system" |
  .annotations."catalog.cattle.io/os" = "linux" |
  .annotations."catalog.cattle.io/permits-os" = "linux" |
  .annotations."catalog.cattle.io/provides-gvr" = "neuvector.com/v1" |
  .annotations."catalog.cattle.io/rancher-version" = strenv(RANGE) |
  .annotations."catalog.cattle.io/release-name" = "neuvector-monitor" |
  .annotations."catalog.cattle.io/type" = "cluster-tool" |
  .annotations."catalog.cattle.io/upstream-version" = strenv(UPSTREAM)
' "$OUT/neuvector-monitor/Chart.yaml"
yq -i '
  .global.cattle.systemDefaultRegistry = "" |
  .global.cattle.imagePullSecrets = [] |
  .exporter.image.repository = "rancher/neuvector-prometheus-exporter"
' "$OUT/neuvector-monitor/values.yaml"

# Usage
# ./scripts/bump_version.sh <chart version> <app version>

sed -i "s/^version:.*/version: $1/g" charts/core/Chart.yaml
sed -i "s/^appVersion:.*/appVersion: $2/g" charts/core/Chart.yaml
sed -i "s/^version:.*/version: $1/g" charts/crd/Chart.yaml
sed -i "s/^appVersion:.*/appVersion: $2/g" charts/crd/Chart.yaml
sed -i "s/^version:.*/version: $1/g" charts/monitor/Chart.yaml
sed -i "s/^tag:.*/tag: $2/g" charts/core/values.yaml

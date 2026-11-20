# The Rancher charts

What the Rancher catalog has that the charts of this repository do not, used by
`.github/workflows/scripts/prepare-charts.sh` to make `neuvector`, `neuvector-crd` and
`neuvector-monitor` from `charts/core`, `charts/crd` and `charts/monitor`.

For each of `core` and `monitor`:

- `patches/` changes the templates so that the images come from the registry Rancher sets
  (`global.cattle.systemDefaultRegistry`) and the pull secrets from `global.cattle`. They
  are patches of the templates, which do not change with the version of a release.
- `overlay/` has the files that only the Rancher chart has: the questions of the UI, its
  readme and, for `core`, the check of the pod security policies.

The changes to `Chart.yaml` and `values.yaml`, which differ in every release, are in the
script, as `yq` edits. When a template changes in a release and a patch no longer applies,
the script stops and names it: make the patch again against that template.

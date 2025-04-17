# NeuVector Override files

## Scan caching
Scan caching can be enabled by editing values.yaml or creating below override file and pass them with "-f" option on HELM commands.
```console
cve:
  scanner:
    volumes:
      - name: scan-cache
        hostPath:
          path: /tmp/
          type: ""
    volumeMounts:
      - mountPath: /tmp/images/caches
        name: scan-cache
```

## Google Autopilot support
Below override files should be used for the deploying NeuVector on GKE Autopilot cluster.
```console
cve:
  scanner:
    podLabels:
      cloud.google.com/matching-allowlist: suse-neuvector-scanner-v1.0
    resources:
      limits:
        ephemeral-storage: "3Gi"
      requests:
        ephemeral-storage: "2Gi"
enforcer:
  podLabels:
    cloud.google.com/matching-allowlist: suse-neuvector-enforcer-v1.0
```

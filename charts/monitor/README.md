# NeuVector Helm Chart

Helm chart for NeuVector's monitoring services.

## Configuration

The following table lists the configurable parameters of the NeuVector chart and their default values.

Parameter | Description | Default | Notes
--------- | ----------- | ------- | -----
`exporter.enabled` | If true, create Prometheus exporter | `false` |
`exporter.image.repository` | exporter image name | `neuvector/prometheus-exporter` |
`exporter.image.tag` | exporter image tag | `latest` |
`exporter.CTRL_USERNAME` | Username to login to the controller. Suggest to replace the default admin user to a read-only user | `admin` |
`exporter.CTRL_PASSWORD` | Passowrd to login to the controller. | `admin` |

---
Contact <support@neuvector.com> for access to Docker Hub and docs.


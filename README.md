# NeuVector

Visibility and Security: The NeuVector ‘Multi-Vector Container Security Platform’

[NeuVector](https://neuvector.com) provides a real-time Kubernetes and OpenShift container security solution that adapts easily to your changing environment and secures containers at their most vulnerable point – during run-time. The declarative security policy ensures that applications scale up or scale down quickly without manual intervention. The NeuVector solution is a Red Hat and Docker Certified container itself which deploys easily on each host, providing a container firewall, host monitoring and security, security auditing with CIS benchmarks, and vulnerability scanning.

The installation will deploy the NeuVector Enforcer container on each worker node as a daemon set, and by default 3 controller containers (for HA, one is elected the leader). The controllers can be deployed on any node, including Master, Infra or management nodes. See the NeuVector docs for node labeling to control where controllers are deployed.

## Prerequisites

- Kubernetes 1.7+

- If you are going to pull images from docker.io and need an image pull secret:

```console
$ kubectl create secret docker-registry regsecret -n neuvector --docker-server=https://index.docker.io/v1/ --docker-username=your-name --docker-password=your-pword --docker-email=your-email
```

## Downloading the Chart

Clone or download this repository.

## Installing the Chart

To install the chart with the release name `my-release` and image pull secret:

```console
$ helm install --name my-release --namespace neuvector ./neuvector-helm/k8s/ --set imagePullSecrets=regsecret
```

If you already pulled neuvector images and saved in your private registry:

```console
$ helm install --name my-release --namespace neuvector ./neuvector-helm/k8s/ --set registry=your-private-registry
```

## Openshift

Replace k8s with openshift, for example:

```console
$ helm install --name my-release --namespace neuvector ./neuvector-helm/openshift/ --set registry=your-private-registry
```

## Rolling upgrade

```console
$ helm upgrade my-release --set tag=2.2.0 ./neuvector-helm/k8s/
```

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the NeuVector chart and their default values.

Parameter | Description | Default
--------- | ----------- | -------
`registry` | image registry | `docker.io`
`tag` | image tag for controller enforcer manager | `latest`
`imagePullSecrets` | image pull secret | `{}`
`controller.enabled` | If true, create controller | `true`
`controller.image.repository` | controller image repository | `neuvector/controller`
`controller.replicas` | controller replicas | `3`
`enforcer.enabled` | If true, create enforcer | `true`
`enforcer.image.repository` | enforcer image repository | `neuvector/enforcer`
`manager.enabled` | If true, create manager | `true`
`manager.image.repository` | manager image repository | `neuvector/manager`
`manager.env.ssl` | enable/disable HTTPS and disable/enable HTTP access  | `on`
`manager.svc.type` | manager service type | `NodePort`
`updater.enabled` | If true, create updater | `true`
`updater.image.repository` | updater image repository | `neuvector/updater`
`updater.image.tag` | image tag for updater | `latest`
`updater.schedule` | cronjob updater schedule | `0 0 * * *` |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```console
$ helm install --name my-release --namespace neuvector ./neuvector-helm/k8s/ --set manager.env.ssl=off
```

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example,

```console
$ helm install --name my-release --namespace neuvector ./neuvector-helm/k8s/ -f values.yaml
```

> **Tip**: You can use the default [values.yaml](k8s/values.yaml)


If you installed neuvector before and manually created the cluster role neuvector-binding, you need to delete this cluster role first. If helm install error because of this, you need to `helm del --purge my-release` before install again.


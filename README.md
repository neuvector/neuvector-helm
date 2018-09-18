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
$ helm install --name my-release --namespace neuvector ./neuvector-helm/ --set imagePullSecrets=regsecret
```

If you already pulled neuvector images and saved in your private registry:

```console
$ helm install --name my-release --namespace neuvector ./neuvector-helm/ --set registry=your-private-registry
```

If you already installed neuvector in your cluster without using helm, please `kubectl delete -f your-neuvector-yaml.yaml` before trying to use helm install.

## Openshift

```console
$ helm install --name my-release --namespace neuvector ./neuvector-helm/ --set openshift=true,registry=your-private-registry
```

## Rolling upgrade

```console
$ helm upgrade my-release --set tag=2.2.0 ./neuvector-helm/
```

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the NeuVector chart and their default values.

Parameter | Description | Default | Notes
--------- | ----------- | ------- | -----
`openshift` | If deploying in Openshift, set this to true | `false` | 
`registry` | image registry | `docker.io` | If Azure, set to my-reg.azurecr.io;<br>if Openshift, set to docker-registry.default.svc:5000
`tag` | image tag for controller enforcer manager | `latest` | 
`imagePullSecrets` | image pull secret | `{}` | 
`controller.enabled` | If true, create controller | `true` | 
`controller.image.repository` | controller image repository | `neuvector/controller` | 
`controller.replicas` | controller replicas | `3` | 
`controller.pvc.enabled` | If true, enable persistence for controller using PVC | `false` | Require persistent volume type RWX, and storage 1Gi
`enforcer.enabled` | If true, create enforcer | `true` | 
`enforcer.image.repository` | enforcer image repository | `neuvector/enforcer` | 
`manager.enabled` | If true, create manager | `true` | 
`manager.image.repository` | manager image repository | `neuvector/manager` | 
`manager.env.ssl` | enable/disable HTTPS and disable/enable HTTP access  | `on`;<br>if ingress is enabled, then default is `off` | 
`manager.svc.type` | set manager service type for native Kubernetes | `NodePort`;<br>if it is Openshift platform or ingress is enabled, then default is `ClusterIP` | set to LoadBalancer if using cloud providers, such as Azure, Amazon, Google
`manager.ingress.enabled` | If true, create ingress, must also set ingress host value | `false` | enable this if ingress controller is installed
`manager.ingress.host` | Must set this host value if ingress is enabled | `{}` | 
`cve.updater.enabled` | If true, create cve updater | `true` | 
`cve.updater.image.repository` | cve updater image repository | `neuvector/updater` | 
`cve.updater.image.tag` | image tag for cve updater | `latest` | 
`cve.updater.schedule` | cronjob cve updater schedule | `0 0 * * *` |  |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```console
$ helm install --name my-release --namespace neuvector ./neuvector-helm/ --set manager.env.ssl=off
```

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example,

```console
$ helm install --name my-release --namespace neuvector ./neuvector-helm/ -f values.yaml
```

> **Tip**: You can use the default [values.yaml](values.yaml)

## RBAC Configuration

If you installed neuvector before and manually created the cluster role and cluster role binding for neuvector-binding, you need to delete the cluster role binding first, then delete the cluster role.

```console
$ kubectl delete clusterrolebinding neuvector-binding
$ kubectl delete clusterrole neuvector-binding
```

If helm install returns error because of an existing cluster role, you need to delete the release before install again.

```console
$ helm delete --purge my-release
```


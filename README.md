# NeuVector Helm charts

A collection of Helm charts for deploying NeuVector product in Kubernetes and Openshift clusters.

## Installing charts

### Helm Charts

This repository contains three Helm charts.
Chart | Description
----- | -----------
core | Deploy NeuVector container security core services. [chart](charts/core)
crd | Deploy CRD services before installing NeuVector container security platform. [chart](charts/crd)
monitor | Deploy monitoring services, such as Prometheus exporter. [chart](charts/monitor)

**IMPORTANT** - Each chart has a set of configuration values, especially for the 'core' chart. Review the Helm chart configuration values [here](charts/core) and make any required changes to the values.yaml file for your deployment.

### Adding chart repo

```console
$ helm repo add neuvector https://neuvector.github.io/neuvector-helm/
$ helm search repo neuvector/core
```

### Versioning

Helm charts for officially released product are published from the release branch of the repository. The main branch is used for the charts of the product in the development. Typically the charts in the main branch are published with the alpha, beta or rc tag. They can be discovered with --devel option.

```console
$ helm search repo neuvector/core 
NAME          	CHART VERSION	APP VERSION	DESCRIPTION
neuvector/core	1.9.2        	4.4.4-s2   	Helm chart for NeuVector's core services

$ helm search repo becitsthere/core --devel
NAME            	CHART VERSION	APP VERSION	DESCRIPTION
neuvector/core	2.2.0-b1     	5.0.0-b1   	Helm chart for NeuVector's core services
neuvector/core	1.9.2        	4.4.4-s2   	Helm chart for NeuVector's core services
```

#### Kubernetes

- Create the NeuVector namespace.
```console
$ kubectl create namespace neuvector
```

- Configure Kubernetes to pull from the NeuVector container registry.
```console
$ kubectl create secret docker-registry regsecret -n neuvector --docker-server=https://index.docker.io/v1/ --docker-username=your-name --docker-password=your-password --docker-email=your-email
```

Where ’your-name’ is your registry username, ’your-password’ is your registry password, ’your-email’ is your email.

To install the chart with the release name `my-release` and image pull secret:

```console
$ helm install my-release --namespace neuvector neuvector/core  --set imagePullSecrets=regsecret
```

#### RedHat OpenShift

- Create a new project.
```console
$ oc new-project neuvector
```

- Create a new service account **if** you don't want to use the 'default'. Specify the service account name in charts' values.yaml file.
```console
$ oc create serviceaccount neuvector -n neuvector
```

- Grant Service Account Access to the Privileged SCC. Please replace the service account name that you plan to use.
```console
$ oc -n neuvector adm policy add-scc-to-user privileged -z default
```

- Configure Openshift to pull from the NeuVector container registry.
```console
$ oc create secret docker-registry regsecret -n neuvector --docker-server=https://index.docker.io/v1/ --docker-username=your-name --docker-password=your-password --docker-email=your-email
```

To install the chart with the release name `my-release`:

```console
$ helm install my-release --namespace neuvector neuvector/core --set openshift=true,imagePullSecrets=regsecret,crio.enabled=true
```

To install the chart with the release name `my-release` and your private registry:

```console
$ helm install my-release --namespace neuvector neuvector/core --set openshift=true,imagePullSecrets=regsecret,crio.enabled=true,registry=your-private-registry
```

If you are using a private registry, and want to enable the updater cronjob, please create a script, run it as a cronjob before midnight or the updater daily schedule.

## Rolling upgrade

```console
$ helm upgrade my-release --set imagePullSecrets=regsecret,tag=4.4.0 neuvector/core
```

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Using private registry

If you are using a private registry, you need pull NeuVector images of the specified version to your own registry and add registry name when installing the chart.

```console
$ helm install my-release --namespace neuvector neuvector/core --set registry=your-private-registry
```

To keep the vulnerability database up-to-date, you want to create a script, run it as a cronjob to pull the updater and scanner images periodically to your own registry.

```console
$ docker login docker.io
$ docker pull docker.io/neuvector/updater
$ docker logout docker.io

$ oc login -u <user_name>
# this user_name is the one when you install neuvector

$ docker login -u <user_name> -p `oc whoami -t` docker-registry.default.svc:5000
$ docker tag docker.io/neuvector/updater docker-registry.default.svc:5000/neuvector/updater
$ docker push docker-registry.default.svc:5000/neuvector/updater
$ docker logout docker-registry.default.svc:5000
```

## Migration

If you are using the previous way to install charts from the source directly, after adding the Helm repo, you can upgrade the current installation by given the same chart name. 

```console
helm upgrade my-release neuvector/core --namespace neuvector --set tag=4.1.0
```

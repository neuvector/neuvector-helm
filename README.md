# NeuVector Helm charts

A collection of Helm charts for deploying NeuVector product in Kubernetes and Openshift clusters.

## Installing charts

### Adding chart repo

```console
$ helm repo add neuvector https://neuvector.github.io/neuvector-helm/
$ helm search neuvector/core
```

#### Kubernetes

- Create the NeuVector namespace.
```console
$ kubectl create namespace neuvector
```

- Configure Kubernetes to pull from the private NeuVector registry on Docker Hub.
```console
$ kubectl create secret docker-registry regsecret -n neuvector --docker-server=https://index.docker.io/v1/ --docker-username=your-name --docker-password=your-pword --docker-email=your-email
```

Where ’your-name’ is your Docker username, ’your-pword’ is your Docker password, ’your-email’ is your Docker email.

To install the chart with the release name `my-release` and image pull secret:

```console
$ helm install --name my-release --namespace neuvector neuvector/core  --set imagePullSecrets=regsecret
```

#### RedHat OpenShift

- Create a new project.
```console
$ oc new-project neuvector
```

- Grant Service Account Access to the Privileged SCC.
```console
$ oc -n neuvector adm policy add-scc-to-user privileged -z neuvector
```

To install the chart with the release name `my-release` and your private registry:

```console
$ helm install --name my-release --namespace neuvector neuvector/core --set openshift=true,registry=your-private-registry
```

If you are using a private registry, and want to enable the updater cronjob, please create a script, run it as a cronjob before midnight or the updater daily schedule.

## Rolling upgrade

```console
$ helm upgrade my-release --set imagePullSecrets=regsecret,tag=4.0.0 neuvector/core
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
$ helm install --name my-release --namespace neuvector neuvector/core --set registry=your-private-registry
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
helm upgrade my-release neuvector/core --namespace neuvector --set tag=4.0.0.s1
```

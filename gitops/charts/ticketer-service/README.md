# user-service
user-service is a node.js socket.io server with basic messaging functionality 

## Installing the Chart

To install the chart with the release name `my-release`:

```console
$ helm upgrade my-release --install user-service 
```

The command deploys user-service on the Kubernetes cluster in the default namespace.

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
$ helm delete --purge my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.
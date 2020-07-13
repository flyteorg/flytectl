# flytectl
Flyte CLI

In order to dev, needing to work with a cluster.

A local one can be setup via --> https://lyft.github.io/flyte/administrator/install/getting_started.html#getting-started


get the gRPC port:
```
FLYTECTL_GRPC_PORT=`kubectl get service -n flyte flyteadmin -o json | jq '.spec.ports[] | select(.name=="grpc").port'`
```

`kubectl port-forward -n flyte service/flyteadmin 8081:$FLYTECTL_GRPC_PORT`

Update config line in https://github.com/lyft/flytectl/blob/master/config.yaml
to dns:///localhost:8081 


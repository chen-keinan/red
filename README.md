## Runtime Dev Cli
Simplify development via IDE with runtime components running in cluster (app-proxy / gitops-operator)
- Setup ngrok tunnels from cluster to local env
- Port forward fro cluster to local env
- Update relevent configmaps with tunnels external ip
- Create env variable files for gitops-operator and app-proxy to be used with IDE

## Getting started

### pre-requisite
- running k8s cluster
- make sure gitops-runtime is deployed to you cluster and its running without any errors
- make sure values.yaml for that installation is avaliable on specific folder
- make sure the `env.sh` script (extract runtime values) is avaliable on specific folder

```sh
git clone https://github.com/chen-keinan/devcli
cd devcli
make build
```

## Basic usage
```sh
./devcli
```

```sh
devcli
Command Options:
-- clean      Clean up resources and delete DevEnv files
-- setup      Setting up app-proxy and gitops-operator DevEnv
```

## Setup - Dev Env
```sh
./devcli --setup

***************************************************************************************************************************

1. Enter Helm Values Path (default: /path/to/runtime/values.yaml):
2. Enter Codefresh Namespace (default: codefresh):
3. Enter Cluster Name (default: kind-local-cluster):
4. Enter Environment Variable Script Path (default: /path/to/env/extarct/script/env.sh):
5. Enter debug-app-proxy (default: y):
6. Enter debug-gitops-operator (default: y):

****************************************************************************************************************************

- Reading Helm Values
- Extracting Values from EnvVar script
- Tunneling 3017 --> Localhost
Forwarding from 127.0.0.1:2746 -> 2746
Forwarding from [::1]:2746 -> 2746
Forwarding from 127.0.0.1:8080 -> 8080
Forwarding from [::1]:8080 -> 8080
- Updating codefresh-cm
- Tunneling 8082 --> Localhost
- Scalling down gitops operator to 0
- Updating gitops-operator-notifications cm
********************************************************
-- output files:
/Users/<name>/.devcli/app-proxy-dev-env.json
/Users/<name>/.devcli/gitops-operator-dev-env.json

******************************************************
```

copy the env var values from the `output files` and put it in your ide (app-proxy and gitops-operator launch setting )

## Cleanup -  Dev Env
```sh
./devcli --clean
- Revert codefresh-cm ingressUrl value
- Revert gitops-operator-notifications-cm gitops-operator value
- Clean up ngrok tunnels
- Clean up port forwards
- Clean up output folder: /Users/<name>/.devcli
```
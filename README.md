## Runtime Dev Env CLI
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
git clone https://github.com/chen-keinan/red
cd red
make install
```

## Basic usage
```sh
red
```

```sh
Runtime Env Dev
Command Options:
--clean      Clean up resources and delete DevEnv files
--setup      Setting up app-proxy and gitops-operator DevEnv
--no-setup   loading setting from /Users/<UserName>/red.json (this option is not valid on 1st setup)
```

## Setup - Dev Env
```sh
red --setup
***************************************************************************************************************************

1. Enter Helm Values Path (default: /Users/<UserName>/workspace/codefresh-values/local.values.yaml):
2. Enter Codefresh Namespace (default: codefresh):
3. Enter Cluster Name (default: kind-codefresh-local-cluster):
4. Enter Environment Variable Script Path (default: /Users/<UserName>/workspace/codefresh-values/env.sh):
5. Enter debug-app-proxy (default: y):
6. Enter debug-gitops-operator (default: y):

****************************************************************************************************************************

- Reading Helm Values
- Extracting Values from EnvVar script
- Tunneling 3017 --> Localhost
- Updating codefresh-cm
- Tunneling 8082 --> Localhost
- Updating gitops-operator-notifications cm
********************************************************
-- Copy the EnvVars values from output files to IDE run setting:

/Users/chenkeinan/.red/app-proxy-dev-env.json
/Users/chenkeinan/.red/gitops-operator-dev-env.json

******************************************************
port forward on ports:
 2746:2746
8080:8080

Enjoy Debugging :)
press Ctrl-c to terminate
```

copy the env var values from the `output files` and put it in your IDE (app-proxy and gitops-operator launch setting)

## Cleanup -  Dev Env
```sh
red --clean
- Revert codefresh-cm configmap
- Revert gitops-operator-notifications-cm configmap
- Clean up ngrok tunnels
- Clean up port forwards
- Clean up output folder: /Users/<name>/.red
```
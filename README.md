# kubectl-mounts
English|[简体中文](./README.zh.md) 

`kubectl-mounts` is a `kubectl` plugin written in Go that displays detailed information about Pod volumes and their mount paths in the current Kubernetes namespace.

## Features

- Displays Pod name, container name, volume name, mount path, and volume type
- Supports ConfigMap, Secret, EmptyDir, and PVC volume types
- Output rendered as a neat, human-readable table using `tablewriter`
- Automatically detects current context namespace

## Usage

```bash
Show Pod Volumes and VolumeMounts in the cluster

Usage:
  kubectl-mounts [flags]
  kubectl-mounts [command]

Available Commands:
  completion  Generate completion script
  help        Help about any command

Flags:
  -h, --help                help for kubectl-mounts
  -k, --kubeconfig string   Path to kubeconfig file (default $HOME/.kube/config)
  -n, --namespace string    Namespace (default is current namespace)
  -o, --output string       Output format: table|yaml|json(default table)
  -p, --pod string          Filter by specific Pod name
  -v, --version             Show version

Use "kubectl-mounts [command] --help" for more information about a command.

```
Example output:

```pgsql

+-----------------------------------------+-----------------------------------+-----------------------------------+-----------------------------------------------+--------------------------------------------------+
|                   Pod                   |             Container             |            Volume Name            |                   MountPath                   |                   Volume Type                    |
+-----------------------------------------+-----------------------------------+-----------------------------------+-----------------------------------------------+--------------------------------------------------+
| default-flyflow-console-8d97c74f-mxmgp  | default-flyflow-console-container | default-flyflow-console-nginxconf | /conf                                         | ConfigMap(cm-default-flyflow-console-nagix.conf) |
+                                         +                                   +-----------------------------------+-----------------------------------------------+--------------------------------------------------+
|                                         |                                   | kube-api-access-gtfww             | /var/run/secrets/kubernetes.io/serviceaccount | Projected                                        |
+-----------------------------------------+-----------------------------------+-----------------------------------+                                               +                                                  +
| ephemeral-demo                          | ephemeral-demo                    | kube-api-access-spnk2             |                                               |                                                  |
+-----------------------------------------+-----------------------------------+-----------------------------------+-----------------------------------------------+--------------------------------------------------+                                         
```
## Installation
### Homebrew

If you use Homebrew you can install like this:
```bash
brew tap yeqianmen/kubectl-mounts
brew install kubectl-mounts
```

### Using Krew (kubectl Plugin Manager)

If you use kubectl and have Krew installed, you can install Kelper as a kubectl plugin:
```bash
kubectl krew install --manifest-url https://github.com/yeqianmen/kubectl-mounts/releases/download/v0.1.1/mounts.yaml
```


### Build manually

```bash
git clone https://github.com/yeqianmen/kubectl-mounts.git
cd kubectl-mounts
go build -o kubectl-mounts
sudo mv kubectl-mounts /usr/local/bin/
```

## Completion

Supported Completions
- -n: Auto-completes available namespaces in the cluster

- -p: Auto-completes pod names in the selected namespace

- -o: Supports table, yaml, and json formats
 
### Bash

Load completion for current session:
```bash
source <(kubectl-mounts completion bash)
```
### Zsh

Load completion for current session:
```bash
source <(kubectl-mounts completion zsh)
```


# kubectl-mounts

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

Or build manually:

```bash
git clone https://github.com/yeqianmen/kubectl-mounts.git
cd kubectl-mounts
go build -o kubectl-mounts
sudo mv kubectl-mounts /usr/local/bin/
```
## Completion

Supported Completions
- --namespace: Auto-completes available namespaces in the cluster

- --pod: Auto-completes pod names in the selected namespace

- --output: Supports table, yaml, and json formats
 
**Bash**

Load completion for current session:
```bash
source <(kubectl-mounts completion bash)
```
**Zsh**

Load completion for current session:
```bash
source <(kubectl-mounts completion zsh)
```


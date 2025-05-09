# kubectl-mounts

`kubectl-mounts` is a `kubectl` plugin written in Go that displays detailed information about Pod volumes and their mount paths in the current Kubernetes namespace.

## Features

- Displays Pod name, container name, volume name, mount path, and volume type
- Supports ConfigMap, Secret, EmptyDir, and PVC volume types
- Output rendered as a neat, human-readable table using `tablewriter`
- Automatically detects current context namespace

## Usage

```bash
kubectl mounts



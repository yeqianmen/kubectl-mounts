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
```
Example output:

```pgsql

+------------+-----------+--------------+----------------------+---------------------+
| POD        | CONTAINER | VOLUME NAME  | MOUNT PATH           | VOLUME TYPE         |
+------------+-----------+--------------+----------------------+---------------------+
| mypod-1234 | app       | config-vol   | /etc/config          | ConfigMap(app-cm)   |
| mypod-1234 | app       | log-vol      | /var/log             | EmptyDir            |
+------------+-----------+--------------+----------------------+---------------------+
```
Installation
```bash
go install github.com/yeqianmen/kubectl-mounts@latest
```
Or build manually:

```bash
git clone https://github.com/yeqianmen/kubectl-mounts.git
cd kubectl-mounts
go build -o kubectl-mounts
sudo mv kubectl-mounts /usr/local/bin/
```

# kubectl-mounts

简体中文| [English](./README.md)

`kubectl-mounts` 是一个用 Go 编写的 `kubectl` 插件，用于展示当前 Kubernetes 命名空间中 Pod 卷（Volume）及其挂载路径的详细信息。

## 功能特性

- 显示 Pod 名称、容器名称、卷名称、挂载路径和卷类型
- 支持 ConfigMap、Secret、EmptyDir 和 PVC 类型的卷
- 使用 `tablewriter` 以人类可读的表格方式美观输出
- 自动检测当前上下文的命名空间

## 使用方式

```bash
在集群中显示 Pod 的 Volume 和 VolumeMount 信息

用法:
  kubectl-mounts [flags]
  kubectl-mounts [command]

可用命令:
  completion  生成命令补全脚本
  help        获取命令帮助信息

参数:
  -h, --help                查看帮助信息
  -k, --kubeconfig string   指定 kubeconfig 文件路径（默认为 $HOME/.kube/config）
  -n, --namespace string    指定命名空间（默认为当前命名空间）
  -o, --output string       输出格式：table|yaml|json（默认为 table）
  -p, --pod string          按指定 Pod 名过滤
  -v, --version             显示版本信息

使用 "kubectl-mounts [command] --help" 获取某个命令的更多信息。
```
示例输出：

```pgsql

+-----------------------------------------+-----------------------------------+-----------------------------------+-----------------------------------------------+--------------------------------------------------+
|                   Pod                   |             Container             |            Volume Name            |                   MountPath                   |                   Volume Type                    |
+-----------------------------------------+-----------------------------------+-----------------------------------+-----------------------------------------------+--------------------------------------------------+
| default-flyflow-console-8d97c74f-mxmgp  | default-flyflow-console-container | default-flyflow-console-nginxconf | /conf                                         | ConfigMap(cm-default-flyflow-console-nagix.conf) |
+-----------------------------------------+-----------------------------------+-----------------------------------+-----------------------------------------------+--------------------------------------------------+
| hostpath-demo                           | busybox                           | host-volume                       | /data/host                                    | HostPath(/data)                                  |
+-----------------------------------------+-----------------------------------+-----------------------------------+-----------------------------------------------+--------------------------------------------------+
| mysql-operator-0                        | orchestrator                      | data                              | /var/lib/orchestrator                         | PVC(data-mysql-operator-0)                       |
+-----------------------------------------+-----------------------------------+-----------------------------------+-----------------------------------------------+--------------------------------------------------+
| nfs-client-provisioner-844c6c9489-pl8bn | nfs-client-provisioner            | nfs-client-root                   | /persistentvolumes                            | NFS(192.168.150.183:/data/nfs)                   |
+-----------------------------------------+-----------------------------------+-----------------------------------+-----------------------------------------------+--------------------------------------------------+
| nginx-0                                 | nginx                             | kube-api-access-vflg6             | /var/run/secrets/kubernetes.io/serviceaccount | ServiceAccountToken                              |
+-----------------------------------------+-----------------------------------+-----------------------------------+-----------------------------------------------+--------------------------------------------------+

```
## 安装方式
### 使用 Homebrew 安装
如果你使用 Homebrew，可以通过以下方式安装：

```bash

brew tap yeqianmen/kubectl-mounts
brew install kubectl-mounts
```
### 使用 Krew（kubectl 插件管理器）
如果你使用 kubectl 并已安装 Krew，可通过以下方式安装本插件：

```bash
kubectl krew  install mounts
```
### 手动构建安装
```bash
git clone https://github.com/yeqianmen/kubectl-mounts.git
cd kubectl-mounts
go build -o kubectl-mounts
sudo mv kubectl-mounts /usr/local/bin/
```
## 命令补全（Completion）
支持以下命令补全功能：

- -n：自动补全集群中的命名空间

- -p：自动补全当前命名空间下的 Pod 名称

- -o：支持的输出格式有 table、yaml 和 json

### Bash

为当前会话加载补全脚本：

```bash
source <(kubectl-mounts completion bash)
```
### Zsh

为当前会话加载补全脚本：

```bash
source <(kubectl-mounts completion zsh)
```

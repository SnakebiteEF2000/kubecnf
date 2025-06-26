# kubecnf

Add and remove kubeconfigs from the main config file

## Usage

**Info:** With unspecified config file (-c / --config) default value is used (~/.kube/config)

### Add a new cluster config

From a file:
```
kubecnf [-c /path/to/main/config] add /path/to/new/cluster/config
```

From piped input:
```
cat /path/to/new/cluster/config | kubecnf [-c /path/to/main/config] add
```

```
kubectl config view --raw | kubecnf [-c /path/to/main/config] add
```

```
echo "apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://cluster.example.com
  name: example-cluster
contexts:
- context:
    cluster: example-cluster
    user: example-user
  name: example-context
users:
- name: example-user
  user:
    token: your-token-here" | kubecnf add
```

**Note:** You cannot use both a file argument and piped input simultaneously.

### Remove a cluster config
```
kubecnf [-c /path/to/main/config] remove cluster-name
```

### List all cluster configs
```
kubecnf [-c /path/to/main/config] list
```

### Rollback to the previous config
```
kubecnf [-c /path/to/main/config] rollback
```

## Bash Completion

To enable bash completion, source the provided completion script:

```
source <(kubecnf completion)
```

To make it permanent, add the above line to your `~/.bashrc` file.

## Installation

1. Clone the repository
2. Build the binary:
   ```
   go build -o kubecnf
   ```
3. Move the binary to a directory in your PATH, e.g.:
   ```
   sudo mv kubecnf /usr/local/bin/
   ```

Alternatively, you can download the pre-built binary directly from the GitHub [release page](https://github.com/SnakebiteEF2000/kubecnf/releases) rename it and move it to your PATH as shown above.

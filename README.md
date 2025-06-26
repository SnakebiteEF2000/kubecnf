# kubecnf

[![CI](https://github.com/SnakebiteEF2000/kubecnf/workflows/CI/badge.svg)](https://github.com/SnakebiteEF2000/kubecnf/actions/workflows/ci.yml)
[![Release](https://github.com/SnakebiteEF2000/kubecnf/workflows/Release/badge.svg)](https://github.com/SnakebiteEF2000/kubecnf/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/SnakebiteEF2000/kubecnf)](https://goreportcard.com/report/github.com/SnakebiteEF2000/kubecnf)
[![Docker](https://img.shields.io/badge/docker-available-blue)](https://github.com/SnakebiteEF2000/kubecnf/pkgs/container/kubecnf)

Add and remove kubeconfigs from the main config file with support for piped input and comprehensive shell completions.

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

## Shell Completion

kubecnf supports completion for both bash and zsh shells.

### Bash Completion

To enable bash completion, source the completion script:

```bash
source <(kubecnf completion bash)
```

To make it permanent, add the above line to your `~/.bashrc` file.

Alternatively, install system-wide (requires sudo):
```bash
kubecnf completion bash | sudo tee /etc/bash_completion.d/kubecnf
```

### Zsh Completion

To enable zsh completion, add the completion script to your fpath:

```zsh
# Add to ~/.zshrc
kubecnf completion zsh > "${fpath[1]}/_kubecnf"
```

Or source it directly:
```zsh
source <(kubecnf completion zsh)
```

For Oh My Zsh users, you can place the completion file in the completions directory:
```zsh
mkdir -p ~/.oh-my-zsh/completions
kubecnf completion zsh > ~/.oh-my-zsh/completions/_kubecnf
```

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

### Using Docker

```bash
# Pull and run from GitHub Container Registry
docker pull ghcr.io/snakebiteef2000/kubecnf:latest
docker run --rm -v ~/.kube:/root/.kube ghcr.io/snakebiteef2000/kubecnf:latest --help
```

### Using Makefile (for development)

```bash
# Install development tools and build
make dev-setup
make all

# Install locally
make install
```

## 🏗️ Development

For development instructions, testing, and contributing guidelines, see [DEVELOPMENT.md](DEVELOPMENT.md).

## 🚀 Features

- ✅ Add kubeconfigs from files or piped input
- ✅ Remove cluster configurations by name
- ✅ List all configured clusters
- ✅ Rollback to previous configurations
- ✅ Shell completions for bash and zsh
- ✅ Comprehensive input validation
- ✅ Automatic backups
- ✅ Duplicate detection with warnings
- ✅ Docker container support
- ✅ Multi-platform releases (Linux, macOS, Windows, FreeBSD)

## 📊 Project Status

This project follows professional development practices:

- 🔄 **Continuous Integration**: Automated testing, linting, and security scanning
- 📦 **Automated Releases**: Multi-platform binaries and Docker images
- 🔒 **Security**: Regular security scans and dependency updates
- 📈 **Code Quality**: Comprehensive linting and test coverage
- 🤖 **Automation**: Dependabot for dependency management

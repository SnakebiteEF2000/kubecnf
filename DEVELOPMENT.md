# Development Guide

This document explains how to develop, test, and release kubecnf using the professional CI/CD pipeline.

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- Docker (optional, for container builds)
- Make

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/SnakebiteEF2000/kubecnf.git
cd kubecnf

# Set up development tools
make dev-setup

# Run all checks (lint, test, build)
make all
```

## 🛠️ Development Workflow

### Local Development

```bash
# Build the binary
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Format code
make fmt

# Build for all platforms
make build-all
```

### Version Management

The project uses semantic versioning with automatic version detection:

- Version is automatically derived from Git tags
- Build time is embedded in binaries
- Use `./kubecnf version` to see version info

## 🔄 CI/CD Pipeline

### Continuous Integration (CI)

The CI pipeline runs on every push and pull request to `main` and `develop` branches:

1. **Test Stage**: Runs unit tests with race detection and coverage reporting
2. **Lint Stage**: Runs golangci-lint for code quality
3. **Security Stage**: Runs Gosec security scanner
4. **Build Stage**: Builds for multiple platforms (Linux, macOS, Windows, FreeBSD)
5. **Integration Test Stage**: Runs integration tests on built binaries

### Release Pipeline

The release pipeline triggers on version tags (e.g., `v1.2.3`):

1. **CI Checks**: Runs full CI pipeline first
2. **Multi-Platform Build**: Builds binaries for all supported platforms
3. **Docker Image**: Builds and pushes container image to GitHub Container Registry
4. **GitHub Release**: Creates release with binaries and checksums
5. **Homebrew**: Updates Homebrew formula (if configured)

## 📦 Creating a Release

### Step 1: Prepare Release

```bash
# Ensure you're on main branch
git checkout main
git pull origin main

# Update version in any relevant files if needed
# The version will be automatically set from the git tag
```

### Step 2: Create and Push Tag

```bash
# Create a new tag (use semantic versioning)
git tag v1.2.3

# Push the tag to trigger release
git push origin v1.2.3
```

### Step 3: Monitor Release

1. Go to GitHub Actions tab to monitor the release pipeline
2. Check the GitHub Releases page for the new release
3. Verify all platforms are built and uploaded

## 🐳 Docker

### Build Docker Image

```bash
# Build locally
make docker-build

# Run container
make docker-run
```

### Published Images

Docker images are automatically published to GitHub Container Registry:

- `ghcr.io/snakebiteef2000/kubecnf:latest` - Latest release
- `ghcr.io/snakebiteef2000/kubecnf:v1.2.3` - Specific version

```bash
# Pull and run
docker pull ghcr.io/snakebiteef2000/kubecnf:latest
docker run --rm ghcr.io/snakebiteef2000/kubecnf:latest --help
```

## 🔧 Code Quality

### Linting

The project uses golangci-lint with comprehensive checks:

```bash
# Run locally
make lint

# Auto-fix some issues
golangci-lint run --fix
```

### Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# View coverage report
open coverage.html
```

### Security

```bash
# Run security scanner
make security
```

## 📋 Makefile Targets

| Target | Description |
|--------|-------------|
| `all` | Run clean, lint, test, and build |
| `build` | Build the binary |
| `build-all` | Build for multiple platforms |
| `test` | Run tests |
| `test-coverage` | Run tests with coverage report |
| `lint` | Run linter |
| `fmt` | Format code |
| `tidy` | Tidy dependencies |
| `verify` | Verify dependencies |
| `security` | Run security scan |
| `clean` | Clean build artifacts |
| `install` | Install the binary |
| `docker-build` | Build Docker image |
| `docker-run` | Run Docker container |
| `completions` | Generate completion scripts |
| `dev-setup` | Set up development environment |
| `help` | Show help |

## 🤖 Automated Dependency Updates

Dependabot is configured to automatically:

- Update Go modules weekly
- Update GitHub Actions weekly  
- Update Docker base images weekly
- Create PRs with proper labels and assignees

## 📊 Monitoring and Metrics

### Code Coverage

- Target: 80%+ coverage
- Coverage reports uploaded to CI artifacts
- HTML reports generated locally with `make test-coverage`

### Security Scanning

- Gosec runs on every CI build
- SARIF results uploaded to GitHub Security tab
- Review security alerts in GitHub Security tab

## 🚨 Troubleshooting

### CI/CD Issues

1. **Build Failures**: Check GitHub Actions logs
2. **Test Failures**: Run `make test` locally
3. **Lint Failures**: Run `make lint` and fix issues
4. **Security Issues**: Check GitHub Security tab

### Local Development Issues

1. **Missing Tools**: Run `make dev-setup`
2. **Build Issues**: Ensure Go 1.23+ is installed
3. **Test Issues**: Check if dependencies are up to date with `make tidy`

## 🎯 Best Practices

1. **Commit Messages**: Use conventional commits (feat:, fix:, chore:, etc.)
2. **Pull Requests**: Ensure CI passes before merging
3. **Releases**: Always test locally before creating tags
4. **Security**: Review Dependabot PRs promptly
5. **Documentation**: Update docs when adding features

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/) 
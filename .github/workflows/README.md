# GitHub Actions Workflows

This directory contains GitHub Actions workflows for continuous integration and deployment.

## Workflows

### 🔄 CI Workflow (`ci.yml`)
**Main workflow for continuous integration**

- **Triggers**: Push to `main`/`develop` branches, Pull Requests
- **Jobs**:
  - **Test**: Runs comprehensive test suite across Go versions 1.21, 1.22, 1.23
  - **Lint**: Code quality checks using golangci-lint
  - **Build**: Compilation verification across all supported Go versions
- **Features**:
  - Go module caching for faster builds
  - Coverage reporting to Codecov
  - Dependency verification
  - Multi-version Go support

### 🧪 Test Workflow (`test.yml`)
**Comprehensive testing and security workflow**

- **Triggers**: Push to `main`/`develop` branches, Pull Requests
- **Jobs**:
  - **Test**: Full test suite with coverage reporting
  - **Lint**: Static code analysis
  - **Build**: Multi-version build verification
  - **Security**: Security scanning with Gosec
- **Features**:
  - SARIF security report upload
  - Detailed coverage analysis
  - Security vulnerability detection

### 🔍 Lint Workflow (`lint.yml`)
**Code quality and style checking**

- **Triggers**: Push and Pull Request events
- **Features**:
  - golangci-lint integration
  - Multiple linter configurations
  - Code style enforcement

## Badges

The following badges are available for the README:

```markdown
[![CI](https://github.com/itsLeonB/ezutil/workflows/CI/badge.svg)](https://github.com/itsLeonB/ezutil/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsLeonB/ezutil)](https://goreportcard.com/report/github.com/itsLeonB/ezutil)
[![codecov](https://codecov.io/gh/itsLeonB/ezutil/branch/main/graph/badge.svg)](https://codecov.io/gh/itsLeonB/ezutil)
```

## Configuration

### Required Secrets

For full functionality, the following GitHub secrets should be configured:

- `CODECOV_TOKEN`: Token for uploading coverage reports to Codecov (optional)

### Supported Go Versions

All workflows test against:
- Go 1.21
- Go 1.22  
- Go 1.23

### Caching Strategy

Workflows use GitHub Actions caching to speed up builds:
- Go build cache (`~/.cache/go-build`)
- Go module cache (`~/go/pkg/mod`)
- Cache keys include Go version and go.sum hash for optimal cache hits

## Local Testing

To run the same checks locally:

```bash
# Run tests
make test-coverage

# Run linter
make lint

# Build verification
go build -v ./...
cd test && go build -v ./...
```

## Workflow Status

Check the status of all workflows in the [Actions tab](https://github.com/itsLeonB/ezutil/actions) of the repository.

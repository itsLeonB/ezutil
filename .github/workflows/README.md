# GitHub Actions Workflows

This directory contains GitHub Actions workflows for continuous integration, testing, and security scanning.

## Workflows

### 🔄 CI Workflow (`ci.yml`)
**Main workflow for continuous integration**

- **Triggers**: Push to `main`/`develop` branches, Pull Requests
- **Jobs**:
  - **Test**: Runs comprehensive test suite with coverage across Go versions 1.23-1.24
  - **Lint**: Code quality checks using golangci-lint
  - **Build**: Compilation verification across all supported Go versions
- **Features**:
  - Go module caching for faster builds
  - Coverage reporting to Codecov
  - Dependency verification
  - Multi-version Go support (1.23, 1.24)

### 🧪 Test Workflow (`test.yml`)
**Comprehensive testing and security workflow**

- **Triggers**: Push to `main`/`develop` branches, Pull Requests
- **Jobs**:
  - **Test**: Full test suite with coverage reporting across Go versions 1.23-1.24
  - **Lint**: Static code analysis with golangci-lint
  - **Build**: Multi-version build verification
  - **Security**: Security scanning with Gosec
- **Features**:
  - SARIF security report upload to GitHub Security tab
  - Detailed coverage analysis with Codecov integration
  - Security vulnerability detection
  - Optimized single test run (no redundant executions)

### 🔍 Lint Workflow (`lint.yml`)
**Code quality and style checking**

- **Triggers**: Push and Pull Request events
- **Features**:
  - golangci-lint integration
  - Multiple linter configurations
  - Code style enforcement
  - Go best practices validation

## Security Features

### Gosec Integration
- **Static Security Analysis**: Automated security vulnerability scanning
- **SARIF Reports**: Security findings integrated with GitHub's security tab
- **Continuous Monitoring**: Security checks on every push and PR

### Permissions
All workflows use minimal required permissions:
```yaml
permissions:
  contents: read        # Read repository contents
  security-events: write # Upload security scan results
  actions: read         # Read workflow metadata
```

## Badges

The following badges are available for the README:

```markdown
[![CI](https://github.com/itsLeonB/ezutil/workflows/CI/badge.svg)](https://github.com/itsLeonB/ezutil/actions)
[![Tests](https://github.com/itsLeonB/ezutil/workflows/Tests/badge.svg)](https://github.com/itsLeonB/ezutil/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsLeonB/ezutil)](https://goreportcard.com/report/github.com/itsLeonB/ezutil)
[![codecov](https://codecov.io/gh/itsLeonB/ezutil/branch/main/graph/badge.svg)](https://codecov.io/gh/itsLeonB/ezutil)
[![Go Reference](https://pkg.go.dev/badge/github.com/itsLeonB/ezutil.svg)](https://pkg.go.dev/github.com/itsLeonB/ezutil)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
```

## Configuration

### Required Secrets

For full functionality, the following GitHub secrets should be configured:

- `CODECOV_TOKEN`: Token for uploading coverage reports to Codecov (optional but recommended)

### Supported Go Versions

All workflows test against:
- Go 1.23 (minimum supported)
- Go 1.24 (latest)

### Caching Strategy

Workflows use GitHub Actions caching to speed up builds:
- Go build cache (`~/.cache/go-build`)
- Go module cache (`~/go/pkg/mod`)
- Cache keys include Go version and go.sum hash for optimal cache hits
- Separate caches for main module and test module

## Performance Optimizations

### Efficient Test Execution
- **Single Test Run**: Tests run once with both verbose output and coverage
- **Conditional Coverage Upload**: Only uploads coverage from Go 1.24 to avoid duplicates
- **Parallel Jobs**: Different jobs run in parallel for faster overall execution

### Smart Caching
- **Version-Specific Caches**: Separate caches for each Go version
- **Dependency Caching**: Both main and test module dependencies cached
- **Restore Keys**: Fallback cache keys for partial cache hits

## Local Testing

To run the same checks locally:

```bash
# Run tests with coverage (same as CI)
make test-coverage

# Run tests with verbose output
make test-verbose

# Run linter (same as CI)
make lint

# Build verification
go build -v ./...
cd test && go build -v ./...

# Security scan (requires gosec installation)
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...
```

## Test Coverage

The project maintains comprehensive test coverage:
- **200+ test cases** covering all exported functions
- **Multiple test scenarios**: Happy path, error conditions, edge cases
- **Database testing**: In-memory SQLite for isolation
- **HTTP testing**: Complete request/response cycle testing
- **Configuration testing**: Environment variable processing
- **Transaction testing**: Database transaction management

## Workflow Status

Check the status of all workflows in the [Actions tab](https://github.com/itsLeonB/ezutil/actions) of the repository.

## Troubleshooting

### Common Issues

1. **Test Failures**: Check the test logs in the Actions tab
2. **Lint Failures**: Run `make lint` locally to see specific issues
3. **Coverage Issues**: Ensure tests are properly written and cover edge cases
4. **Security Scan Issues**: Review Gosec findings in the Security tab

### Getting Help

- Check the [Issues](https://github.com/itsLeonB/ezutil/issues) for known problems
- Review the [Contributing Guide](../README.md#contributing) for development setup
- Look at recent [Pull Requests](https://github.com/itsLeonB/ezutil/pulls) for examples

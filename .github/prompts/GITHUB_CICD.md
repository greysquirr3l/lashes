# GitHub Actions CI/CD Best Practices

This document provides concise guidance for implementing effective CI/CD workflows with GitHub Actions.

<!-- REF: https://docs.github.com/en/actions -->
<!-- REF: https://www.linkedin.com/pulse/mastering-cicd-best-practices-github-actions-tatiana-sava -->

## 📋 Workflow Structure

### Basic Workflow Anatomy

```yaml
name: CI Pipeline

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up environment
        uses: actions/setup-node@v3
        with:
          node-version: '16'
      - name: Install dependencies
        run: npm ci
      - name: Run tests
        run: npm test
```

### Key Components

- **Workflow file**: YAML files in `.github/workflows/`
- **Triggers** (`on`): Events that start the workflow
- **Jobs**: Groups of steps that run on the same runner
- **Steps**: Individual tasks that run commands or actions
- **Actions**: Reusable units of code (e.g., `actions/checkout@v3`)
- **Runners**: VMs that execute the jobs

## 🚀 Best Practices

### 1. Workflow Design

- **Keep workflows focused**: One workflow per logical CI/CD phase
- **Use descriptive names**: Clear workflow, job, and step names
- **Modularize with composite actions**: Create reusable components
- **Trigger precision**: Limit workflow runs to relevant files/branches

```yaml
# Good practice: Focused triggers
on:
  push:
    branches: [ main ]
    paths:
      - 'src/**'
      - 'package.json'
  pull_request:
    branches: [ main ]
    paths-ignore:
      - '**/*.md'
```

### 2. Performance Optimization

- **Use dependency caching**: Speed up builds by caching dependencies
- **Conditional execution**: Skip unnecessary steps
- **Job parallelization**: Run independent jobs concurrently
- **Matrix builds**: Test across multiple configurations in parallel

```yaml
# Example: Dependency caching
steps:
  - uses: actions/checkout@v3
  - name: Cache dependencies
    uses: actions/cache@v3
    with:
      path: ~/.npm
      key: ${{ runner.os }}-node-${{ hashFiles('**/package-lock.json') }}
      restore-keys: ${{ runner.os }}-node-

# Example: Matrix builds
jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        node-version: [14.x, 16.x, 18.x]
```

### 3. Security

- **Protect secrets**: Use GitHub Secrets for sensitive data
- **Restrict permissions**: Apply principle of least privilege
- **Pin action versions**: Use specific versions (SHA is safest)
- **Scan for vulnerabilities**: Include security scanning in workflows
- **Review third-party actions**: Audit external actions for security risks

```yaml
# Good practice: Define specific permissions
permissions:
  contents: read
  pull-requests: write

# Good practice: Pin actions to specific SHAs
steps:
  - uses: actions/checkout@a81bbbf8298c0fa03ea29cdc473d45769f953675 # v3.2.0
```

### 4. Workflow Management

- **Document workflows**: Add comments explaining complex steps
- **Self-contained workflows**: Minimize external dependencies
- **Re-usable workflows**: Define common workflows in separate files
- **Timeout limits**: Add timeouts to prevent stuck jobs

```yaml
# Example: Job timeout
jobs:
  build:
    runs-on: ubuntu-latest
    timeout-minutes: 15
    steps:
      # ...
```

## 📝 Code Quality and Testing

- **Linting**: Check code quality in early stages
- **Automated tests**: Run unit, integration, and e2e tests
- **Coverage reports**: Track test coverage over time
- **Build artifacts**: Upload build artifacts for verification

```yaml
steps:
  - name: Run tests
    run: npm test

  - name: Upload coverage reports
    uses: actions/upload-artifact@v3
    with:
      name: coverage-report
      path: coverage/
```

## 🚢 Deployment

### Environments and Approvals

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v3
      - name: Deploy
        run: ./deploy.sh
```

### Continuous Deployment Practices

- **Environment-specific workflows**: Separate workflows for dev/staging/prod
- **Approval gates**: Require approval for sensitive environments
- **Deployment verification**: Add post-deployment verification steps
- **Rollback capability**: Plan for failed deployments

## 🔄 Advanced Patterns

### 1. Workflow Composition

```yaml
# Reusable workflow in .github/workflows/reusable-build.yml
name: Reusable build
on:
  workflow_call:
    inputs:
      config:
        required: true
        type: string

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build with config
        run: ./build.sh ${{ inputs.config }}

# Calling workflow
jobs:
  call-build:
    uses: ./.github/workflows/reusable-build.yml
    with:
      config: 'production'
```

### 2. Custom Actions

Create custom actions for project-specific tasks:

```bash
my-project/
├── .github/
│   └── actions/
│       └── custom-build/
│           ├── action.yml
│           └── index.js
```

### 3. CI/CD Metrics Tracking

Monitor workflow performance, success rates, and durations.

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitHub Actions Marketplace](https://github.com/marketplace?type=actions)
- [GitHub Actions Security Hardening](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions)
- [GitHub Actions Workflow Syntax Reference](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions)

## References

1. [GitHub Actions Documentation](https://docs.github.com/en/actions)
2. [Mastering CI/CD Best Practices](https://www.linkedin.com/pulse/mastering-cicd-best-practices-github-actions-tatiana-sava)
3. [GitHub Actions Marketplace](https://github.com/marketplace?type=actions)
4. [GitHub Actions Security Best Practices](https://blog.gitguardian.com/github-actions-security-best-practices/)

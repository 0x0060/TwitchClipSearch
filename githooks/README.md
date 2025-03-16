# Git Hooks

This directory contains Git hooks for maintaining code quality and automating development workflows.

## Available Hooks

### Pre-commit
- Code formatting check
- Lint checks
- Unit test execution
- Commit message validation

### Pre-push
- Integration test execution
- Security checks
- Build verification

### Post-merge
- Dependency updates
- Database migrations

## Installation

```bash
# Make hooks executable
chmod +x githooks/*

# Install hooks
git config core.hooksPath githooks
```

## Configuration

Hook behavior can be configured in `.githooks.yaml`:

```yaml
hooks:
  pre-commit:
    skip_tests: false
    lint_threshold: "warning"
  pre-push:
    run_integration_tests: true
    security_scan: true
```

## Skip Hooks

To skip hook execution (not recommended):
```bash
git commit --no-verify
git push --no-verify
```
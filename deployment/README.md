# Deployment

This directory contains deployment configurations and scripts for the TwitchClipSearch application.

## Structure

- `docker/`: Docker-related files and configurations
  - `Dockerfile`: Main application Dockerfile
  - `docker-compose.yml`: Development and testing environment setup
  - `docker-compose.prod.yml`: Production environment setup
- `kubernetes/`: Kubernetes manifests and configurations
  - `base/`: Base Kubernetes configurations
  - `overlays/`: Environment-specific overlays
- `terraform/`: Infrastructure as Code configurations
  - `modules/`: Reusable Terraform modules
  - `environments/`: Environment-specific configurations
- `scripts/`: Deployment automation scripts

## Deployment Environments

- Development
- Staging
- Production

## Deployment Process

1. Build and test application
2. Create Docker image
3. Push to container registry
4. Deploy using Kubernetes/Terraform
5. Run post-deployment checks

## Configuration

Environment-specific configurations are managed through:
- Environment variables
- ConfigMaps
- Secrets

## Monitoring

- Health checks
- Resource monitoring
- Logging
- Alerting
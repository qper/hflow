# Deployment notes

## Runtime assumptions
- The target cluster is assumed to be a self-hosted K3s environment with ingress support via Traefik.
- No existing monitoring stack was specified in the repository or prompt, so the deployment scaffolding exposes Prometheus-compatible metrics on `/metrics` and leaves the scraping integration opt-in via a ServiceMonitor or external scraper.
- Akeyless integration is modeled through an ExternalSecret resource and requires the External Secrets Operator plus a configured SecretStore in the cluster.

## Verification
- `go test ./...`
- `helm template ./deploy/helm/habitflow`

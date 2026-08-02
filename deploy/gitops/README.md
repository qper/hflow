# GitOps layout

This repository keeps deployment manifests in [deploy/helm/habitflow](../helm/habitflow) and treats the Helm values file as the GitOps source of truth.

Recommended flow:
1. CI builds and publishes multi-arch images to GHCR.
2. CI updates the image tag in the Helm values file.
3. Argo CD (or Flux) syncs the repository into the K3s cluster.

For a self-hosted K3s cluster, Argo CD is a pragmatic choice because it supports Helm-based apps and pull-based reconciliation without requiring a managed control plane.

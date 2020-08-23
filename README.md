# Kustomize helm plugins

This repo help to generate [Helm Chart](https://helm.sh) from [Kustomize](https://kustomize.io).

Read more about [Kustomize Plugins](https://github.com/kubernetes-sigs/kustomize/tree/master/docs/plugins).

Version: v0.1.0

## Installation
### build image
```
docker build -t kustomize-helm-plugins .
```
### usage

Add transformers to kustomization.yaml
```yaml
transformers:
- helm.yaml
```
helm.yaml
```yaml
apiVersion: helm/v1
kind: HelmTransformer
metadata:
  name: external-dns
ChartName: external-dns
ChartVersion: 0.1.0
appVersion: 0.2.0
Values:
  DOMAIN: example.com
  PROVIDER: aws
```
## Example
### [kustomize for external-dns](https://github.com/kubernetes-sigs/external-dns/tree/master/kustomize) to Helm Chart
```
cd example
mkdir -p external-dns/templates
docker run -v ${PWD}:/working kustomize-helm-plugins overlays/helm > external-dns/templates/external-dns.yaml
helm template --release-name kustomize --namespace external-dns-ns -f external-dns/values.yaml external-dns/
```

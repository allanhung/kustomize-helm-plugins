module github.com/allanhung/kustomize-helm-plugins/helm/v1

go 1.14

require (
  github.com/pkg/errors v0.9.1
	gopkg.in/yaml.v3 v3.0.0-20200121175148-a6ecf24a6d71
	sigs.k8s.io/kustomize/api v0.5.1
	sigs.k8s.io/yaml v1.2.0
)

replace (
	sigs.k8s.io/kustomize/api => /go/src/github.com/kubernetes-sigs/kustomize/api
)

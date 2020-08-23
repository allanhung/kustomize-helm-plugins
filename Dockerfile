FROM golang:1.14 AS builder

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOPATH=/go
ENV KUSTOMIZE_VERSION=3.8.1
ENV API_VERSION=0.5.1

RUN apt-get update && apt-get install -y curl gettext g++ git 

WORKDIR /go/src/github.com/kubernetes-sigs
RUN git clone -b api/v${API_VERSION} https://github.com/kubernetes-sigs/kustomize.git
RUN cd kustomize/kustomize && go build \
  -ldflags="-w -s -X sigs.k8s.io/kustomize/api/provenance.version=${KUSTOMIZE_VERSION} -X sigs.k8s.io/kustomize/api/provenance.gitCommit=$(git rev-parse HEAD) -X sigs.k8s.io/kustomize/api/provenance.buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
  -o /go/bin/kustomize

WORKDIR /go/src/github.com/kustomize-plugins
COPY helm/v1/helmtransformer/go.mod ./
RUN go mod download

COPY helm/v1/helmtransformer/HelmTransformer.go ./
RUN go build -buildmode plugin -o /go/bin/HelmTransformer.so HelmTransformer.go

FROM debian:stretch-slim

RUN apt-get update && apt-get install -y git

ENV XDG_CONFIG_HOME=/opt

COPY --from=builder /go/bin/HelmTransformer.so /opt/kustomize/plugin/helm/v1/helmtransformer/HelmTransformer.so
COPY --from=builder /go/bin/kustomize /usr/bin/kustomize

WORKDIR /working 
ENTRYPOINT ["/usr/bin/kustomize", "build", "--enable_alpha_plugins"]

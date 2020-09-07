package main

import (
	"io/ioutil"
	"log"
	"os"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
	"sigs.k8s.io/yaml"
)

type chart struct {
	ApiVersion  string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string `json:"version,omitempty" yaml:"version,omitempty"`
	AppVersion  string `json:"appVersion,omitempty" yaml:"appVersion,omitempty"`
}

type plugin struct {
	ChartName    string                 `json:"chartName,omitempty" yaml:"chartName,omitempty"`
	ChartVersion string                 `json:"chartVersion,omitempty" yaml:"chartVersion,omitempty"`
	AppVersion   string                 `json:"appVersion,omitempty" yaml:"appVersion,omitempty"`
	Values       map[string]interface{} `json:"values,omitempty" yaml:"values,omitempty"`

	Logger *log.Logger
	h      *resmap.PluginHelpers
}

var KustomizePlugin plugin

func (p *plugin) Config(
	h *resmap.PluginHelpers, c []byte) (err error) {
	p.h = h

	err = yaml.Unmarshal(c, p)
	if err != nil {
		return nil
	}
	p.Logger = log.New(os.Stdout, "[DEBUG] ", log.Lshortfile)
	return nil
}

func (p *plugin) Transform(m resmap.ResMap) (err error) {
	var resources []*resource.Resource
	lb := make(map[string]string)
	lb["helm.sh/chart"] = "{{ printf \"%s-%s\" .Chart.Name .Chart.Version | replace \"+\" \"_\" | trunc 63 | trimSuffix \"-\" }}"
	lb["app.kubernetes.io/name"] = "{{ .Chart.Name }}"
	lb["app.kubernetes.io/instance"] = "{{ .Release.Name }}"
	lb["app.kubernetes.io/version"] = "{{ .Chart.Version }}"
	lb["app.kubernetes.io/managed-by"] = "{{ .Release.Service }}"
	for _, res := range m.Resources() {
		if res.GetKind() != "Namespace" {
			resources = append(resources, res)
			newlb := make(map[string]string)
			for k, v := range lb {
				newlb[k] = v
			}
			for k, v := range res.GetLabels() {
				newlb[k] = v
			}
			res.SetLabels(newlb)
			res.SetName("{{ .Release.Name }}-" + res.GetName())
			res.SetNamespace("{{ .Release.Namespace }}")
		}
	}
	m.Clear()
	for _, r := range resources {
		m.Append(r)
	}

	if err := p.createChartYaml(); err != nil {
		return err
	}
	if err := p.createChartValue(); err != nil {
		return err
	}
	return nil
}

func (p *plugin) createChartYaml() (err error) {
	os.MkdirAll(p.ChartName+"/templates", os.ModePerm)
	chartYaml := &chart{
		ApiVersion:  "v3",
		Name:        p.ChartName,
		Description: p.ChartName + " Helm chart for Kubernetes",
		Version:     p.ChartVersion,
		AppVersion:  p.AppVersion,
	}
	chartBytes, err := yaml.Marshal(chartYaml)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(p.ChartName+"/Chart.yaml", chartBytes, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (p *plugin) createChartValue() (err error) {
	chartValues, err := yaml.Marshal(p.Values)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(p.ChartName+"/values.yaml", chartValues, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

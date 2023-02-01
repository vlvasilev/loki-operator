package config

import (
	"embed"
	"text/template"
)

const (
	// ValiConfigFileName is the name of the config file in the configmap
	ValiConfigFileName = "config.yaml"
	// ValiConfigMountDir is the path that is mounted from the configmap
	ValiConfigMountDir = "/etc/vali/config"
)

var (
	//go:embed vali-config.yaml
	valiConfigYAMLTmplFile embed.FS
	ValiConfigYAMLTmpl     = template.Must(template.ParseFS(valiConfigYAMLTmplFile, "vali-config.yaml"))

	//go:embed curator-config.yaml
	curatorConfigYAMLTmplFile embed.FS
	CuratorConfigYAMLTmpl     = template.Must(template.ParseFS(curatorConfigYAMLTmplFile, "curator-config.yaml"))
)

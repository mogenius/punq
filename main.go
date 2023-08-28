package main

import (
	"embed"

	"github.com/mogenius/punq/cmd"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/operator"
	"github.com/mogenius/punq/utils"
)

//go:embed ui/dist/*
var htmlDirFs embed.FS

//go:embed config/local.yaml
var localConfig string

//go:embed config/dev.yaml
var devConfig string

//go:embed config/prod.yaml
var prodConfig string

//go:embed yaml-templates
var yamlTemplatesFolder embed.FS

func main() {
	utils.DefaultConfigLocalFile = localConfig
	utils.DefaultConfigFileDev = devConfig
	utils.DefaultConfigFileProd = prodConfig
	utils.YamlTemplatesFolder = yamlTemplatesFolder

	operator.HtmlDirFs = htmlDirFs

	logger.Init()
	cmd.Execute()
}

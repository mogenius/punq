package main

import (
	"embed"

	"github.com/mogenius/punq/cmd"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/operator"
	"github.com/mogenius/punq/utils"
)

// SWAGGER <-- DO NOT REMOVE
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Use the user/login route to get the token.
// SWAGGER --> DO NOT REMOVE

//go:embed ui/dist/*
var htmlDirFs embed.FS

//go:embed config/local.yaml
var localConfig string

//go:embed config/operator.yaml
var operatorConfig string

//go:embed config/prod.yaml
var prodConfig string

//go:embed CHANGELOG.md
var changelog string

func main() {
	utils.DefaultConfigLocalFile = localConfig
	utils.DefaultConfigFileOperator = operatorConfig
	utils.DefaultConfigFileProd = prodConfig

	operator.HtmlDirFs = htmlDirFs

	utils.ChangeLog = changelog

	logger.Init()
	cmd.Execute()
}

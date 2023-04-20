package main

import (
	"os"

	api "github.com/tsanton/tfe-client/tfe"
	apim "github.com/tsanton/tfe-client/tfe/models"
	a "github.com/tsanton/tfe-provider-release-action/action"

	log "github.com/sirupsen/logrus"
)

var (
	hostname string
	token    string
)

var runConfig a.RunConfig

func init() {
	//set variables with defaults
	//providerName allow empty -> from metadata
	//versionNumber from metadata

	//get and serialize manifest input environment
	//get and serialize artifacts input environment
	//protocolVersion: '["5.0", "6.0"]'
	runConfig = a.RunConfig{}
}

func main() {
	//configure client
	logger := log.New()
	cli, err := api.NewClient(logger, &apim.ClientConfig{
		Address: hostname,
		Token:   token,
	})
	if err != nil {
		panic("unable to configure ")
	}
	logger.Info("Starting the release action")
	err = a.Run(cli, logger, &runConfig)
	if err != nil {
		logger.Info("Unsuccessful upload")
		os.Exit(1)
	}
	logger.Info("Done!")
}

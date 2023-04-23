package main

import (
	"os"

	api "github.com/tsanton/tfe-client/tfe"
	apim "github.com/tsanton/tfe-client/tfe/models"
	a "github.com/tsanton/tfe-provider-release-action/action"
	u "github.com/tsanton/tfe-provider-release-action/utilities"

	log "github.com/sirupsen/logrus"
)

var (
	logger   u.ILogger
	hostname string
	token    string
)

var runConfig *a.RunConfig

func init() {
	logger = log.New()
	hostname = u.GetEnv("TFE_HOSTNAME", "")
	token = u.GetEnv("TFE_TOKEN", "")
	workdir := u.GetEnv("APP_WORKDIR", "")
	organization := u.GetEnv("TFE_ORGANIZATION", "")
	namespace := u.GetEnv("TFE_NAMESPACE", "")
	providerName := u.GetEnv("TFE_PROVIDER_NAME", "")
	gpgKeyId := u.GetEnv("TFE_GPG_KEY_ID", "")
	runConfig = a.NewRunConfig(workdir, organization, namespace, providerName, gpgKeyId)
	_ = runConfig.ParseGoreleaseArtifacts(logger, os.Getenv("GORELEASE_ARTIFACTS"))
	_ = runConfig.ParseGoreleaserMetadata(logger, os.Getenv("GORELEASE_METADATA"))
}

func main() {
	//configure client
	cli, err := api.NewClient(logger, &apim.ClientConfig{
		Address: hostname,
		Token:   token,
	})
	if err != nil {
		panic("unable to configure ")
	}
	//TODO: calidate input
	// logger.Info("Validating the input")

	// logger.Info("Valid input")

	logger.Info("Starting the release action")
	err = a.Run(cli, logger, runConfig)
	if err != nil {
		logger.Info("Unsuccessful upload")
		os.Exit(1)
	}
	logger.Info("Done!")
}

package action

import (
	"encoding/json"
	"fmt"

	m "github.com/tsanton/tfe-provider-release-action/action/models"
	u "github.com/tsanton/tfe-provider-release-action/utilities"
)

type RunConfig struct {
	/* Environment inputs */
	Organization string
	Namespace    string
	ProviderName string
	GpgKeyId     string
	/*Parsed from goreleaser*/
	// Version                  string
	GoreleaserMetadata       m.GoreleaserMetadata
	ProviderVersionPlatforms []m.GoreleaserArtifactArchive
	Checksum                 m.GoreleaserArtifactFile
	Signature                m.GoreleaserArtifactFile
	RegistryManifest         m.GoreleaserArtifactFile
	ProtocolVersions         []string
}

func NewRunConfig(organization, namespace, providerName, gpgKeyId string) *RunConfig {
	return &RunConfig{
		Organization: organization,
		Namespace:    namespace,
		ProviderName: providerName,
		GpgKeyId:     gpgKeyId,
	}
}

func (c *RunConfig) ParseGoreleaserMetadata(logger u.ILogger, manifest string) error {
	var ret m.GoreleaserMetadata
	logger.Infof("Trying to unmarshal the following string: %s", manifest)
	err := json.Unmarshal([]byte(manifest), &ret)
	if err != nil {
		return err
	}
	c.GoreleaserMetadata = ret
	return nil
}

func (c *RunConfig) ParseGoreleaseArtifacts(logger u.ILogger, artifacts string) error {
	releases := []m.GoreleaserArtifactArchive{}
	var checksum, signature, registryManifest m.GoreleaserArtifactFile
	logger.Infof("Trying to unmarshal the following string: %s", artifacts)
	var maps []map[string]interface{}
	err := json.Unmarshal([]byte(artifacts), &maps)
	if err != nil {
		return err
	}

	for _, mp := range maps {
		internalType, ok := mp["internal_type"].(float64)
		if !ok {
			return fmt.Errorf("invalid internal_type: %v", mp["internal_type"])
		}

		switch int(internalType) {
		case 1: //type: Archive
			artifact, err := u.MapToStruct[m.GoreleaserArtifactArchive](mp)
			if err != nil {
				return err
			}
			releases = append(releases, artifact)
		case 3: //type: File -> should be terraform-registry-manifest.json
			registryManifest, err = u.MapToStruct[m.GoreleaserArtifactFile](mp)
			if err != nil {
				return err
			}
		case 12: //type: Checksum
			checksum, err = u.MapToStruct[m.GoreleaserArtifactFile](mp)
			if err != nil {
				return err
			}
		case 13: //type: Signature
			signature, err = u.MapToStruct[m.GoreleaserArtifactFile](mp)
			if err != nil {
				return err
			}
		}
	}

	c.ProviderVersionPlatforms = releases
	c.Checksum = checksum
	c.Signature = signature
	c.RegistryManifest = registryManifest

	return nil
}

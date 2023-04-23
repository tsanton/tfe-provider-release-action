package action

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/samber/lo"

	api "github.com/tsanton/tfe-client/tfe"
	areq "github.com/tsanton/tfe-client/tfe/models/request"
	aresp "github.com/tsanton/tfe-client/tfe/models/response"
	m "github.com/tsanton/tfe-provider-release-action/action/models"
	u "github.com/tsanton/tfe-provider-release-action/utilities"
)

func Run(cli *api.TerraformEnterpriseClient, logger u.ILogger, config *RunConfig) error {
	ctx := context.Background()
	// 0. Check that the provider exists. If it does not, abort

	// 1. Check if the provider version exists. If it does not, create a provider version
	var providerVersion *aresp.ProviderVersion
	logger.Infof("Checking if provider version %s exists", config.GoreleaserMetadata.Version)
	providerVersion, err := checkProviderVersionExists(ctx, cli, config) //Get -> return nil if 404
	if err != nil {
		logger.Error("Error checking if provider version exists")
		return err //TODO: custom error message
	} else if providerVersion == nil {
		logger.Infof("Creating provider version %s", config.GoreleaserMetadata.Version)
		providerVersion, err = createProviderVersion(ctx, cli, config)
		if err != nil {
			logger.Errorf("Error creating provider version:", err)
			return err
		}
	}

	// 2. Upload the SHA256SUMS and SHA256SUMS.sig files if not uploaded:
	if !providerVersion.Data.Attributes.ShasumsUploaded {
		logger.Info("Uploading SHA256SUMS file")
		fullPath := path.Join(config.Workdir, config.Checksum.Path)
		fileContent, err := ReadFileToByteArray(fullPath)
		if err != nil {
			logger.Errorf("Error reading SHA256SUMS file from path: %s", fullPath)
			return err
		}
		request, err := CreateRequestFromByteArray(http.MethodPost, providerVersion.Data.Links.ShasumsUploadUrl, fileContent, config.Checksum.Name)
		if err != nil {
			return err
		}
		request.Header.Set("Content-Type", "application/vnd.api+json")
		request.Header.Set("Accept", "application/json")
		_, err = api.Do[interface{}](ctx, cli, 200, request)
		if err != nil {
			logger.Errorf("Error Uploading SHA256SUMS file")
			return err
		}
	}
	if !providerVersion.Data.Attributes.ShasumsSigUploaded {
		logger.Info("Uploading SHA256SUMS signature file")
		fullPath := path.Join(config.Workdir, config.Signature.Path)
		fileContent, err := ReadFileToByteArray(fullPath)
		if err != nil {
			logger.Errorf("Error reading SHA256SUMS signature file from path: %s", fullPath)
			return err
		}
		request, err := CreateRequestFromByteArray(http.MethodPost, providerVersion.Data.Links.ShasumsSigUploadUrl, fileContent, config.Signature.Name)
		if err != nil {
			return err
		}
		request.Header.Set("Content-Type", "application/vnd.api+json")
		request.Header.Set("Accept", "application/json")
		_, err = api.Do[interface{}](ctx, cli, 200, request)
		if err != nil {
			logger.Errorf("Error Uploading SHA256SUMS signature file")
			return err
		}
	}

	// 3. Get provider platform for the version
	logger.Info("Getting existing provider platform versions")
	providerPlatforms, err := getProviderVersionPlatforms(ctx, cli, config)
	if err != nil {
		return err
	}
	for _, platform := range config.ProviderVersionPlatforms {
		logger.Infof("Processing provider platform version: %s %s", platform.Os, platform.Goarch)
		var tfplatform *aresp.ProviderVersionPlatformData
		var exists bool
		// 4. Check if the provider platform exists. If it does not, create a provider platform (exist by match on os + arch, and provider-binary-uploaded: false)
		plat, exists := lo.Find(*providerPlatforms, func(i aresp.ProviderVersionPlatformData) bool {
			return i.Attributes.Os == platform.Os && i.Attributes.Arch == platform.Goarch
		})
		if !exists {
			logger.Infof("Creating provider platform version: %s %s", platform.Os, platform.Goarch)
			tfplatform, err = createProviderPlatformVersion(ctx, cli, config, platform)
			if err != nil {
				logger.Errorf("Error creating provider version platform %s %s", platform.Os, platform.Goarch)
				return err
			}
		} else {
			tfplatform = &plat
		}
		// 5. Upload the provider binary zip file if not uploaded
		// - when provider-binary-uploaded: false, instead of including a link to provider-binary-download, the response will include an upload link provider-binary-upload
		if !tfplatform.Attributes.ProviderBinaryUploaded {
			logger.Infof("Uploading provider platform binary: %s %s", platform.Os, platform.Goarch)
			fullPath := path.Join(config.Workdir, platform.Path)
			fileContent, err := ReadFileToByteArray(fullPath)
			if err != nil {
				logger.Errorf("Error reading provider platform binary at path: %s", fullPath)
				return err
			}
			request, err := CreateRequestFromByteArray(http.MethodPost, tfplatform.Links.ProviderBinaryUpload, fileContent, platform.Name)
			if err != nil {
				return err
			}
			_, err = api.Do[interface{}](ctx, cli, 200, request)
			if err != nil {
				logger.Errorf("Error Uploading provider platform binary")
				return err
			}
			logger.Info("Successfully uploaded provider platform binary")
		}
	}
	logger.Info("Provider version and platform created/uploaded successfully!")
	return nil
}

func checkProviderVersionExists(ctx context.Context, cli *api.TerraformEnterpriseClient, config *RunConfig) (*aresp.ProviderVersion, error) {
	v, err := cli.ProviderVersionService.Read(ctx, config.Organization, config.Namespace, config.ProviderName, config.GoreleaserMetadata.Version)
	if err != nil && strings.Contains(err.Error(), "returned non 200 response: 404") {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &v, nil
}

func createProviderVersion(ctx context.Context, cli *api.TerraformEnterpriseClient, config *RunConfig) (*aresp.ProviderVersion, error) {
	request := areq.ProviderVersion{
		Data: areq.ProviderVersionData{
			Type: "registry-provider-versions",
			Attributes: areq.ProviderVersionDataAttributes{
				Version:   config.GoreleaserMetadata.Version,
				KeyId:     config.GpgKeyId,
				Protocols: config.ProtocolVersions,
			},
		},
	}
	providerVersion, err := cli.ProviderVersionService.Create(ctx, config.Organization, config.Namespace, config.ProviderName, &request)
	if err != nil {
		fmt.Println("Error creating provider version:", err)
		return nil, err
	}
	return &providerVersion, nil
}

func getProviderVersionPlatforms(ctx context.Context, cli *api.TerraformEnterpriseClient, config *RunConfig) (*[]aresp.ProviderVersionPlatformData, error) {
	platforms, err := cli.ProviderVersionPlatformService.List(ctx, config.Organization, config.Namespace, config.ProviderName, config.GoreleaserMetadata.Version)
	if err != nil {
		return nil, err
	}
	return &platforms.Data, nil
}

func createProviderPlatformVersion(ctx context.Context, cli *api.TerraformEnterpriseClient, config *RunConfig, platform m.GoreleaserArtifactArchive) (*aresp.ProviderVersionPlatformData, error) {
	request := areq.ProviderVersionPlatform{
		Data: areq.ProviderVersionPlatformData{
			Type: "registry-provider-platform-versions",
			Attributes: areq.ProviderVersionPlatformDataAttributes{
				Os:       platform.Os,
				Arch:     platform.Goarch,
				Filename: platform.Name,
				Shasum:   platform.ShaSum,
			},
		},
	}
	response, err := cli.ProviderVersionPlatformService.Create(ctx, config.Organization, config.Namespace, config.ProviderName, config.GoreleaserMetadata.Version, &request)
	if err != nil {
		return nil, err
	}
	return &response.Data, nil
}

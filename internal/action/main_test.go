package action_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	logt "github.com/sirupsen/logrus/hooks/test"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"

	api "github.com/tsanton/tfe-client/tfe"
	apim "github.com/tsanton/tfe-client/tfe/models"
	me "github.com/tsanton/tfe-client/tfe/models/enum"
	mreq "github.com/tsanton/tfe-client/tfe/models/request"
	mresp "github.com/tsanton/tfe-client/tfe/models/response"
	u "github.com/tsanton/tfe-client/tfe/utilities"
)

var (
	logger u.ILogger
)

func TestMain(m *testing.M) {
	log.Println("Main test setup")

	logger, _ = logt.NewNullLogger()

	log.Println("Running tests")
	exitVal := m.Run()
	log.Println("Main test teardown")

	os.Exit(exitVal)
}

/*###########################
### Test utility function ###
###########################*/

func runnerValidator(t *testing.T) (string, string) {
	orgName := u.GetEnv("TFE_ORG_NAME", "")
	token := u.GetEnv("TFE_TOKEN", "")
	run := u.GetEnv("TFE_RUN_LIVE_TESTS", false)
	if !run && orgName != "" && token != "" {
		t.Skip("Skipping test 'Test_live_provider_service_version_lifecycle'")
	}
	return orgName, token
}

func liveClientSetup(t *testing.T, host, token string) *api.TerraformEnterpriseClient {
	cli, err := api.NewClient(logger, &apim.ClientConfig{
		Address: host,
		Token:   token,
	})
	if err != nil {
		t.Errorf("Error creating client: %s", err)
		t.FailNow()
	}
	return cli
}

func boostrapGpgKey(t *testing.T, cli *api.TerraformEnterpriseClient, orgName, comment string) *mresp.GpgKey {
	entity, err := openpgp.NewEntity(orgName, comment, "donotreply@gruntwork.com", &packet.Config{RSABits: 4096})
	if err != nil {
		panic("unable to generate GPG entity in boostrapGpgKey")
	}

	/* Generate GPG key */
	publicKeyString, err := generateGpgKey(entity)
	if err != nil {
		panic("unable to generate GPG key in boostrapGpgKey")
	}

	/* Create GPG key*/
	request := &mreq.Gpg{
		Data: mreq.GpgData{
			Type: "gpg-keys",
			Attributes: mreq.GpgDataAttributes{
				AsciiArmor: publicKeyString,
				Namespace:  orgName,
			},
		},
	}
	cResp, err := cli.GpgService.Create(context.Background(), request)
	if err != nil {
		panic("unable to boostrap gpg key")
	}
	return &cResp
}

func bootstrapProvider(t *testing.T, cli *api.TerraformEnterpriseClient, orgName, providerName string) *mresp.Provider {
	request := mreq.Provider{
		Data: mreq.ProviderData{
			Type: "registry-providers",
			Attributes: mreq.ProviderDataAttributes{
				Name:         providerName,
				Namespace:    orgName,
				RegistryName: me.RegistryTypePrivate,
			},
		},
	}
	cResp, err := cli.ProviderService.Create(context.Background(), orgName, &request)
	if err != nil {
		panic("unable to boostrap provider")
	}
	return &cResp
}

func generateGpgKey(entity *openpgp.Entity) (string, error) {
	var publicKeyBuf bytes.Buffer
	err := entity.Serialize(&publicKeyBuf)
	if err != nil {
		fmt.Println("Error serializing public key:", err)
		return "", err
	}

	// Convert the public key to an armored string
	publicKeyArmorBuf := bytes.Buffer{}
	w, err := armor.Encode(&publicKeyArmorBuf, "PGP PUBLIC KEY BLOCK", nil)
	if err != nil {
		fmt.Println("Error encoding public key:", err)
		return "", err
	}
	_, err = w.Write(publicKeyBuf.Bytes())
	if err != nil {
		fmt.Println("Error writing public key to armored buffer:", err)
		return "", err
	}
	w.Close()

	return publicKeyArmorBuf.String(), nil
}

func getGoreleaserMetadataString() string {
	return `
{
	"project_name":"terraform-provider-tfepatch",
	"tag":"0.1.8",
	"previous_tag":"0.1.7",
	"version":"0.1.8",
	"commit":"d9e65cf460c1686df0741db13198709c8edb6c69",
	"date":"2023-04-18T18:05:24.933860632Z",
	"runtime":{
		"goos":"linux",
		"goarch":"amd64"
	}
}`
}

func getGoreleaserArtifactString() string {
	return `
	[
		{
		   "name":"terraform-provider-tfepatch_v0.1.8",
		   "path":"/home/runner/work/terraform-provider-tfepatch/terraform-provider-tfepatch/internal/dist/terraform-provider-tfepatch_linux_amd64_v1/terraform-provider-tfepatch_v0.1.8",
		   "goos":"linux",
		   "goarch":"amd64",
		   "goamd64":"v1",
		   "internal_type":4,
		   "type":"Binary",
		   "extra":{
			  "Binary":"terraform-provider-tfepatch_v0.1.8",
			  "Ext":"",
			  "ID":"terraform-provider-tfepatch"
		   }
		},
		{
		   "name":"terraform-provider-tfepatch_v0.1.8",
		   "path":"/home/runner/work/terraform-provider-tfepatch/terraform-provider-tfepatch/internal/dist/terraform-provider-tfepatch_darwin_amd64_v1/terraform-provider-tfepatch_v0.1.8",
		   "goos":"darwin",
		   "goarch":"amd64",
		   "goamd64":"v1",
		   "internal_type":4,
		   "type":"Binary",
		   "extra":{
			  "Binary":"terraform-provider-tfepatch_v0.1.8",
			  "Ext":"",
			  "ID":"terraform-provider-tfepatch"
		   }
		},
		{
		   "name":"terraform-provider-tfepatch_0.1.8_linux_amd64.zip",
		   "path":"/workspace/testing-assets/terraform-provider-tfepatch_0.1.8_linux_amd64.zip",
		   "goos":"linux",
		   "goarch":"amd64",
		   "goamd64":"v1",
		   "internal_type":1,
		   "type":"Archive",
		   "extra":{
			  "Binaries":[
				 "terraform-provider-tfepatch_v0.1.8"
			  ],
			  "Builds":[
				 {
					"name":"terraform-provider-tfepatch_v0.1.8",
					"path":"/home/runner/work/terraform-provider-tfepatch/terraform-provider-tfepatch/internal/dist/terraform-provider-tfepatch_linux_amd64_v1/terraform-provider-tfepatch_v0.1.8",
					"goos":"linux",
					"goarch":"amd64",
					"goamd64":"v1",
					"internal_type":4,
					"type":"Binary",
					"extra":{
					   "Binary":"terraform-provider-tfepatch_v0.1.8",
					   "Ext":"",
					   "ID":"terraform-provider-tfepatch"
					}
				 }
			  ],
			  "Checksum":"sha256:f79a37934c2e9d793ee874eeef921213e72069480c9df3953c907d79a0f5f034",
			  "Format":"zip",
			  "ID":"default",
			  "Replaces":null,
			  "WrappedIn":""
		   }
		},
		{
		   "name":"terraform-provider-tfepatch_0.1.8_darwin_amd64.zip",
		   "path":"/workspace/testing-assets/terraform-provider-tfepatch_0.1.8_darwin_amd64.zip",
		   "goos":"darwin",
		   "goarch":"amd64",
		   "goamd64":"v1",
		   "internal_type":1,
		   "type":"Archive",
		   "extra":{
			  "Binaries":[
				 "terraform-provider-tfepatch_v0.1.8"
			  ],
			  "Builds":[
				 {
					"name":"terraform-provider-tfepatch_v0.1.8",
					"path":"/home/runner/work/terraform-provider-tfepatch/terraform-provider-tfepatch/internal/dist/terraform-provider-tfepatch_darwin_amd64_v1/terraform-provider-tfepatch_v0.1.8",
					"goos":"darwin",
					"goarch":"amd64",
					"goamd64":"v1",
					"internal_type":4,
					"type":"Binary",
					"extra":{
					   "Binary":"terraform-provider-tfepatch_v0.1.8",
					   "Ext":"",
					   "ID":"terraform-provider-tfepatch"
					}
				 }
			  ],
			  "Checksum":"sha256:f2991fc425fbaa4033f10cf7962243832056762c9fbecb656db11e429ea9551e",
			  "Format":"zip",
			  "ID":"default",
			  "Replaces":null,
			  "WrappedIn":""
		   }
		},
		{
		   "name":"terraform-provider-tfepatch_0.1.8_SHA256SUMS",
		   "path":"/workspace/testing-assets/terraform-provider-tfepatch_0.1.8_SHA256SUMS",
		   "internal_type":12,
		   "type":"Checksum",
		   "extra":{}
		},
		{
		   "name":"terraform-provider-tfepatch_0.1.8_SHA256SUMS.sig",
		   "path":"/workspace/testing-assets/terraform-provider-tfepatch_0.1.8_SHA256SUMS.sig",
		   "internal_type":13,
		   "type":"Signature",
		   "extra":{
			  "ID":"default"
		   }
		},
		{
		   "name":"terraform-provider-tfepatch_0.1.8_manifest.json",
		   "path":"/workspace/testing-assets/terraform-registry-manifest.json",
		   "internal_type":3,
		   "type":"File"
		}
	 ]`
}

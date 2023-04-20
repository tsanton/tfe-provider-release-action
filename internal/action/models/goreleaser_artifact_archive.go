package models

import (
	"encoding/json"
)

/*
[
   {
      "name":"terraform-provider-tfepatch_0.1.8_linux_amd64.zip",
      "path":"dist/terraform-provider-tfepatch_0.1.8_linux_amd64.zip",
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
]
*/

type GoreleaserArtifactArchive struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Os      string `json:"goos"`
	Goarch  string `json:"goarch"`
	GoAmd64 string `json:"goamd64"`
	Type    string `json:"type"`
	ShaSum  string
}

func (a *GoreleaserArtifactArchive) UnmarshalJSON(data []byte) error {
	var aux struct {
		Name    string                 `json:"name"`
		Path    string                 `json:"path"`
		Os      string                 `json:"goos"`
		Goarch  string                 `json:"goarch"`
		GoAmd64 string                 `json:"goamd64"`
		Type    string                 `json:"type"`
		Extra   map[string]interface{} `json:"extra"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	a.Name = aux.Name
	a.Path = aux.Path
	a.Os = aux.Os
	a.Goarch = aux.Goarch
	a.GoAmd64 = aux.GoAmd64
	a.Type = aux.Type

	sha, ok := aux.Extra["Checksum"].(string)
	if ok {
		// a.ShaSum = strings.Split(sha, ":")[1]
		a.ShaSum = sha
	}

	return nil
}

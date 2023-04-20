package models

/*
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
}
*/

type GoreleaserMetadata struct {
	ProviderName string `json:"project_name"`
	Version      string `json:"version"`
}

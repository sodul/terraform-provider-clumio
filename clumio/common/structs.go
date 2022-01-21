package common

import clumioConfig "github.com/clumio-code/clumio-go-sdk/config"

// ApiClient defines the APIs/connections required by the resources.
type ApiClient struct {
	SnsAPI       SNSAPI
	S3API        S3API
	ClumioConfig clumioConfig.Config
}

// The payload in the status file read from S3.
type StatusObject struct {
	Status string            `json:"Status"`
	Reason *string           `json:"Reason,omitempty"`
	Data   map[string]string `json:"Data,omitempty"`
}

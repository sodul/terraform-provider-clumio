// Copyright 2021. Clumio, Inc.

// This file contains the functions related to provider definition and initialization.

package clumio

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	awsCredentials "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

// New is the factory method that returns a function which, when called,
// creates a new Provider instance.
func New(isTest bool) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"clumio_callback_resource": clumioCallback(),
			},
		}
		p.ConfigureContextFunc = configure(p, isTest)
		p.Schema = map[string]*schema.Schema{
			"region" : {
				Type: schema.TypeString,
				Optional: true,
				Description: "AWS Region.",
			},
			"clumio_region" : {
				Type: schema.TypeString,
				Optional: true,
				Description: "Clumio Control Plane AWS Region.",
			},
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AWS Access Key.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AWS Secret Key.",
			},
			"session_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "AWS Session Token.",
			},
		}
		return p
	}
}

// apiClient defines the APIs/connections required by the resources.
type apiClient struct {
	snsAPI SNSAPI
	s3API S3API
}

// configure is a factory method to configure the Provider.
func configure(_ *schema.Provider, isTest bool) func(context.Context,
	*schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{},
	diag.Diagnostics) {
		if isTest{
			return &apiClient{
				snsAPI: SNSClient{},
				s3API:  S3Client{},
			}, nil
		}

		accessKey := getStringValue(d, "access_key")
		secretKey := getStringValue(d, "secret_key")
		sessionToken := getStringValue(d, "session_token")

		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRetryer(func() aws.Retryer {
			return retry.AddWithMaxAttempts(retry.NewStandard(),kMaxRetries)
		}))
		if err != nil {
			return nil, diag.Errorf("Error loading default config for AWS Provider: %v",
				err)
		}
		region := getStringValue(d, "region")
		if region != "" {
			cfg.Region = region
		}

		if accessKey != "" && (secretKey != "" || sessionToken != ""){
			cfg.Credentials = awsCredentials.NewStaticCredentialsProvider(accessKey,
				secretKey, sessionToken)
		}
		regionalSns := sns.NewFromConfig(cfg)
		clumioRegion := getStringValue(d, "clumio_region")
		if clumioRegion != "" {
			cfg.Region = clumioRegion
		}
		s3obj := s3.NewFromConfig(cfg)
		return &apiClient{
			snsAPI: regionalSns,
			s3API: s3obj,
		}, nil
	}
}

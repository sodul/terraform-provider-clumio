// Copyright 2021. Clumio, Inc.

// clumio_callback_resource definition and CRUD implementation.

package clumio

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/smithy-go"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	keyRegion           = "Region"
	keyToken            = "Token"
	keyType             = "Type"
	keyAccountID        = "AccountId"
	keyRoleID           = "RoleId"
	keyRoleArn          = "RoleArn"
	keyExternalID       = "RoleExternalId"
	keyClumioEventPubID = "ClumioEventPubId"
	keyCanonicalUser    = "CanonicalUser"
	keyTemplateConfig   = "TemplateConfiguration"

	// Number of retries that we will perform before giving up a AWS request.
	kMaxRetries = 8
	requestTypeCreate = "Create"
	requestTypeDelete = "Delete"
	requestTypeUpdate = "Update"

	//Status strings
	statusFailed        = "FAILED"
	statusSuccess       = "SUCCESS"

	//bucket key format
	bucketKeyFormat = "acmtfstatus/%s/clumio-status.json"
)

var (
	// protectInfoMap is the mapping of the the datasource to the resource parameter and
	// if a config section is required, then isConfig will be true.
	protectInfoMap = map[string]sourceConfigInfo{
		"ebs": {
			sourceKey: "protect_ebs_version",
			isConfig: false,
		},
		"rds": {
			sourceKey: "protect_rds_version",
			isConfig: false,
		},
		"ec2_mssql": {
			sourceKey: "protect_ec2_mssql_version",
			isConfig: false,
		},
		"warm_tier": {
			sourceKey: "protect_warm_tier_version",
			isConfig: true,
		},
	}
	// warmtierInfoMap is the mapping of the the warm tier datasource to the resource
	// parameter and if a config section is required, then isConfig will be true.
	warmtierInfoMap = map[string]sourceConfigInfo{
		"dynamodb": {
			sourceKey: "protect_warm_tier_dynamodb_version",
			isConfig: false,
		},
	}
)

// SNSEvent is the event payload to be sent to the topic
type SNSEvent struct {
	RequestType        string                 `json:"RequestType"`
	ServiceToken       string                 `json:"ServiceToken"`
	ResponseURL        string                 `json:"ResponseURL"`
	StackID            string                 `json:"StackId"`
	RequestID          string                 `json:"RequestId"`
	LogicalResourceID  string                 `json:"LogicalResourceId"`
	ResourceType       string                 `json:"ResourceType"`
	ResourceProperties map[string]interface{} `json:"ResourceProperties"`
}

// The payload in the status file read from S3.
type StatusObject struct {
	Status             string            `json:"Status"`
	Reason             *string           `json:"Reason,omitempty"`
	Data               map[string]string `json:"Data,omitempty"`
}

// clumioCallback returns the resource for Clumio Callback. This resource is similar to
// the cloud formation custom resource. It will publish an event to the specified SNS
// topic and then wait for the status payload in the given S3 bucket.
func clumioCallback() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Callback Resource used while on-boarding AWS clients.",

		CreateContext: clumioCallbackCreate,
		ReadContext:   clumioCallbackRead,
		UpdateContext: clumioCallbackUpdate,
		DeleteContext: clumioCallbackDelete,

		Schema: map[string]*schema.Schema{
			"sns_topic": {
				Description: "SNS Topic to publish event.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"token": {
				Type:        schema.TypeString,
				Description: "The AWS integration ID token.",
				Required: true,
			},
			"role_external_id": {
				Type:        schema.TypeString,
				Description: "A key that must be used by Clumio to assume the service role" +
					" in your account. This should be a secure string, like a password," +
					" but it does not need to be remembered (random characters are best).",
				Required: true,
			},
			"account_id": {
				Type:        schema.TypeString,
				Description: "The AWS Customer Account ID.",
				Required: true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "The AWS Region.",
				Required: true,
			},
			"role_id": {
				Type:        schema.TypeString,
				Description: "Clumio IAM Role ID.",
				Required: true,
			},
			"role_arn": {
				Type:        schema.TypeString,
				Description: "Clumio IAM Role Arn.",
				Required: true,
			},
			"clumio_event_pub_id": {
				Type:        schema.TypeString,
				Description: "Clumio Event Pub SNS topic ID.",
				Required: true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Registration Type.",
				Required: true,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Description: "S3 bucket name where the status file is written.",
				Required: true,
			},
			"canonical_user": {
				Type:        schema.TypeString,
				Description: "Canonical User ID of the account.",
				Required: true,
			},
			"config_version": {
				Type: schema.TypeString,
				Description: "Clumio Config version.",
				Required: true,
			},
			"discover_enabled": {
				Type: schema.TypeBool,
				Description: "Is Clumio Discover enabled.",
				Optional: true,
			},
			"discover_version": {
				Type: schema.TypeString,
				Description: "Clumio Discover version.",
				Required: true,
			},
			"protect_enabled": {
				Type: schema.TypeBool,
				Description: "Is Clumio Protect enabled.",
				Optional: true,
			},
			"protect_config_version": {
				Type: schema.TypeString,
				Description: "Clumio Protect Config version.",
				Optional: true,
			},
			"protect_ebs_version": {
				Type: schema.TypeString,
				Description: "Clumio EBS Protect version.",
				Optional: true,
			},
			"protect_rds_version": {
				Type: schema.TypeString,
				Description: "Clumio RDS Protect version.",
				Optional: true,
			},
			"protect_ec2_mssql_version": {
				Type: schema.TypeString,
				Description: "Clumio EC2 MSSQL Protect version.",
				Optional: true,
			},
			"protect_s3_version": {
				Type: schema.TypeString,
				Description: "Clumio S3 Protect version.",
				Optional: true,
			},
			"protect_warm_tier_version": {
				Type: schema.TypeString,
				Description: "Clumio Warmtier Protect version.",
				Optional: true,
			},
			"protect_warm_tier_dynamodb_version": {
				Type: schema.TypeString,
				Description: "Clumio DynamoDB Warmtier Protect version.",
				Optional: true,
			},
			"properties": {
				Type: schema.TypeMap,
				Description: "Properties to be passed in the SNS event.",
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// clumioCallbackCreate handles the Create action for the Clumio Callback Resource.
func clumioCallbackCreate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return clumioCallbackCommon(ctx, d, meta, requestTypeCreate)
}

// clumioCallbackCreate handles the Read action for the Clumio Callback Resource.
func clumioCallbackRead(_ context.Context, _ *schema.ResourceData,
	_ interface{}) diag.Diagnostics {
	// Nothing to Read for this resource
	return nil
}

// clumioCallbackCreate handles the Update action for the Clumio Callback Resource.
func clumioCallbackUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return clumioCallbackCommon(ctx, d, meta, requestTypeUpdate)
}

// clumioCallbackCreate handles the Delete action for the Clumio Callback Resource.
func clumioCallbackDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return clumioCallbackCommon(ctx, d, meta, requestTypeDelete)
}

// clumioCallbackCommon will construct the event payload from the resource properties and
// publish the event to the SNS topic and then wait for the status payload in the given
// S3 bucket.
func clumioCallbackCommon(ctx context.Context, d *schema.ResourceData, meta interface{},
	eventType string) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*apiClient)
	bucketName := d.Get("bucket_name").(string)
	accountId := fmt.Sprintf("%v", d.Get("account_id"))

	regionalSns := client.snsAPI
	event := SNSEvent{
		RequestType:        eventType,
		ServiceToken:       fmt.Sprintf("%v", d.Get("token")),
		ResourceProperties: nil,
	}
	resourceProperties := make(map[string]interface{})
	resourceProperties[keyAccountID] = fmt.Sprintf("%v", d.Get("account_id"))
	resourceProperties[keyToken] =fmt.Sprintf("%v", d.Get("token"))
	resourceProperties[keyType] = fmt.Sprintf("%v", d.Get("type"))
	resourceProperties[keyAccountID] = accountId
	resourceProperties[keyRegion] = fmt.Sprintf("%v", d.Get("region"))
	resourceProperties[keyRoleID] = fmt.Sprintf("%v", d.Get("role_id"))
	resourceProperties[keyRoleArn] = fmt.Sprintf("%v", d.Get("role_arn"))
	resourceProperties[keyExternalID] =
		fmt.Sprintf("%v", d.Get("role_external_id"))
	resourceProperties[keyClumioEventPubID] =
		fmt.Sprintf("%v", d.Get("clumio_event_pub_id"))
	resourceProperties[keyCanonicalUser] = fmt.Sprintf("%v", d.Get("canonical_user"))

	templateConfigs := getTemplateConfiguration(d)
	resourceProperties[keyTemplateConfig] = templateConfigs
	if val, ok := d.GetOk("properties"); ok && len(val.(map[string]interface{})) > 0 {
		properties := val.(map[string]interface{})
		for key, value := range properties {
			resourceProperties[key] = value.(string)
		}
	}
	event.ResourceProperties = resourceProperties
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return diag.Errorf("Error occurred in marshalling event: %v", err)
	}
	// Publish event to SNS.
	publishInput := &sns.PublishInput{
		Message:  aws.String(string(eventBytes)),
		TopicArn: aws.String(fmt.Sprintf("%v", d.Get("sns_topic"))),
	}
	startTime := time.Now()
	_, err = regionalSns.Publish(ctx, publishInput)
	if err != nil {
		return diag.Errorf("Error occurred in SNS Publish Request: %v", err)
	}
	if eventType == requestTypeCreate {
		d.SetId(uuid.New().String())
	}
	s3obj := client.s3API
	endTime := startTime.Add(5 * time.Minute)
	timeOut := false
	processingDone := false
	// Poll the s3 bucket for the clumio-status.json file. Keep retrying every 5 seconds
	// till the last modified time on the file is greater than the startTime and less than
	// the end time.
	for {
		if time.Now().After(endTime) {
			timeOut = true
			break
		}
		time.Sleep(5 * time.Second)
		// HeadObject call to get the last modified time of the file.
		headObject, err := s3obj.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fmt.Sprintf(bucketKeyFormat, accountId)),
		})
		if err != nil {
			var aerr smithy.APIError
			if errors.As(err, &aerr) {
				if aerr.ErrorCode() == "Forbidden" {
					log.Println(aerr.Error())
					continue
				}
				return diag.Errorf("Error retrieving metadata of clumio-status.json. "+
					"Error Code : %v, Error message: %v, origError: %v",
					aerr.ErrorCode(),
					aerr.ErrorMessage(), err)
			}
			return diag.Errorf("Error retrieving metadata of clumio-status.json: %v", err)
		} else if headObject.LastModified.After(startTime) {
			// Get the clumio-status.json object.
			statusObj, err := s3obj.GetObject(ctx, &s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(fmt.Sprintf(bucketKeyFormat, accountId)),
			})
			if err != nil {
				return diag.Errorf("Error retrieving clumio-status.json: %v", err)
			}

			var status StatusObject
			statusObjBytes := new(bytes.Buffer)
			_, err = statusObjBytes.ReadFrom(statusObj.Body)
			if err != nil {
				return diag.Errorf("Error reading status object: %v", err)
			}
			err = json.Unmarshal(statusObjBytes.Bytes(), &status)
			if err != nil {
				return diag.Errorf("Error unmarshalling status object: %v", err)
			}
			if status.Status == statusFailed {
				return diag.Errorf("Processing of Clumio Event failed. "+
					"Error Message : %s", *status.Reason)
			} else if status.Status == statusSuccess {
				processingDone = true
				break
			}
		}
	}
	if !processingDone && timeOut {
		return diag.Errorf("Timeout occurred waiting for status.")
	}
	return nil
}

type sourceConfigInfo struct {
	sourceKey string
	isConfig bool
}
// getTemplateConfiguration returns the template configuration.
func getTemplateConfiguration(d *schema.ResourceData) map[string]interface{} {
	templateConfigs := make(map[string]interface{})
	configMap := getConfigMapForKey(d, "config_version", false)
	if configMap == nil{
		return templateConfigs
	}
	templateConfigs["config"] = configMap
	discoverMap := getConfigMapForKey(d, "discover_version", true)
	if discoverMap == nil {
		return templateConfigs
	}
	templateConfigs["discover"] = discoverMap
	protectMap := getConfigMapForKey(d, "protect_config_version", true)
	if protectMap == nil {
		return templateConfigs
	}
	populateConfigMap(d, protectInfoMap, protectMap)

	if protectWarmtierMap, ok := protectMap["warm_tier"]; ok {
		populateConfigMap(d, warmtierInfoMap, protectWarmtierMap.(map[string]interface{}))
	}
	templateConfigs["protect"] = protectMap
	return templateConfigs
}

func populateConfigMap(d *schema.ResourceData, configInfoMap map[string]sourceConfigInfo,
	configMap map[string]interface{}) {
	for source, sourceInfo := range configInfoMap {
		protectSourceMap := getConfigMapForKey(d, sourceInfo.sourceKey, sourceInfo.isConfig)
		if protectSourceMap != nil {
			configMap[source] = protectSourceMap
		}
	}
}

// getConfigMapForKey returns a config map for the key if it exists in ResourceData.
func getConfigMapForKey(
	d *schema.ResourceData, key string, isConfig bool) map[string]interface{}{
	var mapToReturn map[string]interface{}
	if val, ok := d.GetOk(key); ok{
		keyMap := make(map[string]interface{})
		if keyVersion, ok := val.(string); ok {
			keyMap["enabled"] = true
			keyMap["version"] = keyVersion
		}
		mapToReturn = keyMap
		// If isConfig is true it wraps the keyMap with another map with "config" as the key.
		if isConfig{
			configMap := make(map[string]interface{})
			configMap["config"] = keyMap
			mapToReturn = configMap
		}
	}
	return mapToReturn
}

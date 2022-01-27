// Copyright 2021. Clumio, Inc.

// clumio_policy resource definition and CRUD implementation.
package clumio_policy

import (
	"context"
	policyDefinitions "github.com/clumio-code/clumio-go-sdk/controllers/policy_definitions"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	resBackupWindow = &schema.Resource{
		Schema: map[string]*schema.Schema{
			schemaStartTime: {
				Type: schema.TypeString,
				Description: "The time when the backup window opens." +
					" Specify the start time in the format hh:mm," +
					" where hh represents the hour of the day and" +
					" mm represents the minute of the day based on" +
					" the 24 hour clock.",
				Required: true,
			},
			schemaEndTime: {
				Type: schema.TypeString,
				Description: "The time when the backup window closes." +
					" Specify the end time in the format hh:mm," +
					" where hh represents the hour of the day and" +
					" mm represents the minute of the day based on" +
					" the 24 hour clock.",
				Required: true,
			},
		},
	}

	resUnitValue = &schema.Resource{
		Schema: map[string]*schema.Schema{
			schemaUnit: {
				Type:     schema.TypeString,
				Required: true,
				Description: "The measurement unit of the SLA parameter. Values include" +
					" hours, days, months, and years.",
			},
			schemaValue: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The measurement value of the SLA parameter.",
			},
		},
	}

	resSla = &schema.Resource{
		Schema: map[string]*schema.Schema{
			schemaRetentionDuration: {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Description: "The retention time for this SLA. " +
					"For example, to retain the backup for 1 month," +
					" set unit=months and value=1.",
				Set:  schema.HashResource(resUnitValue),
				Elem: resUnitValue,
			},
			schemaRpoFrequency: {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Required: true,
				Description: "The minimum frequency between " +
					"backups for this SLA. Also known as the " +
					"recovery point objective (RPO) interval. For" +
					" example, to configure the minimum frequency" +
					" between backups to be every 2 days, set " +
					"unit=days and value=2. To configure the SLA " +
					"for on-demand backups, set unit=on_demand " +
					"and leave the value field empty.",
				Set:  schema.HashResource(resUnitValue),
				Elem: resUnitValue,
			},
		},
	}

	resOperation = &schema.Resource{
		Schema: map[string]*schema.Schema{
			schemaActionSetting: {
				Type: schema.TypeString,
				Description: "Determines whether the policy should take action" +
					" now or during the specified backup window. Valid values:" +
					"immediate: to start backup process immediately" +
					"window: to start backup in the specified window",
				Required: true,
			},
			schemaOperationType: {
				Type: schema.TypeString,
				Description: "The operation to be performed for this SLA set." +
					"Each SLA set corresponds to one and only one operation.",
				Required: true,
			},
			schemaBackupWindow: {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Description: "The start and end times for the customized" +
					" backup window.",
				Set:  schema.HashResource(resBackupWindow),
				Elem: resBackupWindow,
			},
			schemaSlas: {
				Type:     schema.TypeSet,
				Required: true,
				Description: "The service level agreement (SLA) for the policy." +
					" A policy can include one or more SLAs. For example, " +
					"a policy can retain daily backups for a month each, " +
					"and monthly backups for a year each.",
				Set:  schema.HashResource(resSla),
				Elem: resSla,
			},
		},
	}
)

// ClumioPolicy returns the resource for Clumio Policy Definition.
func ClumioPolicy() *schema.Resource {

	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Policy Resource used to schedule backups on Clumio supported" +
			" data sources.",

		CreateContext: clumioPolicyCreate,
		ReadContext:   clumioPolicyRead,
		UpdateContext: clumioPolicyUpdate,
		DeleteContext: clumioPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			schemaId: {
				Description: "Policy Id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			schemaLockStatus: {
				Description: "Policy Lock Status.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			schemaName: {
				Description: "The unique name of the policy.",
				Type:        schema.TypeString,
				Required:    true,
			},
			schemaActivationStatus: {
				Type: schema.TypeString,
				Description: "The status of the policy. Valid values are:" +
					"activated: Backups will take place regularly according to the policy SLA." +
					"deactivated: Backups will not begin until the policy is reactivated." +
					" The assets associated with the policy will have their compliance" +
					" status set to deactivated.",
				Optional: true,
				Computed: true,
			},
			schemaOrganizationalUnitId: {
				Type: schema.TypeString,
				Description: "The Clumio-assigned ID of the organizational unit" +
					" associated with the policy.",
				Optional: true,
				Computed: true,
			},
			schemaOperations: {
				Type:     schema.TypeSet,
				Required: true,
				Set:      schema.HashResource(resOperation),
				Elem:     resOperation,
			},
		},
	}
}

// clumioPolicyCreate handles the Create action for the Clumio Policy Resource.
func clumioPolicyCreate(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	pd := policyDefinitions.NewPolicyDefinitionsV1(client.ClumioConfig)
	activationStatus := common.GetStringValue(d, schemaActivationStatus)
	name := common.GetStringValue(d, schemaName)
	operationsVal, ok := d.GetOk(schemaOperations)
	if !ok {
		return diag.Errorf("Operations is a required attribute")
	}
	policyOperations, _ := mapSchemaOperationsToClumioOperations(operationsVal)
	orgUnitId := common.GetStringValue(d, schemaOrganizationalUnitId)
	pdRequest := &models.CreatePolicyDefinitionV1Request{
		ActivationStatus:     &activationStatus,
		Name:                 &name,
		Operations:           policyOperations,
		OrganizationalUnitId: &orgUnitId,
	}
	res, apiErr := pd.CreatePolicyDefinition(pdRequest)
	if apiErr != nil {
		return diag.Errorf("Error creating policy definition %v. Error: %v",
			d.Id(), string(apiErr.Response))
	}
	d.SetId(*res.Id)
	return nil
}

// clumioPolicyRead handles the Read action for the Clumio Policy Resource.
func clumioPolicyRead(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	pd := policyDefinitions.NewPolicyDefinitionsV1(client.ClumioConfig)
	res, apiErr := pd.ReadPolicyDefinition(d.Id(), nil)
	if apiErr != nil {
		return diag.Errorf("Error retrieving policy definition %v. Error: %v",
			d.Id(), string(apiErr.Response))
	}
	err := d.Set(schemaLockStatus, *res.LockStatus)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaLockStatus, err)
	}
	err = d.Set(schemaName, *res.Name)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaName, err)
	}
	if res.ActivationStatus != nil {
		err = d.Set(schemaActivationStatus, *res.ActivationStatus)
		if err != nil {
			return diag.Errorf(common.SchemaAttributeSetError, schemaActivationStatus, err)
		}
	}
	if res.OrganizationalUnitId != nil {
		err = d.Set(schemaOrganizationalUnitId, *res.OrganizationalUnitId)
		if err != nil {
			return diag.Errorf(common.SchemaAttributeSetError, schemaOrganizationalUnitId, err)
		}
	}
	err = d.Set(schemaOperations, mapClumioOperationsToSchemaOperations(res.Operations))
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaOperations, err)
	}
	return nil
}

// clumioPolicyUpdate handles the Update action for the Clumio Policy Resource.
func clumioPolicyUpdate(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	pd := policyDefinitions.NewPolicyDefinitionsV1(client.ClumioConfig)
	activationStatus := common.GetStringValue(d, schemaActivationStatus)
	name := common.GetStringValue(d, schemaName)
	operationsVal, ok := d.GetOk(schemaOperations)
	if !ok {
		return diag.Errorf("Operations is a required attribute")
	}
	policyOperations, _ := mapSchemaOperationsToClumioOperations(operationsVal)
	orgUnitId := common.GetStringValue(d, schemaOrganizationalUnitId)
	pdRequest := &models.UpdatePolicyDefinitionV1Request{
		ActivationStatus:     &activationStatus,
		Name:                 &name,
		Operations:           policyOperations,
		OrganizationalUnitId: &orgUnitId,
	}
	res, apiErr := pd.UpdatePolicyDefinition(d.Id(), nil, pdRequest)
	if apiErr != nil {
		diag.Errorf("Error updating Policy Definition %v. Error: %v",
			d.Id(), string(apiErr.Response))
	}
	d.SetId(*res.Id)
	return nil
}

// clumioPolicyDelete handles the Delete action for the Clumio Policy Resource.
func clumioPolicyDelete(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	pd := policyDefinitions.NewPolicyDefinitionsV1(client.ClumioConfig)
	_, apiErr := pd.DeletePolicyDefinition(d.Id())
	if apiErr != nil {
		return diag.Errorf("Error deleting policy definition %v. Error: %v",
			d.Id(), string(apiErr.Response))
	}
	return nil
}

// mapSchemaOperationsToClumioOperations maps the schema operations to the Clumio API
// request operations.
func mapSchemaOperationsToClumioOperations(
	operations interface{}) ([]*models.PolicyOperation, error) {
	operationsSlice := operations.(*schema.Set).List()
	policyOperations := make([]*models.PolicyOperation, 0)
	for _, operation := range operationsSlice {
		operationAttrMap := operation.(map[string]interface{})
		actionSetting := operationAttrMap[schemaActionSetting].(string)
		operationType := operationAttrMap[schemaOperationType].(string)
		backupWindowIface, ok := operationAttrMap[schemaBackupWindow]
		var backupWindow *models.BackupWindow
		schemaBackupWindowSlice := backupWindowIface.(*schema.Set).List()
		if ok && len(schemaBackupWindowSlice) > 0 {
			schemaBackupWindow := schemaBackupWindowSlice[0].(map[string]interface{})
			schemaBackupWindowStartTime := schemaBackupWindow[schemaStartTime].(string)
			schemaBackupWindowEndTime := schemaBackupWindow[schemaEndTime].(string)
			backupWindow = &models.BackupWindow{
				EndTime:   &schemaBackupWindowEndTime,
				StartTime: &schemaBackupWindowStartTime,
			}
		}
		backupSLAs := make([]*models.BackupSLA, 0)
		slasIface := operationAttrMap[schemaSlas]
		schemaSlas := slasIface.(*schema.Set).List()
		for _, slaIface := range schemaSlas {
			schemaSla := slaIface.(map[string]interface{})
			retDurationIface := schemaSla[schemaRetentionDuration]
			var retentionDuration *models.RetentionBackupSLAParam
			schemaRetDurationSlice := retDurationIface.(*schema.Set).List()
			schemaRetDuration := schemaRetDurationSlice[0].(map[string]interface{})
			unit := schemaRetDuration[schemaUnit].(string)
			value := int64(schemaRetDuration[schemaValue].(int))
			retentionDuration = &models.RetentionBackupSLAParam{
				Unit:  &unit,
				Value: &value,
			}
			var rpoFrequency *models.RPOBackupSLAParam
			rpoFrequencyIface := schemaSla[schemaRpoFrequency]
			schemaRpoFrequencySlice := rpoFrequencyIface.(*schema.Set).List()
			schemaRpoFrequency := schemaRpoFrequencySlice[0].(map[string]interface{})
			rpoUnit := schemaRpoFrequency[schemaUnit].(string)
			rpoValue := int64(schemaRpoFrequency[schemaValue].(int))
			rpoFrequency = &models.RPOBackupSLAParam{
				Unit:  &rpoUnit,
				Value: &rpoValue,
			}
			backupSLA := &models.BackupSLA{
				RetentionDuration: retentionDuration,
				RpoFrequency:      rpoFrequency,
			}
			backupSLAs = append(backupSLAs, backupSLA)
		}

		policyOperation := &models.PolicyOperation{
			ActionSetting: &actionSetting,
			BackupWindow:  backupWindow,
			Slas:          backupSLAs,
			ClumioType:    &operationType,
		}
		policyOperations = append(policyOperations, policyOperation)
	}
	return policyOperations, nil
}

func mapClumioOperationsToSchemaOperations(operations []*models.PolicyOperation) interface{} {
	schemaOperations := &schema.Set{F: schema.HashResource(resOperation)}
	for _, operation := range operations {
		operationAttrMap := make(map[string]interface{})
		operationAttrMap[schemaActionSetting] = *operation.ActionSetting
		operationAttrMap[schemaOperationType] = *operation.ClumioType

		if operation.BackupWindow != nil {
			backupWindowMap := make(map[string]interface{})
			backupWindowMap[schemaStartTime] = *operation.BackupWindow.StartTime
			backupWindowMap[schemaEndTime] = *operation.BackupWindow.EndTime
			backupWindowSet := &schema.Set{F: schema.HashResource(resBackupWindow)}
			backupWindowSet.Add(backupWindowMap)
			operationAttrMap[schemaBackupWindow] = backupWindowSet
		}

		backupSlas := &schema.Set{F: schema.HashResource(resSla)}
		for _, sla := range operation.Slas {
			backupSla := make(map[string]interface{})

			backupSlaRetentionDuration := make(map[string]interface{})
			backupSlaRetentionDuration[schemaUnit] = *sla.RetentionDuration.Unit
			backupSlaRetentionDuration[schemaValue] = int(*sla.RetentionDuration.Value)
			backupSlaRetentionDurationSet := &schema.Set{F: schema.HashResource(resUnitValue)}
			backupSlaRetentionDurationSet.Add(backupSlaRetentionDuration)
			backupSla[schemaRetentionDuration] = backupSlaRetentionDurationSet

			backupSlaRpoFrequency := make(map[string]interface{})
			backupSlaRpoFrequency[schemaUnit] = *sla.RpoFrequency.Unit
			backupSlaRpoFrequency[schemaValue] = int(*sla.RpoFrequency.Value)
			backupSlaRpoFrequencySet := &schema.Set{F: schema.HashResource(resUnitValue)}
			backupSlaRpoFrequencySet.Add(backupSlaRpoFrequency)
			backupSla[schemaRpoFrequency] = backupSlaRpoFrequencySet

			backupSlas.Add(backupSla)
		}
		operationAttrMap[schemaSlas] = backupSlas

		schemaOperations.Add(operationAttrMap)
	}

	return schemaOperations
}

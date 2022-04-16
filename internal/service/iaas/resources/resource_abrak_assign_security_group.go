package resources

import (
	"context"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/arvancloud/terraform-provider-arvan/internal/service/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrakAssignSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrakAssignSecurityGroupCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceAbrakAssignSecurityGroupUpdate,
		DeleteContext: resourceAbrakAssignSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "region code",
				ValidateFunc: validation.StringInSlice(iaas.AvailableRegions, false),
			},
			"abrak_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "uuid of abrak",
				ValidateFunc: validation.IsUUID,
			},
			"security_group_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "uuid of security group",
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceAbrakAssignSecurityGroupCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	uuid := data.Get("abrak_uuid").(string)

	securityGroupId := data.Get("security_group_uuid").(string)
	err := c.Server.Actions.AssignSecurityGroup(region, uuid, securityGroupId)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not add security group %v server %v", securityGroupId, uuid),
		})
		return errors
	}

	data.SetId(securityGroupId)
	return errors
}

func resourceAbrakAssignSecurityGroupUpdate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	if data.HasChange("security_group_uuid") {
		abrakId := data.Get("abrak_uuid").(string)
		err := c.Server.Actions.RemoveSecurityGroup(region, abrakId, data.Id())
		if err != nil {
			errors = append(errors, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("could not remove security group %v server %v", data.Id(), abrakId),
			})
			return errors
		}

		return resourceAbrakAssignSecurityGroupCreate(ctx, data, meta)
	}
	return nil
}

func resourceAbrakAssignSecurityGroupDelete(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	abrakId := data.Get("abrak_uuid").(string)

	err := c.Server.Actions.RemoveSecurityGroup(region, abrakId, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not remove security group %v server %v", data.Id(), abrakId),
		})
		return errors
	}

	return errors
}

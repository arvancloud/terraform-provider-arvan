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

func ResourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityGroupCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceSecurityGroupUpdate,
		DeleteContext: resourceSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Region code",
				ValidateFunc: validation.StringInSlice(iaas.AvailableRegions, false),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of security group",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "description of security group",
			},
			"real_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "real name of security group",
			},
		},
	}
}

func resourceSecurityGroupCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not parse region",
		})
		return errors
	}

	// SecurityGroup Options
	SecurityGroup := &iaas.SecurityGroupOpts{
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
	}

	response, err := c.SecurityGroup.Create(region, SecurityGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.ID)
	return errors
}

func resourceSecurityGroupUpdate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	if data.HasChange("name") {
		// TODO: do we need to delete the previous one ?
		return resourceSecurityGroupCreate(ctx, data, meta)
	}
	return nil
}

func resourceSecurityGroupDelete(_ context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	err := c.SecurityGroup.Delete(region, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not delete SecurityGroup %v", data.Id()),
		})
		return errors
	}

	return nil
}

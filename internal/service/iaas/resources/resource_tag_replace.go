package resources

import (
	"context"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/arvancloud/terraform-provider-arvan/internal/service/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceTagReplaceBatch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagReplaceBatchCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: helper.DummyResourceAction,
		DeleteContext: resourceTagReplaceBatchDelete,
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
			"instance_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "list of UUID instance",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					Description:  "UUID of instance",
					ValidateFunc: validation.IsUUID,
				},
			},
			"tag_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "list of tag names",
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "tag name",
				},
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "type of instance",
				ValidateFunc: validation.StringInSlice(iaas.SupportedTagTypes, false),
			},
		},
	}
}

func resourceTagReplaceBatchCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	items := data.Get("instance_list").([]any)
	// instanceList
	var instanceList = make([]string, len(items))
	for i, instance := range items {
		instanceList[i] = instance.(string)
	}

	items = data.Get("tag_list").([]any)
	// tagList
	var tagList = make([]string, len(items))
	for i, instance := range items {
		tagList[i] = instance.(string)
	}

	// tagReplaceOpts options
	tagReplaceOpts := &iaas.TagReplaceOpts{
		InstanceList: instanceList,
		TagList:      tagList,
		InstanceType: data.Get("instance_type").(string),
	}

	err := c.Tag.Replace(region, tagReplaceOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(helper.GenUUID())
	return errors
}

func resourceTagReplaceBatchDelete(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	// TODO: we have to remove ONLY the inserted tags via this resource
	return nil
}

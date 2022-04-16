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

func ResourceAbrakSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrakSnapshotCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceAbrakSnapshotUpdate,
		DeleteContext: helper.DummyResourceAction,
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
			"snapshot_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "snapshot name of abrak",
			},
			"snapshot_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "snapshot description of abrak",
			},
		},
	}
}

func resourceAbrakSnapshotCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
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

	snapshotName := data.Get("snapshot_name").(string)
	snapshotDescription := data.Get("snapshot_description").(string)
	err := c.Server.Actions.Snapshot(region, uuid, snapshotName, snapshotDescription)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not take snapshot server %v to %v", uuid, snapshotName),
		})
		return errors
	}

	data.SetId(uuid)
	return errors
}

func resourceAbrakSnapshotUpdate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	if data.HasChanges("snapshot_name") {
		return resourceAbrakSnapshotCreate(ctx, data, meta)
	}
	return nil
}

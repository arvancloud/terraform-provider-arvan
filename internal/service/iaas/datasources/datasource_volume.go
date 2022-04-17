package datasources

import (
	"context"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DatasourceVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceVolumeRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "region code",
				ValidateFunc: validation.StringInSlice(iaas.AvailableRegions, false),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of volume",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "size of volume",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "status of volume",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "created_at of volume",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "description of volume",
			},
			"volume_type_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of volume",
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "snapshot ID of volume",
			},
			"source_volume_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "source id of volume",
			},
			"bootable": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "volume is bootable or not",
			},
		},
	}
}

func datasourceVolumeRead(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	name := data.Get("name").(string)
	volume, err := c.Volume.Find(region, name)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("volume %v not found", name),
		})
		return errors
	}

	data.SetId(volume.ID)
	data.Set("name", volume.Name)
	data.Set("size", volume.Size)
	data.Set("status", volume.Status)
	data.Set("created_at", volume.CreatedAt)
	data.Set("description", volume.Description)
	data.Set("volume_type_name", volume.VolumeTypeName)
	data.Set("snapshot_id", volume.SnapshotId)
	data.Set("source_volume_id", volume.SourceVolumeId)
	data.Set("bootable", volume.Bootable)
	return errors
}

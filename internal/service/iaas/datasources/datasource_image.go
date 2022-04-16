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

func DatasourceImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceImageRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "region code",
				ValidateFunc: validation.StringInSlice(iaas.AvailableRegions, false),
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "type of image",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of image",
			},
		},
	}
}

func datasourceImageRead(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	imageName := data.Get("name").(string)
	imageType := data.Get("type").(string)
	image, err := c.Image.Find(region, imageName, imageType)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("image %v not found", imageName),
		})
		return errors
	}

	switch image.(type) {
	case iaas.ImageDetails:
		data.SetId(image.(iaas.ImageDetails).ID)
	case iaas.ImageServerDetails:
		data.SetId(image.(iaas.ImageServerDetails).ID)
	}

	return errors
}

package iaas

import (
	"context"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAbrak() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrakCreate,
		ReadContext:   resourceAbrakRead,
		UpdateContext: resourceAbrakUpdate,
		DeleteContext: resourceAbrakDelete,
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
				Description: "name of abrak",
			},
			"networks": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Network(s) of abrak",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Plan ID of abrak, you can get list of plan IDs of each region from sizes api",
			},
			"image": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Required:    true,
				Description: "image of abrak",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "type of image",
							ValidateFunc: validation.StringInSlice([]string{
								iaas.ImageTypeServer,
								iaas.ImageTypeSnapshot,
								iaas.ImageTypeDistributions,
							}, false),
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "name of image",
						},
					},
				},
			},
			"security_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "array of security group names",
			},
			"key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SSH Key name (for password: 0)",
			},
			"ssh_key": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Can use ssh key",
			},
			"number": {
				Type:        schema.TypeInt,
				Default:     1,
				Optional:    true,
				Description: "Count of abraks we want to create",
			},
			"create_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "type of abrak creation, " +
					"a flag which shows client sends arguments or we should read from image",
			},
			"disk_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Disk size of abraks we want to create",
			},
			"init_script": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Init script",
			},
			"ha_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Instance HA enabled or not",
			},
		},
	}
}

func resourceAbrakCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var errors diag.Diagnostics
	c := meta.(*client.Client).Iaas
	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	// Abrak Options
	abrak := &iaas.ServerOpts{
		Name:       data.Get("name").(string),
		FlavorId:   data.Get("flavor").(string),
		SshKey:     data.Get("ssh_key").(bool),
		Count:      data.Get("number").(int),
		CreateType: data.Get("create_type").(string),
		DiskSize:   data.Get("disk_size").(int),
		InitScript: data.Get("init_script").(string),
		HAEnabled:  data.Get("ha_enabled").(bool),
	}

	// KeyName
	abrak.KeyName = 0
	if v := data.Get("key_name"); v.(string) != "0" {
		abrak.KeyName = v.(string)
	}

	// ImageID
	if v, ok := data.GetOk("image"); ok && v.(*schema.Set).Len() > 0 {
		var iType, iName string
		for _, imageMap := range v.(*schema.Set).List() {
			imageInfo, ok := imageMap.(map[string]interface{})
			if !ok {
				errors = append(errors, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "invalid image details",
				})
				return errors
			}

			if v, ok := imageInfo["type"]; ok {
				iType = v.(string)
			}
			if v, ok := imageInfo["name"]; ok {
				iName = v.(string)
			}

			if iType != "" && iName != "" {
				imageId, err := c.Image.FindImageId(
					region,
					iName,
					iType,
				)
				if err != nil {
					errors = append(errors, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "image not found",
					})
					return errors
				}
				abrak.ImageId = *imageId
				break
			}
			errors = append(errors, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "image block in incomplete (name and type are required)",
			})
			return errors
		}
	} else {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "image block not found or invalid.",
		})
		return errors
	}

	// NetworkIDs
	if networks, ok := data.GetOk("networks"); !ok {
		opts, err := c.Server.Options(region)
		if err != nil {
			errors = append(errors, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "there is no default network",
			})
			return errors
		}
		abrak.NetworkIds = []string{opts.NetworkId}
	} else {
		var networkIds []string
		for _, network := range networks.([]interface{}) {
			networkId, err := c.Network.FindNetworkId(region, network.(string))
			if err != nil {
				errors = append(errors, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("there is no %v network", network.(string)),
				})
				return errors
			}
			networkIds = append(networkIds, *networkId)
		}
		abrak.NetworkIds = networkIds
	}

	// SecurityGroup
	if securityGroups, ok := data.GetOk("security_groups"); !ok {
		sg, err := c.SecurityGroup.FindSecurityGroupId(region, iaas.DefaultSecurityGroup)
		if err != nil {
			errors = append(errors, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "there is no default security group",
			})
			return errors
		}
		abrak.SecurityGroups = []iaas.ServerSecurityGroupOpts{{
			Name: *sg,
		}}
	} else {
		var sgs []iaas.ServerSecurityGroupOpts
		for _, sg := range securityGroups.([]interface{}) {
			securityGroupId, err := c.SecurityGroup.FindSecurityGroupId(region, sg.(string))
			if err != nil {
				errors = append(errors, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("there is no %v security group", sg.(string)),
				})
				return errors
			}
			sgs = append(sgs, iaas.ServerSecurityGroupOpts{
				Name: *securityGroupId,
			})
		}
		abrak.SecurityGroups = sgs
	}

	response, err := c.Server.Create(region, abrak)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(fmt.Sprint(response.ID))
	return errors
}

func resourceAbrakRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var errors diag.Diagnostics
	c := meta.(*client.Client).Iaas

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	_, err := c.Server.Read(region, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not retrive server %v", data.Id()),
		})
		return errors
	}

	// TODO: we have to store required items

	return nil
}

func resourceAbrakUpdate(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var errors diag.Diagnostics
	// TODO: we have to implement it
	return errors
}

func resourceAbrakDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var errors diag.Diagnostics
	c := meta.(*client.Client).Iaas
	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}
	err := c.Server.Delete(region, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not delete server %v", data.Id()),
		})
		return errors
	}
	return nil
}

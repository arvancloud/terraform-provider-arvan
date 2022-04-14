package helper

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DummyResourceAction(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

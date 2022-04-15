package iaas

import (
	"context"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

type Port struct {
	requester *api.Requester
}

// NewPort - generate communicator with Port
func NewPort(ctx context.Context) *Port {
	return &Port{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

// Enable - enable a port
func (p *Port) Enable(region, id string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/ports/%v/enable", ECCEndPoint, Version, region, id)
	_, err = p.requester.Patch(endpoint, nil, nil)
	return err
}

// Disable - disable a port
func (p *Port) Disable(region, id string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/ports/%v/disable", ECCEndPoint, Version, region, id)
	_, err = p.requester.Patch(endpoint, nil, nil)
	return err
}

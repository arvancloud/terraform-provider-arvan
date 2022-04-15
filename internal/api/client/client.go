package client

import (
	"context"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
)

type Config struct {
	ApiKey string
}

type Client struct {
	IaaS *iaas.IaaS
}

// NewClient - client for communicate with APIs
func NewClient(cfg *Config) (*Client, error) {
	var err error

	ctx := context.WithValue(
		context.Background(),
		api.RequesterContext,
		api.NewRequester(cfg.ApiKey),
	)

	// check user is authenticated or not
	err = ctx.Value(api.RequesterContext).(*api.Requester).CheckAuthenticate()
	if err != nil {
		return nil, err
	}

	return &Client{
		IaaS: iaas.NewIaaS(
			iaas.NewServer(ctx),
			iaas.NewImage(ctx),
			iaas.NewSizes(ctx),
			iaas.NewNetwork(ctx),
			iaas.NewSecurityGroup(ctx),
			iaas.NewVolume(ctx),
			iaas.NewFloatIP(ctx),
			iaas.NewRegion(ctx),
			iaas.NewSSHKey(ctx),
			iaas.NewTag(ctx),
			iaas.NewPort(ctx),
			iaas.NewPtr(ctx),
			iaas.NewQuota(ctx),
		),
	}, err
}

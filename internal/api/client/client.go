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
	Iaas *iaas.IaaS
}

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
		Iaas: iaas.NewIaaS(
			iaas.NewServer(ctx),
			iaas.NewImage(ctx),
			iaas.NewSizes(ctx),
			iaas.NewNetwork(ctx),
			iaas.NewSecurityGroup(ctx),
			iaas.NewVolume(ctx),
			iaas.NewFloatIP(ctx),
			iaas.NewRegion(ctx),
		),
	}, err
}

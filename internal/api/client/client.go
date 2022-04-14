package client

import (
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
)

type Client struct {
	Iaas *iaas.IaaS
}

func NewClient(apiKey string) *Client {
	requester := api.NewRequester(apiKey)
	return &Client{
		Iaas: iaas.NewIaaS(
			iaas.NewServer(requester),
			iaas.NewImage(requester),
			iaas.NewSizes(requester),
			iaas.NewNetwork(requester),
			iaas.NewSecurityGroup(requester),
			iaas.NewVolume(requester),
		),
	}
}

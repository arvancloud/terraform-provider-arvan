package iaas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

type PtrDetails struct {
}

type PtrOpts struct {
	IP     string `json:"ip"`
	Domain string `json:"domain"`
}

type Ptr struct {
	requester *api.Requester
}

// NewPtr - generate communicator with Ptr
func NewPtr(ctx context.Context) *Ptr {
	return &Ptr{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

// Create - create a ptr
func (p *Ptr) Create(region string, opts *PtrOpts) (*PtrDetails, error) {

	endpoint := fmt.Sprintf("/%v/%v/regions/%v/ptr", ECCEndPoint, Version, region)

	data, err := p.requester.Create(endpoint, opts, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details *PtrDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

// Delete - delete a ptr
func (p *Ptr) Delete(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/ptr/%v", ECCEndPoint, Version, region, id)
	return p.requester.Delete(endpoint, nil)
}

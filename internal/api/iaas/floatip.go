package iaas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

type FloatIPOpts struct {
	Description string `json:"description"`
}

type FloatIPDetails struct {
	ID                string         `json:"id"`
	Status            string         `json:"status"`
	FloatingNetworkId string         `json:"floating_network_id"`
	RouterId          string         `json:"router_id"`
	FixedIPAddress    string         `json:"fixed_ip_address"`
	FloatingIPAddress string         `json:"floating_ip_address"`
	PortId            string         `json:"port_id"`
	Description       string         `json:"description"`
	CreatedAt         string         `json:"created_at"`
	UpdatedAt         string         `json:"updated_at"`
	RevisionNumber    int            `json:"revision_number"`
	Server            *ServerDetails `json:"server"`
	Tags              []TagDetails   `json:"tags"`
}

type FloatIPAttachOpts struct {
	ServerId string `json:"server_id"`
	SubnetId string `json:"subnet_id"`
	PortId   string `json:"port_id"`
}

type FloatIPDetachOpts struct {
	PortId string `json:"port_id"`
}

type FloatIP struct {
	requester *api.Requester
}

// NewFloatIP - init communicator with FloatIP
func NewFloatIP(ctx context.Context) *FloatIP {
	return &FloatIP{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

// List - return all floatips
func (f *FloatIP) List(region string) (details []FloatIPDetails, err error) {

	endpoint := fmt.Sprintf("/%v/%v/regions/%v/float-ips", ECCEndPoint, Version, region)

	data, err := f.requester.List(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, &details)
	return details, err
}

// Create - create a floatip
func (f *FloatIP) Create(region string, opts *FloatIPOpts) (details *FloatIPDetails, err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/float-ips", ECCEndPoint, Version, region)

	data, err := f.requester.Create(endpoint, opts, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, &details)
	return details, err
}

// Delete - delete a floatip
func (f *FloatIP) Delete(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/fload-ips/%v", ECCEndPoint, Version, region, id)
	return f.requester.Delete(endpoint, nil)
}

// Attach - attach a floatip to a server
func (f *FloatIP) Attach(region, id string, opts *FloatIPAttachOpts) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/float-ip/%v/attach", ECCEndPoint, Version, region, id)
	_, err = f.requester.Patch(endpoint, opts, nil)
	return err
}

// Detach - attach a floatip from a server
func (f *FloatIP) Detach(region string, opts FloatIPDetachOpts) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/float-ip/detach", ECCEndPoint, Version, region)
	_, err = f.requester.Patch(endpoint, opts, nil)
	return err
}

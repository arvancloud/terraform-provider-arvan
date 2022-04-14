package iaas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

type FloatIPOpts struct {
	Description string `json:"description"`
}

type FloatIPDetails struct {
	ID                string
	Status            string
	FloatingNetworkId string
	RouterId          string
	FixedIPAddress    string
	FloatingIPAddress string
	PortId            string
	Description       string
	CreatedAt         string
	UpdatedAt         string
	RevisionNumber    int
	Server            *ServerDetails
	Tags              []Tag
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

func (f *FloatIP) List(region string) ([]FloatIPDetails, error) {

	endpoint := fmt.Sprintf("/%v/%v/regions/%v/float-ips", ECCEndPoint, Version, region)

	body, err := f.requester.DoRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var successResponse *api.SuccessResponse
	err = json.Unmarshal(body, &successResponse)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(successResponse.Data)
	if err != nil {
		return nil, err
	}

	var details []FloatIPDetails
	err = json.Unmarshal(data, &details)
	if err != nil {
		return nil, err
	}

	return details, nil
}

func (f *FloatIP) Create(region string, opts *FloatIPOpts) (*FloatIPDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/float-ips", ECCEndPoint, Version, region)

	body, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	response, err := f.requester.DoRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var successResponse *api.SuccessResponse
	err = json.Unmarshal(response, &successResponse)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(successResponse.Data)
	if err != nil {
		return nil, err
	}

	var details *FloatIPDetails
	err = json.Unmarshal(data, &details)
	if err != nil {
		return nil, err
	}

	return details, nil
}

func (f *FloatIP) Delete(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/fload-ips/%v", ECCEndPoint, Version, region, id)

	_, err := f.requester.DoRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}

func (f *FloatIP) Attach(region, id string, opts *FloatIPAttachOpts) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/float-ip/%v/attach", ECCEndPoint, Version, region, id)

	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	_, err = f.requester.DoRequest("PATCH", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return nil
}

func (f *FloatIP) Detach(region string, opts FloatIPDetachOpts) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/float-ip/detach", ECCEndPoint, Version, region)

	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	_, err = f.requester.DoRequest("PATCH", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return nil
}

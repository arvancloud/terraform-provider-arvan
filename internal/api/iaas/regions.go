package iaas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

type RegionDetails struct {
	Flag          string `json:"flag"`
	Country       string `json:"country"`
	CityCode      string `json:"city_code"`
	DcCode        string `json:"dc_code"`
	Dc            string `json:"dc"`
	Code          string `json:"code"`
	Region        string `json:"region"`
	Create        bool   `json:"create"`
	Soon          bool   `json:"soon"`
	Default       bool   `json:"default,omitempty"`
	VolumeBackend bool   `json:"volume_backend"`
	New           bool   `json:"new"`
	Beta          bool   `json:"beta"`
	Visible       bool   `json:"visible"`
}

type Region struct {
	requester *api.Requester
}

func NewRegion(ctx context.Context) *Region {
	return &Region{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

func (r *Region) List() ([]RegionDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/details", ECCEndPoint, Version)

	data, err := r.requester.List(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details []RegionDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

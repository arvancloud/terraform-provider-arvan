package iaas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

type PlanDetails struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	CpuCount         int     `json:"cpu_count"`
	Disk             int     `json:"disk"`
	DiskInBytes      int64   `json:"disk_in_bytes"`
	BandwidthInBytes int64   `json:"bandwidth_in_bytes"`
	Memory           int     `json:"memory"`
	MemoryInBytes    int64   `json:"memory_in_bytes"`
	PricePerHour     string  `json:"price_per_hour"`
	PricePerMonth    float64 `json:"price_per_month"`
	Generation       string  `json:"generation"`
	Type             string  `json:"type"`
	Subtype          string  `json:"subtype"`
	BasePackage      string  `json:"base_package"`
	CpuShare         string  `json:"cpu_share"`
	Order            string  `json:"order"`
	PPS              []int   `json:"pps"`
	IOpsMaxHDD       int     `json:"iops_max_hdd"`
	IOpsMaxSSD       int     `json:"iops_max_ssd"`
	Off              string  `json:"off"`
	OffPercent       string  `json:"off_percent"`
	Throughput       int64   `json:"throughput"`
	CreateType       string  `json:"create_type"`
	Canary           bool    `json:"canary"`
	Outbound         int64   `json:"outbound"`
}

type Sizes struct {
	requester *api.Requester
}

// NewSizes - init communicator with Sizes
func NewSizes(ctx context.Context) *Sizes {
	return &Sizes{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

// ListPlans - return all plans
func (s *Sizes) ListPlans(region string) ([]PlanDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/sizes", ECCEndPoint, Version, region)

	data, err := s.requester.List(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details []PlanDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

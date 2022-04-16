package iaas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
	"strings"
)

const (
	ImageTypeServer        = "server"
	ImageTypeSnapshot      = "snapshot"
	ImageTypeDistributions = "distributions"
)

type ImageServerDetails struct {
	ID              string       `json:"id"`
	Status          string       `json:"status"`
	Name            string       `json:"name"`
	MinRam          int          `json:"min_ram"`
	MinDisk         int          `json:"min_disk"`
	DiskFormat      string       `json:"disk_format"`
	Size            int64        `json:"size"`
	RealSize        int64        `json:"real_size"`
	Checksum        string       `json:"checksum"`
	CreatedAt       string       `json:"created_at"`
	ContainerFormat string       `json:"container_format"`
	Tags            []TagDetails `json:"tags"`
	Abrak           string       `json:"abrak"`
	AbrakId         string       `json:"abrak_id"`
	ImageType       string       `json:"image_type"`
}

type ImageDetails struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	DistributionName string `json:"distribution_name"`
	OSDescription    string `json:"os_description"`
	Disk             int    `json:"disk"`
	Ram              int    `json:"ram"`
	SSHKey           bool   `json:"ssh_key"`
	SSHPassword      bool   `json:"ssh_password"`
}

type ImageDistributionDetails struct {
	Name   string         `json:"name"`
	Images []ImageDetails `json:"images"`
}

type MarketPlaceDetail struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	DistributionName  string `json:"distrubation_name"`
	OSDescription     string `json:"OSDescription"`
	BaseOSDistro      string `json:"base_os_distro"`
	BaseOSVersion     string `json:"base_os_version"`
	Disk              int    `json:"disk"`
	Ram               int    `json:"ram"`
	SSHKey            bool   `json:"ssh_key"`
	SSHPassword       bool   `json:"ssh_password"`
	GenreId           string `json:"genre_id"`
	GenreName         string `json:"genre_name"`
	PartnerId         string `json:"partner_id"`
	PartnerName       string `json:"partner_name"`
	ImageId           string `json:"image_id"`
	ImageName         string `json:"image_name"`
	ImageVersion      string `json:"image_version"`
	ImageBuildNumber  string `json:"image_build_number"`
	ImageLatestUpdate string `json:"image_latest_update"`
	ImageSize         int    `json:"image_size"`
	IsIranian         bool   `json:"is_iranian"`
	Price             struct {
		PPH float64 `json:"pph"`
		PPM float64 `json:"ppm"`
	}
	ImageDocuments string `json:"image_documents"`
}

type Image struct {
	requester *api.Requester
}

// NewImage - init communicator with Image
func NewImage(ctx context.Context) *Image {
	return &Image{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

// Find - looking for an image by name
func (i *Image) Find(region, name, imageType string) (any, error) {
	images, err := i.List(region, imageType)
	if err != nil {
		return nil, err
	}

	imageName := strings.Split(name, "/")
	if len(imageName) < 1 {
		return nil, fmt.Errorf("invalid name")
	}

	switch images.(type) {
	case []ImageDistributionDetails:
		iGroup, iName := imageName[0], imageName[1]
		for _, group := range images.([]ImageDistributionDetails) {
			for _, image := range group.Images {
				if strings.ToLower(image.Name) == iName && strings.ToLower(group.Name) == iGroup {
					return image, nil
				}
			}
		}
		return nil, fmt.Errorf("not found image %v", name)
	case []ImageServerDetails:
		iName := imageName[0]
		for _, image := range images.([]ImageServerDetails) {
			if iName == strings.ToLower(image.Name) {
				return image, nil
			}
		}
		return nil, fmt.Errorf("not found image %v", name)
	default:
		return nil, fmt.Errorf("invalid image's type")
	}
}

// List - return all images
func (i *Image) List(region, imageType string) (any, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/images", ECCEndPoint, Version, region)

	data, err := i.requester.List(endpoint, map[string]string{
		"type": imageType,
	})
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if imageType == ImageTypeDistributions {
		var details []ImageDistributionDetails
		err = json.Unmarshal(marshal, &details)
		return details, err
	}

	var details []ImageServerDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

// ListMarketPlace - return all images at marketplace
func (i *Image) ListMarketPlace(region string) (details []MarketPlaceDetail, err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/images/marketplace", ECCEndPoint, Version, region)

	data, err := i.requester.List(endpoint, nil)
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

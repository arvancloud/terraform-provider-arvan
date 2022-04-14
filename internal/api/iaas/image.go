package iaas

import (
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
	ID              string `json:"id"`
	Status          string `json:"status"`
	Name            string `json:"name"`
	MinRam          int    `json:"min_ram"`
	MinDisk         int    `json:"min_disk"`
	DiskFormat      string `json:"disk_format"`
	Size            int64  `json:"size"`
	RealSize        int64  `json:"real_size"`
	Checksum        string `json:"checksum"`
	CreatedAt       string `json:"created_at"`
	ContainerFormat string `json:"container_format"`
	Tags            []Tag  `json:"tags"`
	Abrak           string `json:"abrak"`
	AbrakId         string `json:"abrak_id"`
	ImageType       string `json:"image_type"`
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

func NewImage(requester *api.Requester) *Image {
	return &Image{
		requester: requester,
	}
}

func (i *Image) List(region, imageType string) (interface{}, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/images", ECCEndPoint, Version, region)

	data, err := i.requester.DoRequestWithQuery("GET", endpoint, nil, map[string]string{
		"type": imageType,
	})
	if err != nil {
		return nil, err
	}

	var response *api.SuccessResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	if imageType == ImageTypeDistributions {
		var imageDetails []ImageDistributionDetails
		err = json.Unmarshal(dataBytes, &imageDetails)
		if err != nil {
			return nil, err
		}
		return imageDetails, nil
	}

	var imageDetails []ImageServerDetails
	err = json.Unmarshal(dataBytes, &imageDetails)
	if err != nil {
		return nil, err
	}
	return imageDetails, nil
}

func (i *Image) FindImageId(region, name, imageType string) (*string, error) {
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
					return &image.ID, nil
				}
			}
		}
		return nil, fmt.Errorf("not found image %v", name)
	case []ImageServerDetails:
		iName := imageName[0]
		for _, image := range images.([]ImageServerDetails) {
			if iName == strings.ToLower(image.Name) {
				return &image.ID, nil
			}
		}
		return nil, fmt.Errorf("not found image %v", name)
	default:
		return nil, fmt.Errorf("invalid image's type")
	}
}

func (i *Image) ListMarketPlace(region string) ([]MarketPlaceDetail, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/images/marketplace", ECCEndPoint, Version, region)

	data, err := i.requester.DoRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response *api.SuccessResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var marketPlaces []MarketPlaceDetail
	err = json.Unmarshal(dataBytes, &marketPlaces)
	if err != nil {
		return nil, err
	}

	return marketPlaces, nil
}

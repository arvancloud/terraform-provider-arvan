package iaas

const (
	ECCEndPoint    = "/ecc/"
	Version        = "v1"
	ForoghRegion   = "ir-thr-c2"
	ShahriarRegion = "ir-tbz-dc1"
	HermanRegion   = "nl-ams-su1"
	AsiaTechRegion = "ir-thr-at1"
)

var (
	AvailableRegions = []string{
		ForoghRegion,
		ShahriarRegion,
		HermanRegion,
		AsiaTechRegion,
	}
)

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type IaaS struct {
	Network       *Network
	Server        *Server
	Image         *Image
	Sizes         *Sizes
	SecurityGroup *SecurityGroup
	Volume        *Volume
}

func NewIaaS(server *Server, image *Image,
	sizes *Sizes, network *Network, securityGroup *SecurityGroup,
	volume *Volume) *IaaS {
	return &IaaS{
		Network:       network,
		Server:        server,
		Image:         image,
		Sizes:         sizes,
		SecurityGroup: securityGroup,
		Volume:        volume,
	}
}

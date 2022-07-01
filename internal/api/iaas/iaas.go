package iaas

const (
	ECCEndPoint    = "/ecc/"
	Version        = "v1"
	SiminRegion    = "ir-thr-c1"
	ForoghRegion   = "ir-thr-c2"
	ShahriarRegion = "ir-tbz-dc1"
	HermanRegion   = "nl-ams-su1"
	BamdadRegion   = "ir-thr-w1"
)

var (
	AvailableRegions = []string{
		SiminRegion,
		ForoghRegion,
		ShahriarRegion,
		HermanRegion,
		BamdadRegion,
	}
)

type IaaS struct {
	FloatIP       *FloatIP
	Image         *Image
	Network       *Network
	Port          *Port
	Quota         *Quota
	Ptr           *Ptr
	Region        *Region
	SSHKey        *SSHKey
	SecurityGroup *SecurityGroup
	Server        *Server
	Sizes         *Sizes
	Tag           *Tag
	Volume        *Volume
}

func NewIaaS(server *Server, image *Image,
	sizes *Sizes, network *Network,
	securityGroup *SecurityGroup, volume *Volume,
	floatIP *FloatIP, region *Region, sshKey *SSHKey,
	tag *Tag, port *Port, ptr *Ptr, quota *Quota) *IaaS {
	return &IaaS{
		FloatIP:       floatIP,
		Image:         image,
		Network:       network,
		Port:          port,
		Quota:         quota,
		Ptr:           ptr,
		Region:        region,
		SSHKey:        sshKey,
		SecurityGroup: securityGroup,
		Server:        server,
		Sizes:         sizes,
		Tag:           tag,
		Volume:        volume,
	}
}

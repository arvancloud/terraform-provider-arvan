package main

import (
	"flag"
	"github.com/arvancloud/terraform-provider-arvan/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"log"
)

func init() {
	logFlags := log.Flags()
	logFlags = logFlags &^ (log.Ldate | log.Ltime)
	log.SetFlags(logFlags)
}

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false,
		"set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
		Debug:        debugMode,
	})

	//ApiKey := "Apikey 09fc3ecf-73c5-5aff-9aaa-14a098b93c9b"
	//
	//c, err := client.NewClient(&client.Config{ApiKey: ApiKey})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//opts := &iaas.NetworkAttachOpts{
	//	ServerId: "31c25a1d-680b-48dd-8dd6-ed79fa56a534",
	//}
	//data, err := c.IaaS.Network.Attach(iaas.ForoghRegion, "2f42d4de-3039-49f8-a76b-f93d7a7627c8", opts)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//fmt.Println(data)

	//i, err := c.IaaS.Image.Find(iaas.ForoghRegion, "debian/11", iaas.ImageTypeDistributions)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//fmt.Println(i.(iaas.ImageDetails).ID)

	//opts := &iaas.ServerOpts{
	//	FlavorId: "g1-1-1-0",
	//	Name:     "asdasd2",
	//	ImageId:  "767ee24e-118c-4447-a04e-a96e82ababf7",
	//	NetworkIds: []string{
	//		"ffc8dbf2-bdc1-4d9a-b64c-d0951d70d6a6",
	//	},
	//	SecurityGroups: []iaas.ServerSecurityGroupOpts{
	//		{
	//			Name: "b582e17c-a79b-4785-be27-e18f507f1c6c",
	//		},
	//	},
	//	DiskSize: 25,
	//	KeyName:  0,
	//}
	//
	//server, err := c.IaaS.Server.Create(iaas.ForoghRegion, opts)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//fmt.Println(server.ID)
}

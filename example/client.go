package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/paysuper/geoip-service/pkg"
	"github.com/paysuper/geoip-service/pkg/proto"
)

func main() {
	// create a new service
	service := micro.NewService()

	// parse command line flags
	service.Init()

	// Create new greeter client
	client := proto.NewGeoIpService(geoip.ServiceName, service.Client())

	// Call it
	rsp, err := client.GetIpData(context.TODO(), &proto.GeoIpDataRequest{IP: "136.0.16.217"})
	if err != nil {
		fmt.Println(err)
	}

	// Print response
	fmt.Printf("%+v\n", rsp)
}

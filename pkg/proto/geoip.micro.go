// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: geoip.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	geoip.proto

It has these top-level messages:
	GeoIpDataRequest
	GeoIpDataResponse
	GeoIpCity
	GeoIpContinent
	GeoIpCountry
	GeoIpLocation
	GeoIpPostal
	GeoIpRepresentedCountry
	GeoIpSubdivision
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for GeoIpService service

type GeoIpService interface {
	GetIpData(ctx context.Context, in *GeoIpDataRequest, opts ...client.CallOption) (*GeoIpDataResponse, error)
}

type geoIpService struct {
	c    client.Client
	name string
}

func NewGeoIpService(name string, c client.Client) GeoIpService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "proto"
	}
	return &geoIpService{
		c:    c,
		name: name,
	}
}

func (c *geoIpService) GetIpData(ctx context.Context, in *GeoIpDataRequest, opts ...client.CallOption) (*GeoIpDataResponse, error) {
	req := c.c.NewRequest(c.name, "GeoIpService.GetIpData", in)
	out := new(GeoIpDataResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GeoIpService service

type GeoIpServiceHandler interface {
	GetIpData(context.Context, *GeoIpDataRequest, *GeoIpDataResponse) error
}

func RegisterGeoIpServiceHandler(s server.Server, hdlr GeoIpServiceHandler, opts ...server.HandlerOption) error {
	type geoIpService interface {
		GetIpData(ctx context.Context, in *GeoIpDataRequest, out *GeoIpDataResponse) error
	}
	type GeoIpService struct {
		geoIpService
	}
	h := &geoIpServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&GeoIpService{h}, opts...))
}

type geoIpServiceHandler struct {
	GeoIpServiceHandler
}

func (h *geoIpServiceHandler) GetIpData(ctx context.Context, in *GeoIpDataRequest, out *GeoIpDataResponse) error {
	return h.GeoIpServiceHandler.GetIpData(ctx, in, out)
}

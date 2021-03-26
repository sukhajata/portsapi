package api

import (
	"context"
	"github.com/sukhajata/portsapi/service/internal/core"
	"github.com/sukhajata/portsapi/service/pkg/proto"
)

type GRPCServer struct {
	coreService *core.Service
	*proto.UnimplementedPortDomainServiceServer
}

func NewGRPCServer(coreService *core.Service) *GRPCServer {
	return &GRPCServer{
		coreService: coreService,
	}
}

// GetPorts - returns a list of all ports
func (s *GRPCServer) GetPorts(ctx context.Context, req *proto.GetPortsRequest) (*proto.Ports, error) {
	return s.coreService.GetPorts(req)
}

// GetPort - retreive a port by id
func (s *GRPCServer) GetPort(ctx context.Context, req *proto.GetPortRequest) (*proto.Port, error) {
	return s.coreService.GetPort(req)
}

// UpsertPort - insert or update a port
func (s *GRPCServer) UpsertPort(ctx context.Context, req *proto.UpsertPortRequest) (*proto.Port, error) {
	return s.coreService.UpsertPort(req)
}

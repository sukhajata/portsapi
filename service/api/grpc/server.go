package grpc

type GRPCServer struct {
	*UnimplementedPortDomainServiceServer
}

// GetPorts - returns a list of all ports
func (s *GRPCServer) GetPorts(GetPortsRequest)  (Ports) {

}

// GetPort - retreive a port by id
func (s *GRPCServer) GetPort(GetPortRequest) (Port) {

}

// UpsertPort - insert or update a port
func (s *GRPCServer) UpsertPort(UpsertPortRequest) (Port) {

}
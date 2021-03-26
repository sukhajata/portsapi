package core

import "github.com/sukhajata/portsapi/service/pkg/db"

type Service struct {
	dbEngine db.SQlEngine
}

func NewService() {

}

// GetPorts - returns a list of all ports
func (s *Service) GetPorts(GetPortsRequest)  (Ports) {

}

// GetPort - retreive a port by id
func (s *Service) GetPort(GetPortRequest) (Port) {

}

// UpsertPort - insert or update a port
func (s *Service) UpsertPort(UpsertPortRequest) (Port) {

}

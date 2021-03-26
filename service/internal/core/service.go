package core

import (
	"encoding/json"
	"fmt"
	"github.com/sukhajata/portsapi/service/pkg/db"
	"github.com/sukhajata/portsapi/service/pkg/proto"
)

type Service struct {
	dbEngine db.SQLEngine
}

func NewService(dbEngine db.SQLEngine) *Service {
	return &Service{
		dbEngine: dbEngine,
	}
}

// GetPorts - returns a list of all ports
func (s *Service) GetPorts(req *proto.GetPortsRequest) (*proto.Ports, error) {
	sql := "SELECT data FROM ports"
	rows, err := s.dbEngine.QueryTextColumn(sql)
	if err != nil {
		return &proto.Ports{}, err
	}

	var ports []*proto.Port
	for _, v := range rows {
		var port *proto.Port
		err := json.Unmarshal([]byte(v), &port)
		if err != nil {
			return &proto.Ports{}, err
		}
		ports = append(ports, port)
	}

	return &proto.Ports{
		Ports: ports,
	}, nil
}

// GetPort - retreive a port by id
func (s *Service) GetPort(req *proto.GetPortRequest) (*proto.Port, error) {
	sql := "SELECT data FROM ports WHERE id = $1"
	rows, err := s.dbEngine.QueryTextColumn(sql, req.Id)
	if err != nil {
		return &proto.Port{}, err
	}

	if len(rows) == 0 {
		return &proto.Port{}, fmt.Errorf("No port found with id %s", req.Id)
	}

	var port *proto.Port
	err = json.Unmarshal([]byte(rows[0]), &port)
	if err != nil {
		return &proto.Port{}, err
	}

	return port, nil
}

// UpsertPort - insert or update a port
func (s *Service) UpsertPort(req *proto.UpsertPortRequest) (*proto.Port, error) {
	sql := "INSERT INTO ports (id, data) VALUES($1, $2)"
	err := s.dbEngine.Exec(sql, req.Id, req.Port)
	if err != nil {
		return &proto.Port{}, err
	}

	return s.GetPort(&proto.GetPortRequest{Id: req.Id})

}

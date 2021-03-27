package cmd

import (
	"fmt"
	"github.com/sukhajata/portsapi/service/api"
	"github.com/sukhajata/portsapi/service/internal/core"
	"github.com/sukhajata/portsapi/service/pkg/db"
	"github.com/sukhajata/portsapi/service/pkg/proto"
	"google.golang.org/grpc"
	"net"
	"os"
)

var (
	grpcPort = os.Getenv("grpcPort")
	psqlURL  = os.Getenv("psqlURL")
)

func main() {
	// postgres
	dbEngine, err := db.NewPostgresEngine(psqlURL)
	if err != nil {
		panic(err)
	}

	// core service
	coreService := core.NewService(dbEngine)

	// grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	portsDomainServerImpl := api.NewGRPCServer(coreService)
	grpcServer := grpc.NewServer()
	proto.RegisterPortDomainServiceServer(grpcServer, portsDomainServerImpl)
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}

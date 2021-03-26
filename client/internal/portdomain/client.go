package portdomain

import (
	"context"
	"fmt"
	"github.com/sukhajata/portsapi/client/internal/processor"
	"google.golang.org/grpc"
	"sync"
	"time"
	pb "github.com/sukhajata/portsapi/service/pkg/proto"
)

type PortDomainClient struct {
	grpcClient pb.PortDomainServiceClient
}

func NewPortDomainClient(address string) (*PortDomainClient, *grpc.ClientConn, error){
	var conn *grpc.ClientConn
	var err error

	for i := 0; i < 6; i++ {
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 1)
			continue
		}

		break
	}

	if err != nil {
		return nil, nil, err
	}

	portDomainClient := pb.NewPortDomainServiceClient(conn)
	return &PortDomainClient{
		grpcClient: portDomainClient,
	}, conn, nil
}

func (p *PortDomainClient) SendPorts(inChan <-chan *processor.PortItem, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range inChan {
		_, err := p.UpsertPort(&pb.UpsertPortRequest{
			Id:   item.Id,
			Port: &item.Port,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (p *PortDomainClient) UpsertPort(req *pb.UpsertPortRequest) (*pb.Port, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.grpcClient.UpsertPort(ctx, req)
}

func (p *PortDomainClient) GetPort(id string) (*pb.Port, error) {
	req := &pb.GetPortRequest{
		Id: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.grpcClient.GetPort(ctx, req)
}
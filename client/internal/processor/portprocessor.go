package processor

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/sukhajata/portsapi/service/pkg/proto"
	"github.com/sukhajata/portsapi/client/pkg/parser"
)

type PortItem struct {
	Id string
	Port pb.Port
}

func ProcessPorts(ctx context.Context, inChan <-chan *parser.DictionaryItem) <-chan *PortItem {
	outChan := make(chan *PortItem, 2)

	select {
	case <-ctx.Done():
		close(outChan)
	case item := <-inChan:
		var port pb.Port
		err := json.Unmarshal([]byte(item.Value), &port)
		if err != nil {
			fmt.Println(err)
		} else {
			id := item.Key
			outChan <- &PortItem{
				Id: id,
				Port: port,
			}
		}
	}

	return outChan
}

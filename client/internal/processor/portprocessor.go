package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sukhajata/portsapi/client/pkg/parser"
	pb "github.com/sukhajata/portsapi/service/pkg/proto"
)

type PortItem struct {
	Id string
	Port pb.Port
}

func ProcessPorts(ctx context.Context, inChan <-chan *parser.DictionaryItem) <-chan *PortItem {
	outChan := make(chan *PortItem, 2)

	go func() {
		for item := range inChan {
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
		close(outChan)
		// TODO - use select with cancel
		/*select {
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
		}*/
	}()

	return outChan
}

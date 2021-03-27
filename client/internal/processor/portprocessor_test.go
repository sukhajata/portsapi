package processor_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/sukhajata/portsapi/client/internal/processor"
	"github.com/sukhajata/portsapi/client/pkg/parser"
	"testing"
)

func TestProcessPorts(t *testing.T) {
	inChan := make(chan *parser.DictionaryItem, 1)
	outChan := processor.ProcessPorts(context.Background(), inChan)

	go func() {
		for i := 0; i < 5; i++ {
			inChan <- &parser.DictionaryItem{
				Key: "TEST",
				Value: `{
				"name": "Abu Dhabi",
				"coordinates": [
				  54.37,
				  24.47
				],
				"city": "Abu Dhabi",
				"province": "Abu Z¸aby [Abu Dhabi]",
				"country": "United Arab Emirates",
				"alias": ["Test"],
				"regions": ["Asia", "Middle East"],
				"timezone": "Asia/Dubai",
				"unlocs": [
				  "AEAUH"
				],
				"code": "52001"
				}`,
			}
		}
		close(inChan)
	}()

	// assert
	count := 0
	for item := range outChan {
		require.Equal(t, "Abu Dhabi", item.Port.City)
		require.Equal(t, "Asia/Dubai", item.Port.Timezone)
		require.Equal(t, 2, len(item.Port.Regions))
		count++
	}

	require.Equal(t, 5, count)
}
package parser_test

import (
	"bufio"
	"bytes"
	"context"
	"github.com/stretchr/testify/require"
	"github.com/sukhajata/portsapi/client/pkg/parser"
	"strings"
	"testing"
)

func generateContent(count int) string {
	content := `  
		"AEAUH": {
			"name": "Abu Dhabi",
			"coordinates": [
			  54.37,
			  24.47
			],
			"city": "Abu Dhabi",
			"province": "Abu Z¸aby [Abu Dhabi]",
			"country": "United Arab Emirates",
			"alias": [],
			"regions": [],
			"timezone": "Asia/Dubai",
			"unlocs": [
			  "AEAUH"
			],
			"code": "52001"
		},`

	var buffer bytes.Buffer

	for i := 0; i < count; i++ {
		buffer.WriteString(content)
	}

	return "{" + buffer.String() + "}"

}

// TestDictionaryParser_Read_Lite - try to process 10 ports
func TestDictionaryParser_Read_Lite(t *testing.T) {
	// arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b := strings.NewReader(generateContent(10))
	r := bufio.NewReader(b)
	dictionaryParser := parser.NewDictionaryParser(r)

	// execute
	portDictionaryChan := dictionaryParser.Read(ctx)

	count := 0
	for range portDictionaryChan {
		count++
	}

	// assert
	require.Equal(t, 10, count)
}

// TestDictionaryParser_Read_Medium - try to process 10,000 ports
func TestDictionaryParser_Read_Medium(t *testing.T) {
	// arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b := strings.NewReader(generateContent(10000))
	r := bufio.NewReader(b)
	dictionaryParser := parser.NewDictionaryParser(r)

	// execute
	portDictionaryChan := dictionaryParser.Read(ctx)

	count := 0
	for range portDictionaryChan {
		count++
	}

	// assert
	require.Equal(t, 10000, count)
}

// TestDictionaryParser_Read_Large - try to process 3,000,000 ports
func TestDictionaryParser_Read_Large(t *testing.T) {
	// arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b := strings.NewReader(generateContent(3000000))
	r := bufio.NewReader(b)
	dictionaryParser := parser.NewDictionaryParser(r)

	// execute
	portDictionaryChan := dictionaryParser.Read(ctx)

	count := 0
	for range portDictionaryChan {
		count++
	}

	// assert
	require.Equal(t, 3000000, count)
}

// TestDictionaryParser_Read_Corrupt - try to read corrupt json
// channel should close after first read
func TestDictionaryParser_Read_Corrupt(t *testing.T) {
	// arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	content := `
	{
		"AEKM": 1,
		{ 1, 2, 3}
	}`
	b := strings.NewReader(content)
	r := bufio.NewReader(b)
	dictionaryParser := parser.NewDictionaryParser(r)

	// execute
	portDictionaryChan := dictionaryParser.Read(ctx)

	count := 0
	for range portDictionaryChan {
		count++
	}

	// assert
	require.Equal(t, 1, count)
}

func BenchmarkDictionaryParser_Read(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sr := strings.NewReader(generateContent(10))
	r := bufio.NewReader(sr)
	dictionaryParser := parser.NewDictionaryParser(r)

	for n := 0; n < b.N; n++ {
		dictionaryParser.Read(ctx)
	}

	/*var channels []<-chan *parser.DictionaryItem
	for n := 0; n < b.N; n++ {
		outChan := dictionaryParser.Read(ctx)
		channels = append(channels, outChan)
	}

	count := 0
	for _, c := range channels {
		for range c {
			count++
		}
	}*/

}
package parser_test

import (
	"bufio"
	"context"
	"github.com/stretchr/testify/require"
	"github.com/sukhajata/portsapi/client/pkg/parser"
	"os"
	"testing"
)

func TestDictionaryParser_Read(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f, err := os.Open("ports.json")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	dictionaryParser := parser.NewDictionaryParser(r)
	portDictionaryChan := dictionaryParser.Read(ctx)

	count := 0
	for range portDictionaryChan {
		count++
	}

	require.Equal(t, 100, count)
}

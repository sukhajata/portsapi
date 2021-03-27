package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/sukhajata/portsapi/client/api"
	"github.com/sukhajata/portsapi/client/internal/portdomain"
	"github.com/sukhajata/portsapi/client/internal/processor"
	"github.com/sukhajata/portsapi/client/pkg/parser"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	portDomainServiceAddress = os.Getenv("portDomainServiceAddress")
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// port domain service client
	portDomainService, conn, err := portdomain.NewPortDomainClient(portDomainServiceAddress)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	// parse json file into dictionary
	f, err := os.Open("ports.json")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	dictionaryParser := parser.NewDictionaryParser(r)
	portDictionaryChan := dictionaryParser.Read(ctx)

	// unmarshal dictionary values into struct
	portItemChan := processor.ProcessPorts(ctx, portDictionaryChan)

	// send ports to port domain service
	// use multiple go routines
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go portDomainService.SendPorts(portItemChan, &wg)
	}

	// http api
	httpClient := api.NewHTTPServer(portDomainService)

	// handle kill signal
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.

	// graceful shutdown
	httpClient.Ready = false
	cancel()
	wg.Wait()
}
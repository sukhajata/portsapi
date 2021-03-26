package cmd

import (
	"bufio"
	"github.com/sukhajata/portsapi/client/internal/parser"
	pb "github.com/sukhajata/portsapi/service/proto"
	"os"
	"os/signal"
	"sync"
	"syscall"
)


func main() {
	var wg sync.WaitGroup

	f, err := os.Open("ports.json")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	parser := parser.NewParser(r)

	// handle kill signal
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.

	// wait for go routines to finish
	wg.Wait()
}
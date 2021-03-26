package cmd

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup



	// handle kill signal
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.

	// wait for go routines to finish
	wg.Wait()
}
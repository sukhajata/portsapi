package internal

import (
	"os"
	"sync"

)

type Parser struct {}

func (p *Parser) ProcessChunk(f *os.File, buf []byte, wg *sync.WaitGroup) error {
	defer wg.Done()

	return nil
}
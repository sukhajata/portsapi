package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(reader *bufio.Reader) *Parser {
	return &Parser{reader: reader}
}

func (p *Parser) ProcessChunk(f *os.File, buf []byte, wg *sync.WaitGroup) error {
	defer wg.Done()

	return nil
}

func (p *Parser) Read() {
	// first line
	_, err := p.reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	for {
		err = p.ReadPort()
		if err != nil {
			break
		}
		_, err = p.reader.ReadString(',')
		if err != nil {
			break
		}
	}
}

func (p *Parser) ReadPort() error {
	// id
	buf, err := p.reader.ReadString(':')
	if err != nil {
		return err
	}
	id := strings.Trim(buf, ":")
	fmt.Println(id)

	// item
	output, err := p.reader.ReadString('}')
	if err != nil {
		return err
	}

	var port proto.Port
	err = json.Unmarshal([]byte(output), &port)
	if err != nil {
		return err
	}
	fmt.Println(port)

	return nil
}
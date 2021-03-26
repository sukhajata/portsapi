package parser

import (
	"bufio"
	"context"
	"strings"
)

type DictionaryParser struct {
	reader *bufio.Reader
}

type DictionaryItem struct {
	Key string
	Value string
}

func NewDictionaryParser(reader *bufio.Reader) *DictionaryParser {
	return &DictionaryParser{reader: reader}
}

func (p *DictionaryParser) Read(ctx context.Context) <-chan *DictionaryItem {
	outChan := make(chan *DictionaryItem, 2)

	// first line
	_, err := p.reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	for {
		err = p.ReadItem(outChan)
		if err != nil {
			break
		}
		_, err = p.reader.ReadString(',')
		if err != nil {
			break
		}
	}

	return outChan
}

func (p *DictionaryParser) ReadItem(outChan chan<- *DictionaryItem) error {
	// id
	buf, err := p.reader.ReadString(':')
	if err != nil {
		return err
	}
	id := strings.Trim(buf, ":")

	// item
	output, err := p.reader.ReadString('}')
	if err != nil {
		return err
	}

	outChan <- &DictionaryItem{
		Key:   id,
		Value: output,
	}
	return nil
}
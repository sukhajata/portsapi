package parser

import (
	"bufio"
	"context"
	"strings"
	"sync"
)

type DictionaryParser struct {
	reader *bufio.Reader
	stringPool sync.Pool
}

type DictionaryItem struct {
	Key string
	Value string
}

func NewDictionaryParser(reader *bufio.Reader) *DictionaryParser {
	// sync pool for recycling string variables
	stringPool := sync.Pool{
		New: func() interface{} {
			value := ""
			return value
		},
	}

	return &DictionaryParser{
		reader: reader,
		stringPool: stringPool,
	}
}

func (p *DictionaryParser) Read(ctx context.Context) <-chan *DictionaryItem {
	outChan := make(chan *DictionaryItem, 2)
	//quitChan := make(chan int)
	go func() {
		// first line
		_, err := p.reader.ReadString('\n')
		if err != nil {
			close(outChan)
			return
		}

		for {
			err = p.ReadItem(outChan)
			if err != nil {
				close(outChan)
				break
			}
			_, err = p.reader.ReadString(',')
			if err != nil {
				close(outChan)
				break
			}

			// TODO - fix select with cancel
			/*select {
			case <-ctx.Done():
				close(outChan)
				return
			case <-quitChan:
				close(outChan)
				return
			default:
				err = p.ReadItem(outChan)
				if err != nil {
					quitChan <- 1
				}
				// move to next item
				_, err = p.reader.ReadString(',')
				if err != nil {
					quitChan <- 1
				}
			}*/

		}
	}()

	return outChan
}

func (p *DictionaryParser) ReadItem(outChan chan<- *DictionaryItem) error {
	// read id
	buf, err := p.reader.ReadString(':')
	if err != nil {
		return err
	}
	id := strings.Trim(buf, ":")

	// read item
	output := p.stringPool.Get().(string)
	output, err = p.reader.ReadString('}')
	if err != nil {
		return err
	}

	outChan <- &DictionaryItem{
		Key:   id,
		Value: output,
	}
	return nil
}
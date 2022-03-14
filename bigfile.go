package suckutils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
)

type bigFile struct {
	file         *os.File
	block_size   int
	define_block func([]byte) (block []byte, remain []byte, err error)
}

type Block struct {
	Bytes []byte
	Err   error
}

func Read(ctx context.Context, file *os.File, block_size int, define_block func([]byte) ([]byte, []byte, error)) <-chan []byte {
	f := &bigFile{
		file:         file,
		block_size:   block_size,
		define_block: define_block,
	}
	result := make(chan []byte, 10)
	var remain []byte
	var err error
	go func() {
		defer close(result)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				remain, err = f.readNextBlock(remain, result)
				fmt.Println(len(remain), err)
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					//TODO
				}
			}
		}
	}()
	return result
}

func (f *bigFile) readNextBlock(remain []byte, ch chan<- []byte) ([]byte, error) {
	buf := make([]byte, f.block_size+len(remain))
	n, err := f.file.Read(buf[len(remain):])
	if err != nil {
		return nil, err
	}
	if len(remain) > 0 {
		copy(buf, remain)
	}
	if n < f.block_size {
		buf = buf[:n+len(remain)]
	}
	var block []byte
	for {
		block, buf, err = f.define_block(buf)
		if err != nil {
			return buf, err
		}
		if block == nil {
			return buf, nil
		}
		ch <- block
	}
}

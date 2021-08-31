package suckutils

import (
	"context"
	"errors"
	"os"
	"time"
)

type ConcurrentFile struct {
	File         *os.File
	lockFilename string
}

func OpenConcurrentFile(ctx context.Context, filename string, timeout time.Duration) (*ConcurrentFile, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	timer := time.NewTicker(time.Millisecond * 100)
	lockFilename := ConcatTwo(filename, ".lock")
	var err error
	var result *ConcurrentFile
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			if result == nil && err == nil {
				err = errors.New("Timeout")
			}
			cancel()
			return result, err
		case <-timer.C:
			f, err := os.OpenFile(lockFilename, os.O_CREATE|os.O_EXCL, 0664)
			if e, ok := err.(*os.PathError); ok && errors.Is(e.Err, os.ErrExist) {
				break
			} else if err != nil {
				cancel()
				break
			}
			f.Close()
			f, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0664)
			if err == nil {
				result = &ConcurrentFile{File: f, lockFilename: lockFilename}
			}
			cancel()
		}
	}
}

func (file *ConcurrentFile) Close() error {
	err := file.File.Close()
	if err != nil {
		return err
	}
	return os.Remove(file.lockFilename)
}

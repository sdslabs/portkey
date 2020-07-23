package transfer

import (
	"io"
	"io/ioutil"
	"os"
	"sync"

	"github.com/pion/quic"
)

const receiveBufferSize = 100

func ReadLoop(stream *quic.BidirectionalStream, receivePath string, receiveErr chan error, wg *sync.WaitGroup) error {
	defer wg.Done()
	zipfile, err := ioutil.TempFile(os.TempDir(), "portkey*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(zipfile.Name())

	buffer := make([]byte, receiveBufferSize)

	for {
		params, err := stream.ReadInto(buffer)
		if err != nil {
			if err != io.EOF {
				return err
			}
		}
		_, err = zipfile.Write(buffer[:params.Amount])
		if err != nil {
			return err
		}
		if params.Finished {
			break
		}
	}

	if receivePath == "" {
		receivePath, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	err = unzip(zipfile, receivePath)
	return err
}

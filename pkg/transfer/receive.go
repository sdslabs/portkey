package transfer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"

	"github.com/pion/quic"
)

const receiveBufferSize = 100

func ReadLoop(stream *quic.BidirectionalStream, receivePath string, receiveErr chan error, wg *sync.WaitGroup) error {
	defer wg.Done()
	fmt.Println("inside readloop")
	zipfile, err := ioutil.TempFile(os.TempDir(), "portkey*.zip")
	if err != nil {
		fmt.Println("error1", err)
		return err
	}
	defer os.Remove(zipfile.Name())

	buffer := make([]byte, receiveBufferSize)

	for {
		params, err := stream.ReadInto(buffer)
		fmt.Println("receiving", buffer, params)
		if err != nil {
			if err != io.EOF {
				fmt.Println("error2", err)
				return err
			}
		}
		_, err = zipfile.Write(buffer[:params.Amount])
		if err != nil {
			fmt.Println("error3", err)
			return err
		}
		if params.Finished {
			break
		}
	}

	if receivePath == "" {
		receivePath, err = os.Getwd()
		if err != nil {
			fmt.Println("error4", err)
			return err
		}
	}

	err = unzip(zipfile, receivePath)
	if err == nil {
		fmt.Printf("Finished reading from Stream %d\n", stream.StreamID())
	}
	fmt.Println("error5", err)
	return err
}

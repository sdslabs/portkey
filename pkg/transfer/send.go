package transfer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"

	"github.com/pion/quic"
)

const sendBufferSize = 100

func WriteLoop(stream *quic.BidirectionalStream, sendPath string, sendErr chan error, wg *sync.WaitGroup) error {
	defer wg.Done()
	fmt.Println("inside write loop")
	zipfile, err := ioutil.TempFile(os.TempDir(), "portkey*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(zipfile.Name())
	err = zipit(sendPath, zipfile)
	if err != nil {
		fmt.Println("error 1", err)
		return err
	}

	if _, err = zipfile.Seek(0, 0); err != nil {
		fmt.Println("error 2", err)
		return err
	}
	finished := false
	buffer := make([]byte, sendBufferSize)
	for {
		n, err := zipfile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				finished = true
			} else {
				fmt.Println("error 3", err)
				return err
			}
		}

		data := quic.StreamWriteParameters{
			Data:     buffer[:n],
			Finished: finished,
		}
		fmt.Println("writing data:  ", data)
		err = stream.Write(data)
		if err != nil {
			fmt.Println("error 4", err)
			return err
		}
		if finished {
			fmt.Printf("Finished writing to stream %d\n", stream.StreamID())
			return nil
		}
	}
}

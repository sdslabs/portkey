package session

import (
	"io"
	"io/ioutil"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/pion/quic"
	"github.com/sdslabs/portkey/pkg/utils"
)

const sendBufferSize = 100

func WriteLoop(stream *quic.BidirectionalStream, sendPath string, sendErr chan error, wg *sync.WaitGroup) error {
	defer wg.Done()

	zipfile, err := ioutil.TempFile(os.TempDir(), "portkey*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(zipfile.Name())
	err = utils.Zip(sendPath, zipfile)
	if err != nil {
		return err
	}

	if _, err = zipfile.Seek(0, 0); err != nil {
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
				return err
			}
		}

		data := quic.StreamWriteParameters{
			Data:     buffer[:n],
			Finished: finished,
		}
		err = stream.Write(data)
		if err != nil {
			return err
		}
		if finished {
			log.Infof("Finished writing to stream %d\n", stream.StreamID())
			return nil
		}
	}
}

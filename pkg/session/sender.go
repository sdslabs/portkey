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

	tempfile, err := ioutil.TempFile(os.TempDir(), "portkey*")
	if err != nil {
		log.WithError(err).Errorf("Error in sender stream: stream id = %d\n", stream.StreamID())
		return err
	}
	defer os.Remove(tempfile.Name())

	err = utils.Tar(sendPath, tempfile)
	if err != nil {
		log.WithError(err).Errorf("Error in sender stream: stream id = %d\n", stream.StreamID())
		return err
	}

	if _, err = tempfile.Seek(0, 0); err != nil {
		log.WithError(err).Errorf("Error in sender stream: stream id = %d\n", stream.StreamID())
		return err
	}
	finished := false
	buffer := make([]byte, sendBufferSize)
	for {
		n, err := tempfile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				finished = true
			} else {
				log.WithError(err).Errorf("Error in sender stream: stream id = %d\n", stream.StreamID())
				return err
			}
		}

		data := quic.StreamWriteParameters{
			Data:     buffer[:n],
			Finished: finished,
		}
		err = stream.Write(data)
		if err != nil {
			log.WithError(err).Errorf("Error in sender stream: stream id = %d\n", stream.StreamID())
			return err
		}
		if finished {
			log.Infof("Finished writing to stream %d\n", stream.StreamID())
			return nil
		}
	}
}

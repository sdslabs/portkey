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

const writeBufferSize = 100
const receiveBufferSize = 100

func ReadLoop(stream *quic.BidirectionalStream, receivePath string, receiveErr chan error, wg *sync.WaitGroup) error {
	defer wg.Done()
	tempfile, err := ioutil.TempFile(os.TempDir(), "portkey*")
	if err != nil {
		log.WithError(err).Errorf("Error in tempfile creation in receiver stream: stream id = %d\n", stream.StreamID())
		return err
	}
	defer os.Remove(tempfile.Name())

	buffer := make([]byte, receiveBufferSize)

	for {
		params, err := stream.ReadInto(buffer)
		if err != nil {
			if err != io.EOF {
				log.WithError(err).Errorf("Error in reading into buffer in receiver stream: stream id = %d\n", stream.StreamID())
				return err
			}
		}
		_, err = tempfile.Write(buffer[:params.Amount])
		if err != nil {
			log.WithError(err).Errorf("Error in writing fileBuffer in receiver stream: stream id = %d\n", stream.StreamID())
			return err
		}
		if params.Finished {
			break
		}
	}

	if receivePath == "" {
		receivePath, err = os.Getwd()
		if err != nil {
			log.WithError(err).Errorf("Error in finding working directory in receiver stream: stream id = %d\n", stream.StreamID())
			return err
		}
	}

	err = utils.Untar(tempfile, receivePath)
	if err != nil {
		log.WithError(err).Errorf("Error in untaring file in receiver stream: stream id = %d\n", stream.StreamID())
		return err
	}
	log.Infof("Finished reading from Stream %d\n", stream.StreamID())
	return nil
}

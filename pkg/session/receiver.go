package session

import (
	"io"
	"io/ioutil"
	"os"
	"sync"

	"github.com/DataDog/zstd"
	log "github.com/sirupsen/logrus"

	"github.com/pion/quic"
	"github.com/sdslabs/portkey/pkg/utils"
)

const receiveBufferSize = 100
const fileBufferSize = 3 * receiveBufferSize

func ReadLoop(stream *quic.BidirectionalStream, receivePath string, receiveErr chan error, wg *sync.WaitGroup) error {
	defer wg.Done()
	tempfile, err := ioutil.TempFile(os.TempDir(), "portkey*")
	if err != nil {
		return err
	}
	defer os.Remove(tempfile.Name())

	receiveBuffer := make([]byte, receiveBufferSize)
	fileBuffer := make([]byte, fileBufferSize)

	for {
		params, err := stream.ReadInto(receiveBuffer)
		if err != nil {
			if err != io.EOF {
				return err
			}
		}

		fileBuffer, err := zstd.Decompress(fileBuffer, receiveBuffer[:params.Amount])
		if err != nil {
			return err
		}

		_, err = tempfile.Write(fileBuffer)
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

	err = utils.Untar(tempfile, receivePath)
	if err == nil {
		log.Infof("Finished reading from Stream %d\n", stream.StreamID())
	}

	return nil
}

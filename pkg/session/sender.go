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

const compressBufferSize = 100

func WriteLoop(stream *quic.BidirectionalStream, sendPath string, sendErr chan error, wg *sync.WaitGroup) error {
	defer wg.Done()

	tempfile, err := ioutil.TempFile(os.TempDir(), "portkey*")
	if err != nil {
		return err
	}
	defer os.Remove(tempfile.Name())

	err = utils.Tar(sendPath, tempfile)
	if err != nil {
		return err
	}

	if _, err = tempfile.Seek(0, 0); err != nil {
		return err
	}
	finished := false
	fileBuffer := make([]byte, compressBufferSize)
	sendBuffer := make([]byte, compressBufferSize)

	for {
		bytesRead, err := tempfile.Read(fileBuffer)
		if err != nil {
			if err == io.EOF {
				finished = true
			} else {
				return err
			}
		}

		sendBuffer, err = zstd.Compress(sendBuffer, fileBuffer[:bytesRead])
		if err != nil {
			return err
		}

		data := quic.StreamWriteParameters{
			Data:     sendBuffer,
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

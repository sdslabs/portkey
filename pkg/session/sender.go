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

const sendBufferSize = 100

func WriteLoop(stream *quic.BidirectionalStream, sendPath string, sendErr chan error, wg *sync.WaitGroup) error {
	defer wg.Done()

	tempfile, err := ioutil.TempFile(os.TempDir(), "portkey*")
	if err != nil {
		log.WithError(err).Errorf("Error in tempfile creation in sender stream: stream id = %d\n", stream.StreamID())
		return err
	}
	defer os.Remove(tempfile.Name())

	err = utils.Tar(sendPath, tempfile)
	if err != nil {
		log.WithError(err).Errorf("Error in making tarball in sender stream: stream id = %d\n", stream.StreamID())
		return err
	}

	if _, err = tempfile.Seek(0, 0); err != nil {
		log.WithError(err).Errorf("Error in going to tempfile start in sender stream: stream id = %d\n", stream.StreamID())
		return err
	}

	var originalSize, compressedSize int64

	fs, err := tempfile.Stat()
	if err != nil {
		log.WithError(err).Errorf("Error in getting stat of tempfile in sender stream: stream id = %d\n", stream.StreamID())
		return err
	}

	originalSize = fs.Size()

	finished := false
	sendBuffer := make([]byte, sendBufferSize)
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()
	zstdWriter := zstd.NewWriter(pipeWriter)

	go func() {
		defer pipeWriter.Close()
		defer zstdWriter.Close()
		_, err = io.Copy(zstdWriter, tempfile)
		if err != nil {
			log.WithError(err).Errorf("Error in writing to tempfile in sender stream: stream id = %d\n", stream.StreamID())
		}
	}()

	compressedSize = 0
	for {
		bytesRead, err := pipeReader.Read(sendBuffer)
		if err != nil {
			if err == io.EOF {
				finished = true
			} else {
				log.WithError(err).Errorf("Error in writing sendBuffer in sender stream: stream id = %d\n", stream.StreamID())
				return err
			}
		}

		data := quic.StreamWriteParameters{
			Data:     sendBuffer[:bytesRead],
			Finished: finished,
		}
		compressedSize += int64(bytesRead)

		err = stream.Write(data)
		if err != nil {
			log.WithError(err).Errorf("Error in writing to sender stream: stream id = %d\n", stream.StreamID())
			return err
		}

		log.Infof("Wrote %d bytes to stream %d\n", bytesRead, stream.StreamID())

		if finished {
			log.Infof("Finished writing to stream %d\n", stream.StreamID())
			log.Infof("Original file size: %d bytes\n", originalSize)
			log.Infof("Compressed file size: %d bytes\n", compressedSize)
			return nil
		}
	}
}

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

// ReadLoop reads the data from the bidirectional stream and calls function for
// decompressing and untarring the files.
func ReadLoop(stream *quic.BidirectionalStream, receivePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	quicgoStream := stream.Detach()

	tempfile, err := ioutil.TempFile(os.TempDir(), "portkey*")
	if err != nil {
		log.WithError(err).Errorf("Error in tempfile creation in receiver stream %d\n", quicgoStream.StreamID())
		return
	}
	defer os.Remove(tempfile.Name())

	zstdReader := zstd.NewReader(quicgoStream)
	bytesWritten, err := io.Copy(tempfile, zstdReader)
	if err != nil {
		log.WithError(err).Errorf("Error in copying from zstdReader to tempfile in receiver stream %d\n", quicgoStream.StreamID())
	}
	log.Infof("Copied %d bytes from zstdReader to tempfile in receiver stream %d", bytesWritten, quicgoStream.StreamID())

	if err = zstdReader.Close(); err != nil {
		log.WithError(err).Errorf("Error in closing zstdWriter in receiver stream %d\n", quicgoStream.StreamID())
	}
	if err = quicgoStream.Close(); err != nil {
		log.WithError(err).Errorf("Error in closing stream %d\n", quicgoStream.StreamID())
	}
	log.Infof("Finished reading from stream %d\n", quicgoStream.StreamID())

	log.Infof("Untaring received file in receiver stream %d ...", quicgoStream.StreamID())
	if receivePath == "" {
		receivePath, err = os.Getwd()
		if err != nil {
			log.WithError(err).Errorf("Error in finding working directory in receiver stream %d\n", quicgoStream.StreamID())
			return
		}
	}

	err = utils.Untar(tempfile, receivePath)
	if err != nil {
		log.WithError(err).Errorf("Error in untaring file in receiver stream %d\n", quicgoStream.StreamID())
	}
}

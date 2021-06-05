package session

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"sync"

	"github.com/DataDog/zstd"
	log "github.com/sirupsen/logrus"

	"github.com/pion/quic"
	"github.com/sdslabs/portkey/pkg/utils"
)

// WriteLoop writes the data on the bidirectional stream after tarring and compressing
// it.
func WriteLoop(stream *quic.BidirectionalStream, sendPath string, wg *sync.WaitGroup) {
	defer wg.Done()

	quicgoStream := stream.Detach()

	tempfile, err := ioutil.TempFile(os.TempDir(), "portkey*")
	if err != nil {
		log.WithError(err).Errorf("Error in tempfile creation in sender stream %d\n", quicgoStream.StreamID())
		return
	}
	defer os.Remove(tempfile.Name())

	// Archiving the file(s).
	err = utils.Tar(sendPath, tempfile)
	if err != nil {
		log.WithError(err).Errorf("Error in making tarball in sender stream %d\n", quicgoStream.StreamID())
		return
	}

	if _, err = tempfile.Seek(0, 0); err != nil {
		log.WithError(err).Errorf("Error in going to tempfile start in sender stream %d\n", quicgoStream.StreamID())
		return
	}

	// Performing compression using zstd.
	zstdWriter := zstd.NewWriter(quicgoStream)
	bytesWritten, err := io.Copy(zstdWriter, tempfile)
	if err != nil {
		log.WithError(err).Errorf("Error in copying from tempfile to zstdWriter in sender stream %d\n", quicgoStream.StreamID())
	}
	log.Infof("Copied %d bytes from tempfile to zstdWriter in sender stream %d\n", bytesWritten, quicgoStream.StreamID())

	if err = zstdWriter.Close(); err != nil {
		log.WithError(err).Errorf("Error in closing zstdWriter in sender stream %d\n", quicgoStream.StreamID())
	}
	if err = quicgoStream.Close(); err != nil {
		log.WithError(err).Errorf("Error in closing sender stream %d", quicgoStream.StreamID())
	}
	log.Infof("Finished writing to sender stream %d", quicgoStream.StreamID())

	log.Infoln("Waiting for peer to close stream...")

	var dummyBuffer bytes.Buffer
	_, _ = io.Copy(&dummyBuffer, quicgoStream)
}

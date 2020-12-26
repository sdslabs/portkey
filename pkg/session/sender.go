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

func WriteLoop(stream *quic.BidirectionalStream, sendPath string, sendErr chan error, wg *sync.WaitGroup) error {
	defer wg.Done()

	quicgoStream := stream.Detach()

	tempfile, err := ioutil.TempFile(os.TempDir(), "portkey*")
	if err != nil {
		log.WithError(err).Errorf("Error in tempfile creation in sender stream %d\n", quicgoStream.StreamID())
		return err
	}
	defer os.Remove(tempfile.Name())

	err = utils.Tar(sendPath, tempfile)
	if err != nil {
		log.WithError(err).Errorf("Error in making tarball in sender stream %d\n", quicgoStream.StreamID())
		return err
	}

	if _, err = tempfile.Seek(0, 0); err != nil {
		log.WithError(err).Errorf("Error in going to tempfile start in sender stream %d\n", quicgoStream.StreamID())
		return err
	}

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

	log.Infoln("Waiting for peer to close stream...")

	var dummyBuffer bytes.Buffer
	_, err = io.Copy(&dummyBuffer, quicgoStream)
	if err != nil {
		log.WithError(err).Infof("Unexpected behaviour in sender stream %d: Received error other than EOF from peer", quicgoStream.StreamID())
	}
	return nil
}

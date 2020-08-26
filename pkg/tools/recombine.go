package tools

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

//RecombineChunks recombines the chunks received
func RecombineChunks() error {
	_, err := ioutil.TempFile(os.TempDir(), "portkey*.zip")
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join(os.TempDir(), "portkey*.zip"), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	var writePosition int64 = 0

	for j := uint64(0); j < 4; j++ {
		currentChunkFileName := "portkey_" + strconv.FormatUint(j, 10) + ".zip"
		newFileChunk, err := os.Open(currentChunkFileName)
		if err != nil {
			return err
		}
		defer newFileChunk.Close()
		chunkInfo, err := newFileChunk.Stat()
		if err != nil {
			return err
		}
		var chunkSize int64 = chunkInfo.Size()
		chunkBufferBytes := make([]byte, chunkSize)
		writePosition = writePosition + chunkSize

		reader := bufio.NewReader(newFileChunk)
		_, err = reader.Read(chunkBufferBytes)
		if err != nil {
			return err
		}
		_, err = file.Write(chunkBufferBytes)
		if err != nil {
			return err
		}
		file.Sync()
		chunkBufferBytes = nil
	}
	file.Close()
	return nil
}

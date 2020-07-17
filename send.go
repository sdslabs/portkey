package main

import (
	"os"
	"io/ioutil"
	"io"

	"github.com/pion/quic"
)

const sendBufferSize = 100 

func WriteLoop(stream *quic.BidirectionalStream, sendPath string) {
	zipfile, err := ioutil.TempFile(os.TempDir(), "portkey*.zip")
	if err != nil { panic(err) }
	defer os.Remove(zipfile.Name())
	err = zipit(sendPath, zipfile)
	if err != nil { panic(err) }

	zipfile.Seek(0, 0)
	finished := false
	buffer := make([]byte, sendBufferSize)
	for {
		n, err := zipfile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				finished = true
			} else {
				panic(err)
			}
		} 
	
		data := quic.StreamWriteParameters{
			Data: buffer[:n],
			Finished: finished,
		}
		err = stream.Write(data)
		if err != nil { panic(err) }
		if finished { break }
	}
}
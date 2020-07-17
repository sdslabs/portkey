package main

import (
	"os"
	"io"
	"io/ioutil"

	"github.com/pion/quic"
)
const receiveBufferSize = 100

func ReadLoop(stream *quic.BidirectionalStream,receivePath string) {
	zipfile, err := ioutil.TempFile(os.TempDir(), "portkey*.zip")
	if err != nil { panic(err) }
	defer os.Remove(zipfile.Name())

	buffer := make([]byte, receiveBufferSize)

	for {
		params, err := stream.ReadInto(buffer)
		if err != nil { 
			if err != io.EOF {
				panic(err) 
			}
		}
		_, err = zipfile.Write(buffer[:params.Amount])
		if err != nil { panic(err) }
		if params.Finished { break }
	}

	if receivePath == "" {
		receivePath, err = os.Getwd()
		if err != nil { panic(err) }
	}

	unzip(zipfile, receivePath)
}
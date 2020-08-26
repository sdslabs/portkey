package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

//ChunkFile splits a file into four equally sized files
func ChunkFile(source string) ([4]string, error) {
	fmt.Println("Dividing into chunks")
	var filenames [4]string
	file, err := os.Open(source)
	if err != nil {
		return filenames, err
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	totalParts := 4
	for i := 0; i < totalParts; i++ {
		partSize := fileSize / int64(totalParts)
		partBuffer := make([]byte, partSize)
		file.Read(partBuffer)
		filenames[i] = "portkeychunk_" + strconv.Itoa(i+1)
		if _, err := os.Create(filenames[i]); err != nil {
			return filenames, err
		}
		if err := ioutil.WriteFile(filenames[i], partBuffer, os.ModeAppend); err != nil {
			return filenames, err
		}
		dir, _ := os.Getwd()
		filenames[i] = filepath.Join(dir, filenames[i])
	}
	fmt.Print("Divided into chunks")
	return filenames, nil
}

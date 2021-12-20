package common

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// DataDir is data dir reader which implements datasetReader
type DataDir string

func (d DataDir) dataset() []dataReader {
	var dr []dataReader
	files, err := ioutil.ReadDir(string(d))
	if err != nil {
		panic(err)
	}
	for _, fileInfo := range files {
		if fileInfo.IsDir() || !strings.HasSuffix(fileInfo.Name(), ".json") {
			continue
		}

		file, err := os.Open(path.Join(string(d), fileInfo.Name()))
		if err != nil {
			panic(err)
		}

		dr = append(dr, &dataFileReader{
			File: file,
			name: strings.TrimSuffix(fileInfo.Name(), ".json"),
		})
	}
	return dr
}

type dataFileReader struct {
	*os.File
	name string
}

func (e *dataFileReader) dataName() string {
	return e.name
}

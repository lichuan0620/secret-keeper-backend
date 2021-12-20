package common

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"
)

// DataDirJoin is data dir aggregation reader which implements datasetReader
type DataDirJoin struct {
	Dir  string
	Name string
}

func (d DataDirJoin) dataset() []dataReader {
	files, err := ioutil.ReadDir(d.Dir)
	if err != nil {
		panic(err)
	}
	dataset := make([]interface{}, 0)
	for _, fileInfo := range files {
		if fileInfo.IsDir() || !strings.HasSuffix(fileInfo.Name(), ".json") {
			continue
		}

		dataSrc, err := ioutil.ReadFile(path.Join(d.Dir, fileInfo.Name()))
		if err != nil {
			panic(err)
		}

		fileData := make([]interface{}, 0)
		if err = json.Unmarshal(dataSrc, &fileData); err != nil {
			panic(err)
		}

		dataset = append(dataset, fileData...)
	}

	dataDst, err := json.Marshal(dataset)
	if err != nil {
		panic(err)
	}

	return []dataReader{&bytesReader{
		Reader: bytes.NewReader(dataDst),
		name:   d.Name,
	}}
}

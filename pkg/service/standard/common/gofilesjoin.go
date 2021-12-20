package common

import (
	"bytes"
	"encoding/json"
)

// GoFilesJoin is a go-file aggregation data reader which implements datasetReader
type GoFilesJoin struct {
	GoFiles GoFilesType
	Name    string
}

func (fj GoFilesJoin) dataset() []dataReader {
	dataset := make([]interface{}, 0)
	for _, f := range fj.GoFiles {
		fileData := make([]interface{}, 0)
		if err := json.Unmarshal(f.data, &fileData); err != nil {
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
		name:   fj.Name,
	}}
}

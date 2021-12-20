package common

import (
	"bytes"
)

// GoFilesType is go-file data reader which implements datasetReader
type GoFilesType []bytesReader

func (f GoFilesType) dataset() []dataReader {
	drs := make([]dataReader, 0)
	for i := range f {
		f[i].Reader = bytes.NewReader(f[i].data)
		drs = append(drs, &f[i])
	}

	return drs
}

// GoFiles is go-file error with all standard error data
var GoFiles GoFilesType

type bytesReader struct {
	*bytes.Reader
	data []byte
	name string
}

func (e *bytesReader) dataName() string {
	return e.name
}

func (e *bytesReader) Close() error {
	return nil
}

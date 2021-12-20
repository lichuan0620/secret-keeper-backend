package common

import (
	"html/template"
	"io"
	"os"
	"path"
)

// Target defines standard error target which can gen file by given parameters
type Target struct {
	TemplatePath  string
	OutDir        string
	DatasetReader datasetReader
	Unmarshaler   DataUnmarshaler
}

type datasetReader interface {
	dataset() []dataReader
}

type dataReader interface {
	dataName() string
	io.ReadCloser
}

// DataUnmarshaler defiles error data unmarshaler with name and data
type DataUnmarshaler interface {
	UnmarshalData(name string, data []byte) error
}

// GenFile can generate go file by given parameters
func (t *Target) GenFile() error {
	for _, reader := range t.DatasetReader.dataset() {
		data, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		if err = reader.Close(); err != nil {
			panic(err)
		}

		if err = t.Unmarshaler.UnmarshalData(reader.dataName(), data); err != nil {
			panic(err)
		}

		tpl, err := template.ParseFiles(t.TemplatePath)
		if err != nil {
			return err
		}
		file, err := os.OpenFile(path.Join(t.OutDir, "generated."+reader.dataName()+".go"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		if err = tpl.Execute(file, t.Unmarshaler); err != nil {
			return err
		}
		if err = file.Close(); err != nil {
			return err
		}
	}

	return nil
}

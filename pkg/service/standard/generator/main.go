package main

import (
	"html/template"

	"github.com/lichuan0620/secret-keeper-backend/pkg/service/standard/common"

	"github.com/pkg/errors"
)

type goFile struct {
	Name string
	Data template.HTML
}

func (d *goFile) UnmarshalData(name string, data []byte) error {
	d.Name = name
	d.Data = template.HTML(data)
	return nil
}

func genGoFileReader() error {
	target := common.Target{
		TemplatePath:  "generator/gofile.tmpl",
		OutDir:        "common/",
		DatasetReader: common.DataDir("data/"),
		Unmarshaler:   &goFile{},
	}
	if err := target.GenFile(); err != nil {
		return errors.WithMessage(err, "gen file")
	}
	return nil
}

func genErrors() error {
	target := common.Target{
		TemplatePath:  "generator/errors.tmpl",
		OutDir:        ".",
		DatasetReader: common.DataDir("data/"),
		Unmarshaler:   &common.ErrorData{},
	}

	if err := target.GenFile(); err != nil {
		return errors.WithMessage(err, "gen file")
	}
	return nil
}

func genErrorTest() error {
	target := common.Target{
		TemplatePath: "generator/errors_test.tmpl",
		OutDir:       ".",
		DatasetReader: common.DataDirJoin{
			Dir:  "data/",
			Name: "errors_test",
		},
		Unmarshaler: &common.ErrorData{},
	}
	if err := target.GenFile(); err != nil {
		return errors.WithMessage(err, "gen file")
	}
	return nil
}

func main() {
	if err := genGoFileReader(); err != nil {
		panic(err)
	}
	if err := genErrors(); err != nil {
		panic(err)
	}
	if err := genErrorTest(); err != nil {
		panic(err)
	}
}

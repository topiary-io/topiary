package main

import (
	"errors"
	"github.com/spf13/hugo/hugolib"
	"io/ioutil"
	"os"
)

func (p *Page) Save() error {
	page, err := hugolib.NewPage(p.Path)

	if err != nil {
		return err
	}

	page.SetSourceContent([]byte(p.Body))

	return page.SafeSaveSourceAs(p.Path)
}

type StaticManager interface {
	Read(fp string) (*Page, error)
	Create(fp string, content []byte) (*Page, error)
	Update(fp string, content []byte) (*Page, error)
	Delete(fp string) error
}

func (p Page) Read(filename string) (*Page, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{
		Path: filename,
		Body: string(content),
	}, nil
}

func (p Page) Create(fp string, content []byte) (*Page, error) {

	// create a new page
	page := &Page{
		Path: fp,
		Body: string(content),
	}

	// save page to disk
	err := page.Save()
	if err != nil {
		return nil, err
	}

	return page, nil
}

// UpdatePage changes the content of an existing page
func (p Page) Update(fp string, content []byte) (*Page, error) {

	// delete existing page
	err := p.Delete(fp)
	if err != nil {
		return nil, err
	}

	// create a new page
	page := &Page{
		Path: fp,
		Body: string(content),
	}

	// save page to disk
	err = page.Save()
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (p Page) Delete(fp string) error {

	// check that file exists
	info, err := os.Stat(fp)
	if err != nil {
		return err
	}

	// that file is a directory
	if info.IsDir() {
		return errors.New("DeletePage cannot delete directories")
	}

	// remove the directory
	return os.Remove(fp)
}

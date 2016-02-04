package main

import (
	"errors"
	"github.com/spf13/cast"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/parser"
	"os"
)

type (
	Frontmatter map[string]interface{}

	Content struct {
		Path     string
		Metadata Frontmatter
		Body     string
	}

	Page struct {
		Path string
		Body string
	}
)

func NewPage() *Page {
	return &Page{}
}

func (c *Content) contentSave() error {
	page, err := hugolib.NewPage(c.Path)

	if err != nil {
		return err
	}

	page.SetSourceMetaData(c.Metadata, '+')
	page.SetSourceContent([]byte(c.Body))

	return page.SafeSaveSourceAs(c.Path)
}

type PageManager interface {
	Read(fp string) (*Content, error)
	Create(fp string, fm Frontmatter, content []byte) (*Content, error)
	Update(fp string, fm Frontmatter, content []byte) (*Content, error)
	Delete(fp string) error
}

func (p Page) contentRead(filename string) (*Content, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	parser, err := parser.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	rawdata, err := parser.Metadata()
	if err != nil {
		return nil, err
	}

	metadata, err := cast.ToStringMapE(rawdata)
	if err != nil {
		return nil, err
	}

	return &Content{
		Path:     filename,
		Metadata: metadata,
		Body:     string(parser.Content()),
	}, nil
}

func (p Page) contentCreate(fp string, fm Frontmatter, content []byte) (*Content, error) {

	// create a new page
	page := &Content{
		Path:     fp,
		Metadata: fm,
		Body:     string(content),
	}

	// save page to disk
	err := page.contentSave()
	if err != nil {
		return nil, err
	}

	return page, nil
}

// UpdatePage changes the content of an existing page
func (p Page) contentUpdate(fp string, fm Frontmatter, content []byte) (*Content, error) {

	// delete existing page
	err := p.contentDelete(fp)
	if err != nil {
		return nil, err
	}

	// create a new page
	page := &Content{
		Path:     fp,
		Metadata: fm,
		Body:     string(content),
	}

	// save page to disk
	err = page.contentSave()
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (p Page) contentDelete(fp string) error {

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

func getTitle(fm Frontmatter) (string, error) {

	// check that title has been specified
	t, ok := fm["title"]
	if ok == false {
		return "", errors.New("page[meta].title must be specified")
	}

	// check that title is a string
	title, ok := t.(string)
	if ok == false {
		return "", errors.New("page[meta].title must be a string")
	}

	return title, nil
}

package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/unidoc/unidoc/pdf/creator"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

type PDFSource struct {
	source
}

func (s PDFSource) MergeTo(c *creator.Creator) error {
	debugInfo(fmt.Sprintf("Adding PDF: %+v", s.source))

	f, _ := os.Open(s.path)
	defer f.Close()

	return addPdfPages(f, s.pages, c)
}

func getReader(rs io.ReadSeeker) (*pdf.PdfReader, error) {

	pdfReader, err := pdf.NewPdfReader(rs)
	if err != nil {
		return nil, err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return nil, err
	}

	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return nil, err
		}
		if !auth {
			return nil, errors.New("Cannot merge encrypted, password protected document")
		}
	}

	return pdfReader, nil
}

func addPdfPages(file *os.File, pages []int, c *creator.Creator) error {
	pdfReader, err := getReader(file)
	if err != nil {
		return err
	}

	if len(pages) > 0 {
		debugInfo(fmt.Sprintf("Adding all pages of PDF: %s", file.Name()))
		for _, pageNo := range pages {
			if page, pageErr := pdfReader.GetPage(pageNo); pageErr != nil {
				return pageErr
			} else {
				err = c.AddPage(page)
			}
		}
	} else {
		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return err
		}
		for i := 0; i < numPages; i++ {
			pageNum := i + 1

			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				return err
			}

			debugInfo(fmt.Sprintf("Adding page %d of PDF: %s", pageNum, file.Name()))
			if err = c.AddPage(page); err != nil {
				return err
			}
		}
	}

	return nil
}

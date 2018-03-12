package main

import (
	"errors"
	"io"

	"github.com/unidoc/unidoc/pdf/creator"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

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

func addPdfPages(file io.ReadSeeker, pages []int, c *creator.Creator) error {
	pdfReader, err := getReader(file)
	if err != nil {
		return err
	}

	if len(pages) > 0 {
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

			if err = c.AddPage(page); err != nil {
				return err
			}
		}
	}

	return nil
}

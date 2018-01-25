package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	unicommon "github.com/unidoc/unidoc/common"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func init() {
	// Debug log level.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Requires at least 3 arguments: output_path and 2 input paths(and optional page numbers) \n")
		fmt.Printf("Usage: merge2pdf output.pdf input1.pdf input2.pdf~1,2,3 ...\n")
		os.Exit(0)
	}

	outputPath := os.Args[1]
	inputPaths := []string{}
	inputPages := [][]int{}

	// Sanity check the input arguments.
	for _, arg := range os.Args[2:] {
		//inputPaths = append(inputPaths, arg)

		fileInputParts := strings.Split(arg, "~")
		inputPaths = append(inputPaths, fileInputParts[0])
		pages := []int{}

		if len(fileInputParts) > 1 {
			for _, e := range strings.Split(fileInputParts[1], ",") {
				pageNo, err := strconv.Atoi(strings.Trim(e, " \n"))
				if err != nil {
					fmt.Errorf("Invalid format! Example of a file input with page numbers: path/to/abc.pdf~1,2,3,5,6")
					os.Exit(1)
				}
				pages = append(pages, pageNo)
			}
		}

		inputPages = append(inputPages, pages)
	}

	// fmt.Println(inputPages)
	// os.Exit(1)

	err := mergePdf(inputPaths, inputPages, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Complete, see output file: %s\n", outputPath)
}

func mergePdf(inputPaths []string, inputPages [][]int, outputPath string) error {
	pdfWriter := pdf.NewPdfWriter()

	for i, inputPath := range inputPaths {

		f, err := os.Open(inputPath)
		if err != nil {
			return err
		}
		defer f.Close()

		fileType, typeError := getFileType(f)
		if typeError != nil {
			return nil
		}

		if fileType == "directory" {
			// @TODO : Read all files in directory
			return errors.New(inputPath + " is a drectory.")
		} else if fileType == "application/pdf" {
			err := addPdfPages(f, inputPages[i], &pdfWriter)
			if err != nil {
				return err
			}
		} else if ok, _ := in_array(fileType, []string{"image/jpg", "image/jpeg", "image/png"}); ok {
			return errors.New(inputPath + " Images is not supproted yet.")
		}

	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
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

func addPdfPages(file io.ReadSeeker, pages []int, writer *pdf.PdfWriter) error {
	pdfReader, err := getReader(file)
	if err != nil {
		return err
	}

	if len(pages) > 0 {
		for _, pageNo := range pages {
			if page, pageErr := pdfReader.GetPage(pageNo); pageErr != nil {
				return pageErr
			} else {
				err = writer.AddPage(page)
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

			if err = writer.AddPage(page); err != nil {
				return err
			}
		}
	}

	return nil
}

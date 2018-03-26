package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/unidoc/unidoc/pdf/creator"
)

type Mergeable interface {
	MergeTo(c *creator.Creator) error
}

type source struct {
	path, sourceType, mime, ext string
	pages                       []int
}

// Initiate new source file from input argument
func NewSource(input string) Mergeable {
	fileInputParts := strings.Split(input, "~")

	path := fileInputParts[0]
	var inputSource Mergeable

	info, err := os.Stat(path)
	if err != nil {
		log.Fatal("Error:", err.Error())
	}

	switch mode := info.Mode(); {
	case mode.IsDir():
		inputSource = getMergeableDir(path)
	case mode.IsRegular():
		pages := []int{}
		if len(fileInputParts) > 1 {
			pages = parsePageNums(fileInputParts[1])
		}
		inputSource = getMergeableFile(path, pages)
	}

	return inputSource
}

func getMergeableDir(path string) Mergeable {
	dir := DirSource{path: path}
	dir.scanMergeables()

	return dir
}

func getMergeableFile(path string, pages []int) Mergeable {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Cannot read source file:", path)
	}

	defer f.Close()

	ext := filepath.Ext(f.Name())
	mime, err := getMimeType(f)
	if err != nil {
		log.Fatal("Error in getting mime type of file:", path)
	}

	sourceType, err := getFileType(mime, ext)
	if err != nil {
		log.Fatal("Error : %s (%s)", err.Error(), path)
	}

	source := source{path, sourceType, mime, ext, pages}

	var m Mergeable
	switch sourceType {
	case "image":
		m = ImgSource{source}
	case "pdf":
		m = PDFSource{source}
	}

	return m
}

func getFileType(mime, ext string) (string, error) {
	pdfExts := []string{".pdf", ".PDF"}
	imgExts := []string{".jpg", ".jpeg", ".gif", ".png", ".tiff", ".tif", ".JPG", ".JPEG", ".GIF", ".PNG", ".TIFF", ".TIF"}

	switch {
	case mime == "application/pdf":
		return "pdf", nil
	case mime[:6] == "image/":
		return "image", nil
	case mime == "application/octet-stream" && in_array(ext, pdfExts):
		return "pdf", nil
	case mime == "application/octet-stream" && in_array(ext, imgExts):
		return "image", nil
	}

	return "error", errors.New("File type not acceptable. ")
}

func parsePageNums(pagesInput string) []int {
	pages := []int{}

	for _, e := range strings.Split(pagesInput, ",") {
		pageNo, err := strconv.Atoi(strings.Trim(e, " \n"))
		if err != nil {
			fmt.Errorf("Invalid format! Example of a file input with page numbers: path/to/abc.pdf~1,2,3,5,6")
			os.Exit(1)
		}
		pages = append(pages, pageNo)
	}

	return pages
}

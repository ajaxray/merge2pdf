package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	flag "github.com/ogier/pflag"
	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/creator"
)

var size, margin string
var scaleH, scaleW bool

const (
	DefaultSize   = "IMG-SIZE"
	DefaultMargin = "1,1,1,1"
)

func init() {
	// Debug log level.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	flag.StringVarP(&size, "size", "s", DefaultSize, "Resize image pages to print size. One of A4, A3, Legal or Letter.")
	flag.StringVarP(&margin, "margin", "m", DefaultMargin, "Comma separated numbers for left, right, top, bottm side margin in inch.")
	flag.BoolVarP(&scaleW, "scale-width", "w", false, "Scale Image to page width. Only if --size specified.")
	flag.BoolVarP(&scaleH, "scale-height", "h", false, "Scale Image to page height. Only if --size specified.")

	flag.Usage = func() {
		fmt.Println("Requires at least 3 arguments: output_path and 2 input paths (and optional page numbers)")
		fmt.Println("Usage: merge2pdf output.pdf input1.pdf input2.pdf~1,2,3 input3.jpg /path/to/dir ...")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	outputPath := args[0]
	inputPaths := []string{}
	inputPages := [][]int{}

	for _, arg := range flag.Args()[1:] {
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

	// fmt.Println(inputPaths)
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
	c := creator.New()

	for i, inputPath := range inputPaths {

		f, err := os.Open(inputPath)
		if err != nil {
			return err
		}
		defer f.Close()

		fileType, typeError := getFileType(f)
		if typeError != nil {
			return typeError
		}

		if fileType == "directory" {
			// @TODO : Read all files in directory
			return errors.New(inputPath + " is a drectory.")
		} else if fileType == "application/pdf" {
			err = addPdfPages(f, inputPages[i], c)
		} else if fileType[:6] == "image/" {
			err = addImage(inputPath, c)

			fmt.Println("%+v", pageMargin)
			fmt.Println("%+v", pageSize)

		} else {
			err = errors.New("Unsupported type:" + inputPath)
		}

		if err != nil {
			return err
		}
	}

	err := c.WriteToFile(outputPath)
	if err != nil {
		return err
	}

	return nil
}

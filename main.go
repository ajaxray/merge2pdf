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
var scaleH, scaleW, verbose, version bool
var JPEGQuality int

const (
	DefaultSize   = "IMG-SIZE"
	DefaultMargin = "1,1,1,1"
	VERSION       = "1.0.1"
)

func init() {

	flag.StringVarP(&size, "size", "s", DefaultSize, "Resize image pages to print size. One of A4, A3, Legal or Letter.")
	flag.StringVarP(&margin, "margin", "m", DefaultMargin, "Comma separated numbers for left, right, top, bottm side margin in inch.")
	flag.BoolVar(&scaleW, "scale-width", false, "Scale Image to page width. Only if --size specified.")
	flag.BoolVar(&scaleH, "scale-height", false, "Scale Image to page height. Only if --size specified.")
	flag.IntVar(&JPEGQuality, "jpeg-quality", 100, "Optimize JPEG Quality.")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Display debug info.")
	flag.BoolVarP(&version, "version", "V", false, "Display debug info.")

	flag.Usage = func() {
		fmt.Println("merge2pdf [options...] <output_file> <input_file> [<input_file>...]")
		fmt.Println("<output_file> should be a PDF file and <input_file> can be PDF or image files. PDF files can have specific page numbers to merge.")
		fmt.Println("Example: merge2pdf output.pdf input1.pdf input_pages.pdf~1,2,3 input3.jpg /path/to/input.png ...")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("merge2pdf version ", VERSION)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	if verbose {
		unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
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

	debugInfo(fmt.Sprintf("Complete, see output file: %s", outputPath))
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
			err = addImage(inputPath, c, fileType)
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

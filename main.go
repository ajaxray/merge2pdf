package main

import (
	"fmt"
	"os"

	flag "github.com/ogier/pflag"
	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/creator"
)

var size, margin string
var scaleH, scaleW, landscape, verbose, version bool
var JPEGQuality int

const (
	DefaultSize   = "IMG-SIZE"
	DefaultMargin = "0,0,0,0"
	VERSION       = "1.2.0"
)

func init() {

	flag.StringVarP(&size, "size", "s", DefaultSize, "Resize image pages to print size. One of A4, A3, Legal or Letter.")
	flag.StringVarP(&margin, "margin", "m", DefaultMargin, "Comma separated numbers for left, right, top, bottm side margin in inch.")
	flag.BoolVar(&scaleW, "scale-width", false, "Scale Image to page width. Only if --size specified.")
	flag.BoolVar(&scaleH, "scale-height", false, "Scale Image to page height. Only if --size specified.")
	flag.IntVar(&JPEGQuality, "jpeg-quality", 100, "Optimize JPEG Quality.")
	flag.BoolVarP(&landscape, "landscape", "l", false, "Set page size in Landscape mode.")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Display debug information.")
	flag.BoolVarP(&version, "version", "V", false, "Display Version information.")

	flag.Usage = func() {
		fmt.Println("merge2pdf [options...] <output_file> <input_file> [<input_file>...]")
		fmt.Println("<output_file> should be a PDF file and <input_file> can be PDF, image or directory. PDF files can have specific page numbers to merge.")
		fmt.Println("Example: merge2pdf output.pdf input1.pdf input_pages.pdf~1,2,3 input3.jpg /path/to/dir ...")
		fmt.Println("Available Options: ")
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
	inputPaths := args[1:]

	c := creator.New()

	for _, arg := range inputPaths {
		err := NewSource(arg).MergeTo(c)
		if err != nil {
			fmt.Printf("Error: %s (%s) \n", err.Error(), arg)
			os.Exit(1)
		}
	}

	err := c.WriteToFile(outputPath)
	if err != nil {
		fmt.Printf("Error: %s \n", err.Error())
	}

	debugInfo(fmt.Sprintf("Complete, see output file: %s", outputPath))
	os.Exit(0)
}

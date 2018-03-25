package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/unidoc/unidoc/pdf/core"
	"github.com/unidoc/unidoc/pdf/creator"
)

var pageMargin [4]float64
var pageSize creator.PageSize
var sizeHasSet, marginHasSet = false, false

type ImgSource struct {
	source
}

func (s ImgSource) MergeTo(c *creator.Creator) error {
	return addImage(s.path, c, s.mime)
}

func addImage(filePath string, c *creator.Creator, fileType string) error {
	debugInfo(fmt.Sprintf("Adding image: %s", filePath))

	img, err := creator.NewImageFromFile(filePath)
	if err != nil {
		return err
	}

	// The following funcs must be called in sequence
	setEncoding(img, fileType)
	setMargin(img, c)
	setSize(img, c)

	c.NewPage()
	err = c.Draw(img)
	if err != nil {
		return err
	}

	return nil
}

func setMargin(img *creator.Image, c *creator.Creator) {

	if !marginHasSet {
		for i, m := range strings.Split(margin, ",") {
			floatVal, err := strconv.ParseFloat(m, 64)
			if err != nil {
				log.Fatalln("Error: -m|--margin MUST be 4 comma separated int/float numbers. %s found.", m)
			}

			pageMargin[i] = floatVal * creator.PPI
		}
		marginHasSet = true
	}

	c.SetPageMargins(pageMargin[0], pageMargin[1], pageMargin[2], pageMargin[3])
	img.SetPos(pageMargin[0], pageMargin[3])
}

func setSize(img *creator.Image, c *creator.Creator) {
	if size == DefaultSize {
		// sizeHasSet will remain false so that it changes for each image

		// Width height with adding margin
		w := img.Width() + pageMargin[0] + pageMargin[1]
		h := img.Height() + pageMargin[2] + pageMargin[3]

		pageSize = creator.PageSize{w, h}
	} else {
		sizeHasSet = true
		switch size {
		case "A4":
			pageSize = creator.PageSizeA4
		case "A3":
			pageSize = creator.PageSizeA3
		case "Legal":
			pageSize = creator.PageSizeLegal
		case "Letter":
			pageSize = creator.PageSizeLetter
		default:
			log.Fatalln("Error: -s|--size MUST be one of A4, A3, Legal or Letter. %s given.", size)
		}

		if scaleH {
			img.ScaleToHeight(pageSize[1] - (pageMargin[2] + pageMargin[3]))
		} else if scaleW {
			img.ScaleToWidth(pageSize[0] - (pageMargin[0] + pageMargin[1]))
		}
	}

	debugInfo(fmt.Sprintf("Created Page Size: %v", pageSize))
	c.SetPageSize(pageSize)
}

// Set appropriate encoding for JPEG and TIFF
// MUST be called before changing image size
func setEncoding(img *creator.Image, fileType string) {

	switch fileType {
	case "image/jpeg":
		debugInfo(fmt.Sprintf("JPEG found. Setting encoding to DCTEncoder."))
		encoder := core.NewDCTEncoder()
		encoder.Quality = JPEGQuality
		// Encoder dimensions must match the raw image data dimensions.
		encoder.Width = int(img.Width())
		encoder.Height = int(img.Height())
		img.SetEncoder(encoder)
	}
}

package main

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/unidoc/unidoc/pdf/creator"
)

func TestSetSize(t *testing.T) {
	sizes := map[string]creator.PageSize{
		"A4":     creator.PageSizeA4,
		"A3":     creator.PageSizeA3,
		"Legal":  creator.PageSizeLegal,
		"Letter": creator.PageSizeLetter,
	}

	var i creator.Image
	c := creator.New()
	check := func(name string, width, height float64) {
		size = name
		setSize(&i, c)
		c.NewPage()

		assert.Equal(t, c.Context().PageWidth, width)
		assert.Equal(t, c.Context().PageHeight, height)
	}

	for sizeName, sizeVal := range sizes {
		check(sizeName, sizeVal[0], sizeVal[1])
	}

	// Now test in landscape mode
	landscape = true
	for sizeName, sizeVal := range sizes {
		check(sizeName, sizeVal[1], sizeVal[0])
	}

}

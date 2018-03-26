package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/unidoc/unidoc/pdf/creator"
)

type DirSource struct {
	path  string
	files []Mergeable
}

func (s DirSource) MergeTo(c *creator.Creator) error {
	debugInfo(fmt.Sprintf("Adding Directory: %s", s.path))
	var err error

	for _, mergeableFile := range s.files {
		err = mergeableFile.MergeTo(c)
		if err != nil {
			break
		}
	}

	return err
}

func (s *DirSource) scanMergeables() {
	filesInDir, _ := ioutil.ReadDir(s.path)

	for _, file := range filesInDir {
		if file.IsDir() {
			continue
		}

		s.files = append(s.files, getMergeableFile(filepath.Join(s.path, file.Name()), []int{}))
	}
}

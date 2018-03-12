package main

import (
	"log"
	"net/http"
	"os"
	"reflect"
)

func in_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func getFileType(file *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	if info, stateErr := file.Stat(); stateErr != nil {
		return "error", stateErr
	} else if info.IsDir() {
		return "directory", nil
	} else {
		buffer := make([]byte, 512)
		_, readError := file.Read(buffer)
		if readError != nil {
			return "error", readError
		}

		// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
		return http.DetectContentType(buffer), nil
	}
}

func debugInfo(message string) {
	if verbose {
		log.Println(message)
	}
}

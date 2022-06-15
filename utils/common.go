//nolint:wrapcheck
package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"time"
)

func Log(data interface{}, overwrite bool) error {
	jsonData := time.Now().Format("2006-01-02 15:04:05") + "  "

	myJSON, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	jsonData += string(myJSON)

	if overwrite {
		err = ioutil.WriteFile("/tmp/debug", []byte(jsonData), 0o600)
		if err != nil {
			return err
		}
	} else {
		f, err := os.OpenFile("/tmp/debug", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()

		_, err = f.Write([]byte(jsonData))
		if err != nil {
			return err
		}

		_, err = f.WriteString("\n\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func IsItemExistedInSlice(item interface{}, slice interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func GetDiff2Slices(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for j := range slice1 {
			found := false

			for k := range slice2 {
				if slice1[j] == slice2[k] {
					found = true

					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, slice1[j])
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

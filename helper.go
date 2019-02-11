package exsort

import (
	"fmt"
	"os"
)

func writeToFile(file *os.File, data interface{}) error {
	//TODO
	// Add customization
	_, err := file.WriteString(fmt.Sprintf("%v\n", data))
	if err != nil {
		return err
	}

	return nil
}

func openFile(path string) (*os.File, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(path)
			if err != nil {
				return nil, err
			}

			return file, nil
		} else {
			return nil, err
		}
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

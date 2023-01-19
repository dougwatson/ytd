package main

import (
	"os"
)

func writeFile(destFile string, b []byte) error {
	f, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	return err
}

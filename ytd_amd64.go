package main

import (
	"testing/fstest"
)

type myMapFS struct {
	fstest.MapFS
}

func GetFS() (myMapFS, error) {
	//just an empty stub so that ytd will compile when building for mac architecture
	return myMapFS{}, nil
}
func (m myMapFS) AddFile(name, content string) error {
	//stub
	return nil
}
func (m myMapFS) AddDir(name string) error {
	//stub
	return nil
}
func runRemote(ptrString string, args []string) {
	println("runRemote NOOP ptrString=", ptrString, args)
}

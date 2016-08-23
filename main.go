// Copyright (c) 2016, Michael Sonntag (sonntag@bio.lmu.de)
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted under the terms of the BSD License. See
// LICENSE file in the root of the Project.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
)

func main() {

	basePath := os.Getenv("GOPATH")
	if basePath == "" {
		panic("Cannot find gopath")
	}
	resPath := filepath.Join(basePath, "src", "spielwiese", "dentry", "res")

	inFile := filepath.Join(resPath, "parse.txt")
	outFile := filepath.Join(resPath, "out")

	fmt.Printf("InFile: '%s', outFile: '%s'\n", inFile, outFile)

	cont, err := ioutil.ReadFile(inFile)
	fmt.Printf("file content: \n---\n%s\n---\n", string(cont))


	fmt.Println("dentry started")

	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Printf("Error occurred: '%s'", err.Error())
	}
	win.SetTitle("Other Dentry")

	win.SetDefaultSize(800,600)

	// required to end program properly; first string needs to be as supported signal e.g. "destroy"
	win.Connect("destroy", func(){
		gtk.MainQuit()
	})

	label, err := gtk.LabelNew(fmt.Sprintf("content of file:\n%s", string(cont)))
	if err != nil {
		fmt.Println("Error creating label")
	}

	win.Add(label)

	// required to show window
	win.ShowAll()

	// required to display window
	gtk.Main()
}

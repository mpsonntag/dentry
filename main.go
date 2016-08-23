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
	"bytes"
	"bufio"
	"strings"
)

type tagEntList struct {
	content []tagEnt
}

type tagEnt struct{
	tags []string
	body string
}

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
	if err != nil {
		panic(fmt.Sprintf("Error reading file: '%s'", err.Error()))
	}

	err = textToEnt(&cont)
	if err != nil {
		panic(fmt.Sprintf("Error splitting content: '%s'", err.Error()))
	}

	app(&cont)
}

func app(cont *[]byte) {
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Printf("Error occurred: '%s'", err.Error())
	}
	win.SetTitle("Dentry")

	win.SetDefaultSize(800,600)

	// required to end program properly; first string needs to be as supported signal e.g. "destroy"
	win.Connect("destroy", func(){
		gtk.MainQuit()
	})

	label, err := gtk.LabelNew(fmt.Sprintf("content of file:\n%s", string(*cont)))
	if err != nil {
		fmt.Println("Error creating label")
	}

	win.Add(label)

	// required to show window
	win.ShowAll()

	// required to display window
	gtk.Main()
}

func textToEnt(cont *[]byte) error {
	tmp := &tagEntList{}

	r := bytes.NewReader(*cont)
	s := bufio.NewScanner(r)
	s.Split(splitOnHash)
	for s.Scan() {
		curr := strings.Replace(s.Text(), "(#)", "", 1)
		currParts := strings.Split(curr,"#)")
		if len(currParts) > 1 {
			currTags := strings.Split(currParts[0], ",")

			for i := range currTags {
				currTags[i] = strings.TrimSpace(currTags[i])
			}
			t := tagEnt{
				tags: currTags,
				body: currParts[1],
			}
			tmp.content = append(tmp.content, t)
		}
	}
	if err := s.Err(); err != nil {
		return  err
	}

	for _, entry := range tmp.content {
		fmt.Printf("\tTags: '%s'\n\tcontent: '%s'\n", entry.tags, entry.body)
	}

	return nil
}

func splitOnHash(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if(data[i] == '(' && data[i+1] == '#' && data[i+2] == ')') {
			return i+3, data[:i+3], nil
		}
	}
	return 0, data, bufio.ErrFinalToken
}

// Copyright (c) 2016, Michael Sonntag (sonntag@bio.lmu.de)
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted under the terms of the BSD License. See
// LICENSE file in the root of the Project.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

// tagEnt contains information stored in data and data associated with
// this information stored in tags
type tagEnt struct {
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

	tagList, err := textToEnt(&cont)
	if err != nil {
		panic(fmt.Sprintf("Error splitting content: '%s'", err.Error()))
	}

	app(tagList)
}

// app is the main function to open the GUI application
func app(tagList *[]tagEnt) {
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Printf("Error occurred: '%s'", err.Error())
	}
	win.SetTitle("Dentry")

	win.SetDefaultSize(800, 600)

	// required to end program properly; first string needs to be as supported signal e.g. "destroy"
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	grid, err := gtk.GridNew()
	if err != nil {
		fmt.Printf("Error creating grid: %s\n", err.Error())
	}

	for i, ent := range *tagList {
		label, err := gtk.LabelNew(ent.body)
		if err != nil {
			fmt.Println("Error creating label")
			return
		}

		grid.Attach(label, 0, i, 1, 1)

		hboxTag, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
		if err != nil {
			fmt.Printf("Error creating hboxTag: %s", err.Error())
		}

		for _, tag := range ent.tags {
			but, err := gtk.ButtonNew()
			if err != nil {
				fmt.Printf("Error creating button: %s", err.Error())
				return
			}
			but.SetLabel(tag)
			hboxTag.Add(but)
		}

		grid.Attach(hboxTag, 1, i, 1, 1)
	}

	win.Add(grid)

	// required to show window
	win.ShowAll()

	// required to display window
	gtk.Main()
}

// testToEnt scans a byte array and splits the content at '(#)' and removes the '(#)' occurrence.
// The resulting pieces are further split at '#)'. If '#)' exists, the first part is further
// split at ',' occurrences, the individual pieces are trimmed of whitespaces and
// stored in the tags field of a new tagEnt instance. The second part is stored in the
// body part of the tagEnt instance.
// All new tagEnt instances are stored in a tagEntList instance and returned if no error
// occurred.
func textToEnt(cont *[]byte) (*[]tagEnt, error) {
	tmp := make([]tagEnt, 0, 32)

	r := bytes.NewReader(*cont)
	s := bufio.NewScanner(r)
	s.Split(splitOnHash)
	for s.Scan() {
		curr := strings.Replace(s.Text(), "(#)", "", 1)
		currParts := strings.Split(curr, "#)")
		if len(currParts) > 1 {
			currTags := strings.Split(currParts[0], ",")

			for i := range currTags {
				currTags[i] = strings.TrimSpace(currTags[i])
			}
			t := tagEnt{
				tags: currTags,
				body: currParts[1],
			}
			tmp = append(tmp, t)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	// TODO for testing only, remove later
	for _, entry := range tmp {
		fmt.Printf("\tTags: '%s'\n\tcontent: '%s'\n", entry.tags, entry.body)
	}

	return &tmp, nil
}

// splitOnHash is a function satisfying bufio SplitFunc
// splitting a byte array at '(#)'.
func splitOnHash(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == '(' && data[i+1] == '#' && data[i+2] == ')' {
			return i + 3, data[:i+3], nil
		}
	}
	return 0, data, bufio.ErrFinalToken
}

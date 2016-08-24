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
	"spielwiese/dentry/lib"
)

func main() {

	win, err := appStart()
	if err != nil {
		panic(fmt.Sprintf("Error creating main window: %s\n", err.Error()))
	}

	basePath := os.Getenv("GOPATH")
	if basePath == "" {
		panic("Cannot find gopath")
	}
	resPath := filepath.Join(basePath, "src", "spielwiese", "dentry", "res")
	inFile := filepath.Join(resPath, "parse.txt")
	cont, err := ioutil.ReadFile(inFile)
	if err != nil {
		panic(fmt.Sprintf("Error reading file: '%s'", err.Error()))
	}

	ok, err := lib.IsTagNote(&cont)
	if err != nil {
		fmt.Printf("Error trying to check file header: %s\n", err.Error())
	}

	if ok {
		tagList, err := lib.TextToEnt(&cont)
		if err != nil {
			panic(fmt.Sprintf("Error splitting content: '%s'", err.Error()))
		}
		appShowTags(win, tagList)
	} else {
		appConsolatoryWin(win)
	}
}

// appStart is the main function to open the GUI application
func appStart() (*gtk.Window, error) {
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}
	win.SetTitle("Dentry")
	win.SetDefaultSize(800, 600)

	// required to end program properly; first string needs to be as supported signal e.g. "destroy"
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	return win, nil
}

// appShowTag displays a list of tags in the main window
func appShowTags(win *gtk.Window, tagList *[]lib.TagEnt) {
	grid, err := gtk.GridNew()
	if err != nil {
		fmt.Printf("Error creating grid: %s\n", err.Error())
	}
	grid.SetRowSpacing(10)
	grid.SetColumnSpacing(10)

	for i, ent := range *tagList {
		label, err := gtk.LabelNew(ent.Content)
		if err != nil {
			fmt.Printf("Error creating label: %s\n", err.Error())
			return
		}

		label.SetHAlign(gtk.ALIGN_START)
		grid.Attach(label, 0, i, 1, 1)

		hboxTag, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
		if err != nil {
			fmt.Printf("Error creating hboxTag: %s", err.Error())
		}

		for _, tag := range ent.Tags {
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

	fn, err := runFileChooser(win)
	if err != nil {
		fmt.Printf("Error running file chooser: %s\n", err.Error())
	}
	fmt.Printf("The following file has been chosen: %s\n", fn)

	// required to show window
	win.ShowAll()

	// required to display window
	gtk.Main()
}

// appConsolatoryWin shows a label, if no tags were loaded.
func appConsolatoryWin(win *gtk.Window) {
	lbl, err := gtk.LabelNew("There were no tags to be displayed, sorry!")
	if err != nil {
		fmt.Printf("Error creating consolatory label: %s", err.Error())
		return
	}
	win.Add(lbl)
	// required to show window
	win.ShowAll()
	// required to display window
	gtk.Main()
}

// runFileChooser creates a gtk FileChooserDialog and returns
// the path of the selected file as string.
func runFileChooser(win *gtk.Window) (string, error) {

	var fn string

	openFile, err := gtk.FileChooserDialogNewWith2Buttons("Open file", win, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Cancel", gtk.RESPONSE_CANCEL,
		"Ok", gtk.RESPONSE_OK)
	if err != nil {
		return "", err
	}

	openFile.SetDefaultSize(50, 50)

	res := openFile.Run()

	if res == int(gtk.RESPONSE_OK) {
		fn = openFile.FileChooser.GetFilename()
		openFile.Destroy()
	} else if res == int(gtk.RESPONSE_DELETE_EVENT) {
		openFile.Destroy()
	} else if res == int(gtk.RESPONSE_CANCEL) {
		openFile.Destroy()
	}

	return fn, nil
}

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

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mpsonntag/dentry/lib"
)

// LogLevel is the definition of various logging levels.
type LogLevel int

const (
	// ERR LogLevel is logged to the Error writer.
	ERR LogLevel = iota
	// WARN LogLevel is logged to the Standard writer.
	WARN
	// INFO LogLevel is logged to the Standard writer.
	INFO
	// DEBUG LogLevel is logged to the Standard writer.
	DEBUG
)
func main() {

	gtk.Init(nil)

	// application ID has to adhere to the rules defined in g_application_id_is_valid
	app, err := gtk.ApplicationNew("org.mps.dentry", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		panic(fmt.Sprintf("Error creating application: %s", err.Error()))
	}

	app.Connect("startup", startup)
	app.Connect("activate", createWin)

	// starts the main loop of the application, waiting for sthg to happen
	gtk.Main()
}

func startup() {
	fmt.Println("do I need this?")
}

func createWin(app *gtk.Application) error {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		panic(fmt.Sprintf("Error creating main window: %s\n", err.Error()))
	}
	// required to end program properly; first string needs to be as supported signal e.g. "destroy"
	// don't know if this is the proper point or way to end the application, but for now it works.
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.SetTitle("Dentry")
	win.SetDefaultSize(800, 600)

	btn, err := gtk.ButtonNewWithLabel("Parse file")
	if err != nil {
		return err
	}

	btn.Connect("clicked", handleFiles, win)
	win.Add(btn)
	win.ShowAll()

	/*
		err = appStart(win)
		if err != nil {
			panic(fmt.Sprintf("Error populating main window: %s\n", err.Error()))
		}
	*/
	app.AddWindow(win)

	return nil
}

// appStart is the main function to open the GUI application
func appStart(win *gtk.Window) error {

	return nil
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

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		fmt.Printf("Error creating vbox: %s\n", err.Error())
		return
	}
	vbox.Add(grid)

	btn, err := gtk.ButtonNewWithLabel("Try again")
	if err != nil {
		fmt.Printf("Error creating file reload button: %s\n", err.Error())
		return
	}
	btn.Connect("clicked", handleFiles, win)

	vbox.Add(btn)
	win.Add(vbox)
	win.ShowAll()
}

// appConsolatoryWin shows a label, if no tags were loaded.
func appConsolatoryWin(win *gtk.Window) {
	lbl, err := gtk.LabelNew("There were no tags to be displayed, sorry!")
	if err != nil {
		fmt.Printf("Error creating consolatory label: %s", err.Error())
		return
	}

	btn, err := gtk.ButtonNewWithLabel("Try again")
	if err != nil {
		fmt.Printf("Error creating file reload button: %s\n", err.Error())
		return
	}
	btn.Connect("clicked", handleFiles, win)

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		fmt.Printf("Error creating vbox: %s\n", err.Error())
	}

	vbox.Add(lbl)
	vbox.Add(btn)
	win.Add(vbox)
	win.ShowAll()
}

// handleFiles runs a file chooser and opens and parses any chosen file.
// If the file is of a parsable format, the results will be displayed,
// otherwise the file chooser can be opened again.
func handleFiles(_ *gtk.Button, win *gtk.Window) error {
	fp, err := runFileChooser(win)
	if err != nil {
		fmt.Printf("Error selecting file: %s\n", err.Error())
		return err
	}

	cont, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Printf("Error reading file: '%s'\n", err.Error())
		return err
	}

	ok, err := lib.IsTagNote(&cont)
	if err != nil {
		fmt.Printf("Error trying to check file header: %s\n", err.Error())
		return err
	}

	child, _ := win.GetChild()
	win.Remove(child)

	if ok {
		tagList, err := lib.TextToEnt(&cont)
		if err != nil {
			fmt.Printf("Error splitting content: '%s'", err.Error())
			return err
		}
		appShowTags(win, tagList)
	} else {
		appConsolatoryWin(win)
	}

	return nil
}

// runFileChooser creates a gtk FileChooserDialog and returns
// the path of the selected file as string.
func runFileChooser(win *gtk.Window) (string, error) {

	var fn string

	openFile, err := gtk.FileChooserDialogNewWith2Buttons("Open file", win, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Cancel", gtk.RESPONSE_CANCEL,
		"Ok", gtk.RESPONSE_OK)
	defer openFile.Destroy()
	if err != nil {
		return "", err
	}

	openFile.SetDefaultSize(50, 50)

	res := openFile.Run()

	if res == int(gtk.RESPONSE_OK) {
		fn = openFile.FileChooser.GetFilename()
	}

	return fn, nil
}

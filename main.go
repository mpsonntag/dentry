// Copyright (c) 2016, Michael Sonntag (michael.p.sonntag@gmail.com)
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of the copyright holder nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mpsonntag/dentry/lib"
)

// LogLevel is the definition of various logging levels.
type LogLevel int

const (
	// DEBUG LogLevel is logged to the Standard writer.
	DEBUG LogLevel = iota
	// INFO LogLevel is logged to the Standard writer.
	INFO
	// WARN LogLevel is logged to the Error writer.
	WARN
	// ERR LogLevel is logged to the Error writer.
	ERR
)

func (lvl LogLevel) String() string {
	switch lvl {

	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARNING"
	case ERR:
		return "ERROR"
	}

	return "UNDEFINED"
}

func log(lvl LogLevel, message string) {

	if lvl > INFO {
		fmt.Fprintf(os.Stderr, "[%s] %s\n", lvl.String(), message)
	}

	fmt.Fprintf(os.Stdout, "[%s] %s\n", lvl.String(), message)
}

func main() {
	log(DEBUG, "Init gtk")
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log(ERR, fmt.Sprintf(" creating main window: %s\n", err.Error()))
		os.Exit(-1)
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
		log(ERR, fmt.Sprintf(" creating main window content: %s\n", err.Error()))
		os.Exit(-1)
	}

	btn.Connect("clicked", handleFiles, win)
	win.Add(btn)
	win.ShowAll()

	// starts the main loop of the application, waiting for sthg to happen
	log(DEBUG, "Start main")
	gtk.Main()
	log(DEBUG, "Main started")
}

// appShowTag displays a list of tags in the main window
func appShowTags(win *gtk.Window, tagList *[]lib.TagEnt) {
	grid, err := gtk.GridNew()
	if err != nil {
		log(ERR, fmt.Sprintf(" creating grid: %s\n", err.Error()))
		return
	}
	grid.SetRowSpacing(10)
	grid.SetColumnSpacing(10)

	for i, ent := range *tagList {
		label, err := gtk.LabelNew(ent.Content)
		if err != nil {
			log(ERR, fmt.Sprintf(" creating label: %s\n", err.Error()))
			return
		}

		label.SetHAlign(gtk.ALIGN_START)
		grid.Attach(label, 0, i, 1, 1)

		hboxTag, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
		if err != nil {
			log(ERR, fmt.Sprintf(" creating hboxTag: %s", err.Error()))
			return
		}

		for _, tag := range ent.Tags {
			but, err := gtk.ButtonNew()
			if err != nil {
				log(ERR, fmt.Sprintf(" creating button: %s", err.Error()))
				return
			}
			but.SetLabel(tag)
			hboxTag.Add(but)
		}

		grid.Attach(hboxTag, 1, i, 1, 1)
	}

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log(ERR, fmt.Sprintf(" creating vbox: %s\n", err.Error()))
		return
	}
	vbox.Add(grid)

	btn, err := gtk.ButtonNewWithLabel("Try again")
	if err != nil {
		log(ERR, fmt.Sprintf(" creating file reload button: %s\n", err.Error()))
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
		log(ERR, fmt.Sprintf(" creating consolatory label: %s", err.Error()))
		return
	}

	btn, err := gtk.ButtonNewWithLabel("Try again")
	if err != nil {
		log(ERR, fmt.Sprintf(" creating file reload button: %s\n", err.Error()))
		return
	}
	btn.Connect("clicked", handleFiles, win)

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log(ERR, fmt.Sprintf(" creating vbox: %s\n", err.Error()))
		return
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
		log(ERR, fmt.Sprintf(" selecting file: %s\n", err.Error()))
		return err
	}

	cont, err := ioutil.ReadFile(fp)
	if err != nil {
		log(ERR, fmt.Sprintf(" reading file: '%s'\n", err.Error()))
		return err
	}

	ok, err := lib.IsTagNote(&cont)
	if err != nil {
		log(ERR, fmt.Sprintf(" checking file header: %s\n", err.Error()))
		return err
	}

	child, _ := win.GetChild()
	win.Remove(child)

	if ok {
		tagList, err := lib.TextToEnt(&cont)
		if err != nil {
			log(ERR, fmt.Sprintf(" splitting content: '%s'", err.Error()))
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

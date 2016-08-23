package main

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	fmt.Println("dentry started")

	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Printf("Error occurred: '%s'", err.Error())
	}
	win.SetTitle("Other Dentry")

	// required to end program properly; first string needs to be as supported signal e.g. "destroy"
	win.Connect("destroy", func(){
		gtk.MainQuit()
	})

	// required to show window
	win.ShowAll()

	// required to display window
	gtk.Main()
}

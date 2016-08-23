package main

import (
	"fmt"

	"github.com/mattn/go-gtk/gtk"
)

func main() {
	fmt.Println("dentry started")

	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)

	window.SetPosition(gtk.WIN_POS_CENTER)

	window.SetTitle("Dentry")

	// required to end program properly; first string needs to be "destroy"
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// required to show window
	window.ShowAll()
	// required to display window
	gtk.Main()
}

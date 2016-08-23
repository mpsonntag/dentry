package main

import (
	"fmt"

	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/glib"
)

func main() {
	fmt.Println("dentry started")

	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)

	window.SetPosition(gtk.WIN_POS_CENTER)

	window.SetTitle("Dentry")

	// required to end program properly; first string needs to be "destroy", other text is required but can be
	// set freely.
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		println("Close Dentry.", ctx.Data().(string))
		gtk.MainQuit()
	}, "Dentry closed")

	window.ShowAll()
	gtk.Main()
}

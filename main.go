package main

import (
	"fmt"

	"github.com/mattn/go-gtk/gtk"
)

func main() {
	fmt.Println("dentry started")

	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.ShowAll()
	gtk.Main()
}

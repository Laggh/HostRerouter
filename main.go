package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	button := widget.NewButton("Click Me", func() {
		w.SetContent(widget.NewLabel("Button Clicked!"))
	})

	// Set initial content
	w.SetContent(button)
	w.ShowAndRun()
}

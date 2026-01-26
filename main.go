package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")
	w.SetTitle("Host Rerouter")

	tamanho := 6
	table := widget.NewTable(
		//lenght
		func() (int, int) {
			return tamanho, 3
		},
		//create
		func() fyne.CanvasObject {
			c := container.NewStack()
			r := canvas.NewRectangle(color.Transparent)
			r.SetMinSize(fyne.NewSize(40, 40))
			c.Add(r)
			return c
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			row, col := i.Row, i.Col
			o.(*fyne.Container).Objects = nil
			o.(*fyne.Container).Add(widget.NewLabel(fmt.Sprintf("R %d C %d", row, col)))
			o.(*fyne.Container).Add(widget.NewLabel(fmt.Sprintf("R %d C %d", row, col)))
		},
	)

	btn := widget.NewButton("Action", func() {
		fmt.Println("Printando algo no console!")
		tamanho++
	})
	split := container.NewHSplit(table, btn)
	split.SetOffset(0.7)

	w.Resize(fyne.NewSize(600, 400))
	w.SetContent(split)
	w.ShowAndRun()
}

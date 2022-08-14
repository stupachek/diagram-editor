package main

import (
	"bytes"
	"log"
	"sem4/figure"
	"sem4/parser"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	var img *canvas.Image = &canvas.Image{}
	entry := widget.NewMultiLineEntry()

	myApp := app.New()
	myApp.Settings().SetTheme(theme.LightTheme())
	w := myApp.NewWindow("Diagram editor")
	w.Resize(fyne.NewSize(400, 400))
	entry.OnChanged = func(s string) {
		r, _ := parser.Parse(s)
		result := figure.DrawBlock(r)
		img = canvas.NewImageFromReader(bytes.NewReader(result), "re.png")
		img.FillMode = canvas.ImageFillContain
		content := container.New(layout.NewMaxLayout(), img)
		row := container.NewGridWithColumns(2, entry, content)
		w.SetContent(row)
	}
	content := container.New(layout.NewCenterLayout(), img)
	row := container.NewMax(entry, content)
	w.SetContent(row)
	openItem := fyne.NewMenuItem("Open", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			buf := bytes.NewBuffer(nil)
			_, err = buf.ReadFrom(reader)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			entry.SetText(buf.String())
		}, w)
		fd.Show()
	})
	saveItem := fyne.NewMenuItem("Save", func() {
		fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if writer == nil {
				log.Println("Cancelled")
				return
			}
			_, err = writer.Write([]byte(entry.Text))
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
		}, w)
		fd.Show()
	})
	savePNGItem := fyne.NewMenuItem("Save .png", func() {
		fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if writer == nil {
				log.Println("Cancelled")
				return
			}
			_, err = writer.Write(img.Resource.Content())
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
		}, w)
		fd.Show()
	})
	fileItem := fyne.NewMenu("File", openItem, saveItem, savePNGItem)
	menu := fyne.NewMainMenu(fileItem)
	w.SetMainMenu(menu)

	w.ShowAndRun()
}

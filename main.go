package main

import (
	"fmt"
	"io"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const NEW_FILE_NAME = "Untitled"

func main() {
	a := app.New()
	w := a.NewWindow("Notepad")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(800, 500))

	// Running
	tabs := AddTabs(w)
	entries := AddFirstTab(tabs)
	AddProgramMenu(w, tabs, entries)

	w.Show()
	a.Run()
}

func NewTextEntry() *widget.Entry {
	textInput := widget.NewMultiLineEntry()
	textInput.SetPlaceHolder("Enter text here...")

	return textInput
}

func AddTabs(w fyne.Window) *container.AppTabs {
	tabs := container.NewAppTabs()
	w.SetContent(tabs)

	return tabs
}

func AddFirstTab(tabs *container.AppTabs) []*widget.Entry {
	entries := []*widget.Entry{}
	firstEntry := NewTextEntry()
	tabs.Append(container.NewTabItem(NEW_FILE_NAME, firstEntry))

	return append(entries, firstEntry)
}

func AddProgramMenu(w fyne.Window, tabs *container.AppTabs, entries []*widget.Entry) {
	new := fyne.NewMenuItem("New", func() {
		entry := NewTextEntry()
		entries = append(entries, entry)
		tabs.Append(container.NewTabItem(NEW_FILE_NAME, entry))
	})

	save := fyne.NewMenuItem("Save", func() {
		tabText := tabs.Selected().Text

		if tabText != NEW_FILE_NAME {
			file, err := os.Create(tabText)

			if err != nil {
				fmt.Println("Error saving file: ", err)
			}

			defer file.Close()

			file.WriteString(entries[tabs.SelectedIndex()].Text)
		} else {
			SaveFile(w, tabs, entries)
		}
	})

	saveAs := fyne.NewMenuItem("Save As...", func() {
		SaveFile(w, tabs, entries)
	})

	open := fyne.NewMenuItem("Open", func() {

	})

	menu := fyne.NewMenu("File", new, save, saveAs, open)
	w.SetMainMenu(fyne.NewMainMenu(menu))
}

func SaveFile(w fyne.Window, tabs *container.AppTabs, entries []*widget.Entry) {
	dialog.ShowFileSave(
		func(uc fyne.URIWriteCloser, err error) {
			if uc == nil {
				return
			}
			
			io.WriteString(uc, entries[tabs.SelectedIndex()].Text)
			tabs.Selected().Text = uc.URI().Path() // change filename in the tab with saved filename
			tabs.Refresh()
		},
		w,
	)
}

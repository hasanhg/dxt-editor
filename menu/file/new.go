package file

import "fyne.io/fyne/v2"

func NewMenuItem(w fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem("New", func() {})
}

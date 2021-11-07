package file

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

func OpenBMPMenuItem(w fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem("Open BMP", func() {
		fd := dialog.NewFileOpen(func(uc fyne.URIReadCloser, e error) {}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".bmp"}))
		fd.Show()
	})
}

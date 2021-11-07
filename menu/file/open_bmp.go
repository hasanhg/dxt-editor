package file

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
)

func OpenBMPMenuItem(w fyne.Window) *fyne.MenuItem {
	ctrlB := &desktop.CustomShortcut{KeyName: fyne.KeyB, Modifier: desktop.ControlModifier}
	w.Canvas().AddShortcut(ctrlB, func(shortcut fyne.Shortcut) {
		openBMP(w)
	})

	return fyne.NewMenuItem("Open BMP\t(CTRL + B)", func() {
		openBMP(w)
	})
}

func openBMP(w fyne.Window) {
	fd := dialog.NewFileOpen(func(uc fyne.URIReadCloser, e error) {}, w)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".bmp"}))
	fd.Show()
}

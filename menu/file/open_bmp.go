package file

import (
	xdOpen "dxt-editor/dialog/open"
	"dxt-editor/global"

	"fyne.io/fyne/v2"
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
	if global.Dialog != nil {
		return
	}
	fd := xdOpen.NewFileOpen("BMP", func(uc fyne.URIReadCloser, e error) {
		defer func() { global.Dialog = nil }()
	}, w)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".bmp"}))

	fd.Show()
	global.Dialog = fd
}

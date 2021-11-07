package tools

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"

	xdSave "dxt-editor/dialog/save"
	"dxt-editor/global"
)

func CropMenuItem(w fyne.Window) *fyne.MenuItem {
	crop := fyne.NewMenuItem("Crop", func() {
		fd := xdSave.NewFileSave(func(uc fyne.URIWriteCloser, e error) {
			if uc == nil {
				return
			}

		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".bmp"}))
		fd.Show()
	})
	global.ActivateAfterOpen = append(global.ActivateAfterOpen, crop)
	return crop
}

package file

import (
	"dxt-editor/dxt"
	"fmt"
	"io/ioutil"

	xdOpen "dxt-editor/dialog/open"
	"dxt-editor/global"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
)

func OpenDXTMenuItem(w fyne.Window) *fyne.MenuItem {
	ctrlD := &desktop.CustomShortcut{KeyName: fyne.KeyD, Modifier: desktop.ControlModifier}
	w.Canvas().AddShortcut(ctrlD, func(shortcut fyne.Shortcut) {
		openDXT(w)
	})

	return fyne.NewMenuItem("Open DXT\t(CTRL + D)", func() {
		openDXT(w)
	})
}

func openDXT(w fyne.Window) {
	if global.Dialog != nil {
		return
	}

	fd := xdOpen.NewFileOpen("DXT", func(uc fyne.URIReadCloser, e error) {
		defer func() { global.Dialog = nil }()
		if uc == nil {
			return
		}

		data, err := ioutil.ReadAll(uc)
		if err != nil {
			dialog.NewError(err, w).Show()
			return
		}

		w.SetTitle(fmt.Sprintf("DXT Editor - %s", uc.URI().Path()))
		global.DXTFile = dxt.NewBuffer(data).ParseDXT(uc.URI().Path())

		imgCanvas := canvas.NewImageFromImage(global.DXTFile.Image)
		imgCanvas.FillMode = canvas.ImageFillOriginal

		//w := fyne.CurrentApp().NewWindow(uc.URI().Name())
		center := container.New(layout.NewCenterLayout(),
			imgCanvas,
		)

		for _, item := range global.ActivateAfterOpen {
			item.Disabled = false
		}

		size := global.DXTFile.Image.Rect.Size()
		center.Resize(fyne.NewSize(float32(size.X), float32(size.Y)))
		w.SetContent(center)
		w.Resize(center.MinSize())
		w.Canvas().Refresh(center)
		w.Content().Refresh()
	}, w)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".dxt"}))

	fd.Show()
	global.Dialog = fd
}

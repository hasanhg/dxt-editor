package file

import (
	"dxt-editor/dxt"
	"image"
	"io/ioutil"
	"os"

	xdOpen "dxt-editor/dialog/open"
	"dxt-editor/global"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"golang.org/x/image/bmp"
)

func OpenDXTMenuItem(w fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem("Open DXT", func() {
		fd := xdOpen.NewFileOpen(func(uc fyne.URIReadCloser, e error) {
			if uc == nil {
				return
			}

			data, err := ioutil.ReadAll(uc)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}

			global.DXTFile = dxt.NewBuffer(data).ParseDXT()

			subimg := global.DXTFile.Image.SubImage(image.Rect(65, 65*6, 65*2, 65*7))

			f, err := os.Create("test.bmp")
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			err = bmp.Encode(f, subimg)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}

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
	})
}

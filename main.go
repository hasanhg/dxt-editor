package main

import (
	"image"
	"image/color"
	"io/ioutil"
	"os"

	xdOpen "dxt-editor/dialog/open"
	xdSave "dxt-editor/dialog/save"
	"dxt-editor/dxt"
	"dxt-editor/uif"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/bmp"
)

var (
	dxtFile *dxt.DXTFile
)

func main() {
	a := app.New()
	w := a.NewWindow("DXT Editor")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(512, 512))

	imgCanvas := canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 0, 0)))
	imgCanvas.FillMode = canvas.ImageFillOriginal

	activateAfterOpen := []*fyne.MenuItem{}

	crop := fyne.NewMenuItem("Crop", func() {
		fd := xdSave.NewFileSave(func(uc fyne.URIWriteCloser, e error) {
			if uc == nil {
				return
			}

		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".bmp"}))
		fd.Show()
	})
	//crop.Disabled = true
	activateAfterOpen = append(activateAfterOpen, crop)

	w.SetMainMenu(
		fyne.NewMainMenu(
			fyne.NewMenu("File",
				fyne.NewMenuItem("New", func() {

				}),
				fyne.NewMenuItem("Open DXT", func() {
					fd := xdOpen.NewFileOpen(func(uc fyne.URIReadCloser, e error) {
						if uc == nil {
							return
						}

						data, err := ioutil.ReadAll(uc)
						if err != nil {
							dialog.NewError(err, w).Show()
							return
						}

						dxtFile = dxt.NewBuffer(data).ParseDXT()

						subimg := dxtFile.Image.SubImage(image.Rect(65, 65*6, 65*2, 65*7))

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

						imgCanvas := canvas.NewImageFromImage(dxtFile.Image)
						imgCanvas.FillMode = canvas.ImageFillOriginal

						//w := fyne.CurrentApp().NewWindow(uc.URI().Name())
						center := container.New(layout.NewCenterLayout(),
							imgCanvas,
						)

						for _, item := range activateAfterOpen {
							item.Disabled = false
						}

						size := dxtFile.Image.Rect.Size()
						center.Resize(fyne.NewSize(float32(size.X), float32(size.Y)))
						w.SetContent(center)
						w.Resize(center.MinSize())
						w.Canvas().Refresh(center)
						w.Content().Refresh()
					}, w)
					fd.SetFilter(storage.NewExtensionFileFilter([]string{".dxt"}))
					fd.Show()
				}),
				fyne.NewMenuItem("Open UIF", func() {
					fd := xdOpen.NewFileOpen(func(uc fyne.URIReadCloser, e error) {
						if uc == nil {
							return
						}

						data, err := ioutil.ReadAll(uc)
						if err != nil {
							dialog.NewError(err, w).Show()
							return
						}

						uifFile := uif.NewBuffer(data).ParseUIF()
						split := container.NewHSplit(
							container.New(layout.NewCenterLayout(), drawObject(uifFile.Root)...),
							drawHierarchy(uifFile),
						)
						split.Offset = 0.5

						w.SetContent(split)
						w.Resize(split.MinSize())
						w.Canvas().Refresh(split)
						w.Content().Refresh()

					}, w)
					fd.SetFilter(storage.NewExtensionFileFilter([]string{".uif"}))
					fd.Show()
				}),
				fyne.NewMenuItem("Open BMP", func() {
					fd := dialog.NewFileOpen(func(uc fyne.URIReadCloser, e error) {}, w)
					fd.SetFilter(storage.NewExtensionFileFilter([]string{".bmp"}))
					fd.Show()
				}),
			),
			fyne.NewMenu("Tools",
				crop,
			),
		),
	)

	w.ShowAndRun()
}

func getImg() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 512, 512))
	for y := 0; y < 512; y++ {
		for x := 0; x < 512; x++ {
			d := uint8((x + y) / 4)
			img.SetRGBA(x, y, color.RGBA{A: 255, R: d, G: d, B: d})
		}
	}
	return img
}

func drawObject(obj *uif.Object) []fyne.CanvasObject {
	arr := []fyne.CanvasObject{}
	for _, ch := range obj.Children {
		arr = append(arr, drawObject(ch)...)
	}
	switch obj.Type {
	case uif.OT_BASE:
	case uif.OT_IMAGE:
		data, _ := ioutil.ReadFile("C:\\Users\\gurso\\Desktop\\PS\\ui\\Icon_Acc.dxt")
		dxtFile := dxt.NewBuffer(data).ParseDXT()
		imgCanvas := canvas.NewImageFromImage(dxtFile.Image)
		imgCanvas.FillMode = canvas.ImageFillOriginal
		arr = append(arr,
			container.New(layout.NewCenterLayout(),
				imgCanvas,
			),
		)
	}

	return arr
}

func drawHierarchy(uifFile *uif.UIFFile) fyne.CanvasObject {
	data := map[string][]string{}
	for _, child := range uifFile.Root.Children {
		getChildren(&data, uifFile.Root, child)
	}

	data[""] = []string{uifFile.Root.ID}
	tree := widget.NewTreeWithStrings(data)
	return container.NewBorder(nil, nil, nil, nil, tree)
}

func getChildren(data *map[string][]string, parent, child *uif.Object) {
	for _, ch := range child.Children {
		getChildren(data, child, ch)
	}
	(*data)[parent.ID] = append((*data)[parent.ID], child.ID)
}

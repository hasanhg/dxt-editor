package file

import (
	xdOpen "dxt-editor/dialog/open"
	"dxt-editor/dxt"
	"dxt-editor/uif"
	"io/ioutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func OpenUIFMenuItem(w fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem("Open UIF", func() {
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
	})
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

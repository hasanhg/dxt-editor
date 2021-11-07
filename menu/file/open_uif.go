package file

import (
	xdOpen "dxt-editor/dialog/open"
	"dxt-editor/dxt"
	"dxt-editor/uif"
	"fmt"
	"image"
	"io/ioutil"
	"log"

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
			w.Resize(split.MinSize().Max(fyne.NewPos(512, 512)))
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

		x0 := obj.Crop.Min.X * float32(dxtFile.Image.Rect.Size().X)
		x1 := obj.Crop.Max.X * float32(dxtFile.Image.Rect.Size().X)
		y0 := obj.Crop.Min.Y * float32(dxtFile.Image.Rect.Size().Y)
		y1 := obj.Crop.Max.Y * float32(dxtFile.Image.Rect.Size().Y)

		img := dxtFile.Image.SubImage(image.Rect(int(x0), int(y0), int(x1), int(y1)))
		imgCanvas := canvas.NewImageFromImage(img)
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
	objects := map[string]*uif.Object{}
	for _, child := range uifFile.Root.Children {
		getChildren(&data, &objects, uifFile.Root, child)
	}

	data[""] = []string{uifFile.Root.ID}
	objects[uifFile.Root.ID] = uifFile.Root
	tree := widget.NewTreeWithStrings(data)
	tree.OpenAllBranches()
	tree.OnSelected = func(uid widget.TreeNodeID) {
		obj, ok := objects[uid]
		if !ok {
			log.Println("Object not found!")
			return
		}

		fmt.Printf("%+v\n", obj)
	}

	return container.NewBorder(nil, nil, nil, nil, tree)
}

func getChildren(data *map[string][]string, objects *map[string]*uif.Object, parent, child *uif.Object) {
	for _, ch := range child.Children {
		getChildren(data, objects, child, ch)
	}
	(*data)[parent.ID] = append((*data)[parent.ID], child.ID)
	(*objects)[child.ID] = child
}

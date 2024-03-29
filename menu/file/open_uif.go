package file

import (
	xdOpen "dxt-editor/dialog/open"
	"dxt-editor/dxt"
	"dxt-editor/global"
	"dxt-editor/uif"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var (
	selectedObject  *uif.Object
	propertiesTree  *widget.Tree
	selectedPropUID string
)

func OpenUIFMenuItem(w fyne.Window) *fyne.MenuItem {
	ctrlU := &desktop.CustomShortcut{KeyName: fyne.KeyU, Modifier: desktop.ControlModifier}
	w.Canvas().AddShortcut(ctrlU, func(shortcut fyne.Shortcut) {
		openUIF(w)
	})

	return fyne.NewMenuItem("Open UIF\t(CTRL + U)", func() {
		openUIF(w)
	})
}

func openUIF(w fyne.Window) {
	if global.Dialog != nil {
		return
	}

	fd := xdOpen.NewFileOpen("UIF", func(uc fyne.URIReadCloser, e error) {
		defer func() { global.Dialog = nil }()
		if uc == nil {
			return
		}

		global.CurrentDirectory = filepath.Dir(uc.URI().Path())
		data, err := ioutil.ReadAll(uc)
		if err != nil {
			dialog.NewError(err, w).Show()
			return
		}

		w.SetTitle(fmt.Sprintf("UIF Editor - %s", uc.URI().Path()))
		global.UIFFile = uif.NewBuffer(data).ParseUIF()
		refreshUIFEditor(w, uc)
	}, w)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".uif"}))
	fd.Show()
	global.Dialog = fd
}

func refreshUIFEditor(w fyne.Window, uc fyne.URIReadCloser) {
	split := container.NewHSplit(
		container.New(layout.NewCenterLayout(), drawObject(uc, global.UIFFile.Root)...),
		drawHierarchy(w, uc, global.UIFFile),
	)
	split.Offset = 0.5
	w.SetContent(split)
	//w.Resize(split.MinSize().Max(fyne.NewPos(512, 512)))
	w.Canvas().Refresh(split)
	w.Content().Refresh()
}

func drawObject(uc fyne.URIReadCloser, obj *uif.Object) []fyne.CanvasObject {
	arr := []fyne.CanvasObject{}
	if !obj.Visible {
		return arr
	}

	for _, ch := range obj.Children {
		arr = append(arr, drawObject(uc, ch)...)
	}
	switch obj.Type {
	case uif.OT_BASE:
	case uif.OT_IMAGE:
		p := filepath.Dir(uc.URI().Path())
		p = filepath.Join(p, obj.Texture)
		data, err := ioutil.ReadFile(p)
		if err != nil {
			p = filepath.Dir(filepath.Dir(uc.URI().Path()))
			p = filepath.Join(p, obj.Texture)
			data, err = ioutil.ReadFile(p)
			if err != nil {
				return nil
			}
		}

		dxtFile := dxt.NewBuffer(data).ParseDXT(uc.URI().Path())

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

func drawHierarchy(w fyne.Window, uc fyne.URIReadCloser, uifFile *uif.UIFFile) fyne.CanvasObject {
	data := map[string][]string{}
	objects := map[string]*uif.Object{}
	for _, child := range uifFile.Root.Children {
		getChildren(&data, &objects, uifFile.Root, child)
	}

	data[""] = []string{uifFile.Root.ID}
	objects[uifFile.Root.ID] = uifFile.Root
	propertiesTree = widget.NewTreeWithStrings(data)
	propertiesTree.OpenAllBranches()

	props := container.NewBorder(nil, nil, nil, nil)
	props.Hide()

	propertiesTree.OnSelected = func(uid widget.TreeNodeID) {
		obj, ok := objects[uid]
		if !ok {
			log.Println("Object not found!")
			return
		}

		selectedPropUID = uid
		drawProperties(w, uc, obj, props)
		w.Canvas().Refresh(props)
		fmt.Printf("%+v\n", obj)
	}

	return container.NewBorder(nil, nil, nil, nil, container.NewVSplit(propertiesTree, container.New(layout.NewAdaptiveGridLayout(1), props)))
}

func getChildren(data *map[string][]string, objects *map[string]*uif.Object, parent, child *uif.Object) {
	for _, ch := range child.Children {
		getChildren(data, objects, child, ch)
	}
	(*data)[parent.ID] = append((*data)[parent.ID], child.ID)
	(*objects)[child.ID] = child
}

func drawProperties(w fyne.Window, uc fyne.URIReadCloser, obj *uif.Object, props *fyne.Container) {
	props.Objects = []fyne.CanvasObject{}

	idWidget := widget.NewEntry()
	idWidget.SetText(obj.ID)

	nameWidget := widget.NewEntry()
	nameWidget.SetText(obj.Name)

	typeWidget := widget.NewEntry()
	typeWidget.SetText(obj.Type.String())

	xRectWidget := widget.NewEntry()
	xRectWidget.SetText(fmt.Sprintf("%d", obj.Rect.Min.X))

	yRectWidget := widget.NewEntry()
	yRectWidget.SetText(fmt.Sprintf("%d", obj.Rect.Min.Y))

	widthRectWidget := widget.NewEntry()
	widthRectWidget.SetText(fmt.Sprintf("%d", obj.Rect.Max.X-obj.Rect.Min.X))

	heightRectWidget := widget.NewEntry()
	heightRectWidget.SetText(fmt.Sprintf("%d", obj.Rect.Max.Y-obj.Rect.Min.Y))

	rectAcc := widget.NewAccordion(
		widget.NewAccordionItem("[Object]",
			widget.NewForm(
				widget.NewFormItem("X", xRectWidget),
				widget.NewFormItem("Y", yRectWidget),
				widget.NewFormItem("Width", widthRectWidget),
				widget.NewFormItem("Height", heightRectWidget),
			),
		),
	)
	rectAcc.MultiOpen = true
	//rectAcc.OpenAll()

	tooltipWidget := widget.NewEntry()
	tooltipWidget.SetText(obj.Tooltip)

	openSoundWidget := widget.NewEntry()
	openSoundWidget.SetText(obj.SoundOpen)

	closeSounWidget := widget.NewEntry()
	closeSounWidget.SetText(obj.SoundClose)

	visiblityWidget := widget.NewCheck("", nil)
	visiblityWidget.SetChecked(obj.Visible)
	visiblityWidget.OnChanged = func(b bool) {
		obj.Visible = b
		refreshUIFEditor(w, uc)
		propertiesTree.Select(selectedPropUID)
	}

	_props := widget.NewAccordion(
		widget.NewAccordionItem("Common",
			widget.NewForm(
				widget.NewFormItem("ID", idWidget),
				widget.NewFormItem("Name", nameWidget),
				widget.NewFormItem("Type", typeWidget),
				widget.NewFormItem("Visible", visiblityWidget),
				widget.NewFormItem("Rectangle", rectAcc),
				widget.NewFormItem("Tooltip", tooltipWidget),
				widget.NewFormItem("Open Sound", openSoundWidget),
				widget.NewFormItem("Close Sound", closeSounWidget),
			),
		),
	)
	_props.MultiOpen = true

	if obj.Type.String() != "" {
		var acItem fyne.CanvasObject

		switch obj.Type {
		case uif.OT_IMAGE:
			dxtWidget := widget.NewEntry()
			dxtWidget.SetText(obj.Texture)

			x0CropWidget := widget.NewEntry()
			x0CropWidget.SetText(fmt.Sprintf("%f", obj.Crop.Min.X))

			y0CropWidget := widget.NewEntry()
			y0CropWidget.SetText(fmt.Sprintf("%f", obj.Crop.Min.Y))

			x1CropWidget := widget.NewEntry()
			x1CropWidget.SetText(fmt.Sprintf("%f", obj.Crop.Max.X))

			y1CropWidget := widget.NewEntry()
			y1CropWidget.SetText(fmt.Sprintf("%f", obj.Crop.Max.Y))

			cropAcc := widget.NewAccordion(
				widget.NewAccordionItem("[Object]",
					widget.NewForm(
						widget.NewFormItem("X0", x0CropWidget),
						widget.NewFormItem("Y0", y0CropWidget),
						widget.NewFormItem("X1", x1CropWidget),
						widget.NewFormItem("Y1", y1CropWidget),
					),
				),
			)
			cropAcc.MultiOpen = true
			//cropAcc.OpenAll()

			animFrameWidget := widget.NewEntry()
			animFrameWidget.SetText(fmt.Sprintf("%f", obj.AnimationFrame))

			acItem = widget.NewForm(
				widget.NewFormItem("Texture", dxtWidget),
				widget.NewFormItem("Crop", cropAcc),
				widget.NewFormItem("Anim. Rate", animFrameWidget),
			)
		}

		if acItem != nil {
			_props.Items = append(_props.Items, widget.NewAccordionItem(obj.Type.String(), acItem))
		}
	}

	_props.OpenAll()

	props.Add(container.NewVScroll(
		container.NewVBox(_props),
	))
	props.Show()
}

package file

import (
	"dxt-editor/dxt"
	"dxt-editor/global"
	"dxt-editor/uif"
	"dxt-editor/utils"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func ExportPNGMenuItem(w fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem("Export PNG", func() {
		fd := dialog.NewFolderOpen(func(lu fyne.ListableURI, e error) {
			if lu == nil {
				return
			}
			export(lu)
		}, w)
		//fd.SetFilter(storage.NewExtensionFileFilter([]string{".png"}))
		fd.Show()
	})
}

func export(lu fyne.ListableURI) {
	if global.DXTFile != nil {
		p := lu.Path()
		fileName := fmt.Sprintf("%s.png", utils.FileNameWithoutExtension(filepath.Base(global.DXTFile.Path)))
		fullPath := filepath.Join(p, fileName)

		file, err := os.Create(fullPath)
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		err = png.Encode(file, global.DXTFile.Image)
		if err != nil {
			log.Println(err)
			return
		}

		return
	}

	if global.UIFFile != nil {
		toExport := []*uif.Object{}
		exportUIF(global.UIFFile.Root, &toExport)
		for _, obj := range toExport {
			p := lu.Path()
			fileName := fmt.Sprintf("%s_%s.png", utils.FileNameWithoutExtension(filepath.Base(obj.Texture)), obj.ID)
			fullPath := filepath.Join(p, fileName)

			dxtPath := filepath.Join(global.CurrentDirectory, obj.Texture)
			data, err := ioutil.ReadFile(dxtPath)
			if err != nil {
				log.Println(err)
				return
			}

			file, err := os.Create(fullPath)
			if err != nil {
				log.Println(err)
				return
			}
			defer file.Close()

			dxtFile := dxt.NewBuffer(data).ParseDXT(dxtPath)
			x0 := obj.Crop.Min.X * float32(dxtFile.Image.Rect.Size().X)
			x1 := obj.Crop.Max.X * float32(dxtFile.Image.Rect.Size().X)
			y0 := obj.Crop.Min.Y * float32(dxtFile.Image.Rect.Size().Y)
			y1 := obj.Crop.Max.Y * float32(dxtFile.Image.Rect.Size().Y)

			img := dxtFile.Image.SubImage(image.Rect(int(x0), int(y0), int(x1), int(y1)))
			err = png.Encode(file, img)
			if err != nil {
				log.Println(err)
				return
			}
		}

		return
	}
}

func exportUIF(obj *uif.Object, toExport *[]*uif.Object) {

	for _, ch := range obj.Children {
		exportUIF(ch, toExport)
	}

	if obj.Texture == "" {
		return
	}

	*toExport = append(*toExport, obj)
}

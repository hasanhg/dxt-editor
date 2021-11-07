package file

import (
	"dxt-editor/dxt"
	"dxt-editor/global"
	"dxt-editor/uif"
	"dxt-editor/utils"
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

		return
	}

	if global.UIFFile != nil {
		toExport := &uif.Object{}
		_, texture := exportUIF(global.UIFFile.Root, &toExport, 0, "")
		p := lu.Path()
		fileName := utils.FileNameWithoutExtension(filepath.Base(texture)) + ".png"
		fullPath := filepath.Join(p, fileName)

		data, err := ioutil.ReadFile(filepath.Join(global.CurrentDirectory, texture))
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

		dxtFile := dxt.NewBuffer(data).ParseDXT()
		x0 := toExport.Crop.Min.X * float32(dxtFile.Image.Rect.Size().X)
		x1 := toExport.Crop.Max.X * float32(dxtFile.Image.Rect.Size().X)
		y0 := toExport.Crop.Min.Y * float32(dxtFile.Image.Rect.Size().Y)
		y1 := toExport.Crop.Max.Y * float32(dxtFile.Image.Rect.Size().Y)

		img := dxtFile.Image.SubImage(image.Rect(int(x0), int(y0), int(x1), int(y1)))
		err = png.Encode(file, img)
		if err != nil {
			log.Println(err)
			return
		}

		return
	}
}

func exportUIF(obj *uif.Object, toExport **uif.Object, resolution int, texture string) (int, string) {

	for _, ch := range obj.Children {
		resolution, texture = exportUIF(ch, toExport, resolution, texture)
	}

	if obj.Texture == "" {
		return resolution, texture
	}

	width := obj.Rect.Max.X - obj.Rect.Min.X
	height := obj.Rect.Max.Y - obj.Rect.Min.Y
	if width*height > resolution {
		resolution = width * height
		texture = obj.Texture
		*toExport = obj
	}

	return resolution, texture
}

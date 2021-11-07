package tools

import "fyne.io/fyne/v2"

func ToolsMenu(w fyne.Window) *fyne.Menu {
	return fyne.NewMenu("Tools",
		CropMenuItem(w),
	)
}

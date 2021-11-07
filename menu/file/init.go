package file

import "fyne.io/fyne/v2"

func FileMenu(w fyne.Window) *fyne.Menu {
	return fyne.NewMenu("File",
		NewMenuItem(w),
	)
}

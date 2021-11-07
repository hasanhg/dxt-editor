package file

import "fyne.io/fyne/v2"

func ExportMenuItem(w fyne.Window) *fyne.MenuItem {
	openMenu := fyne.NewMenuItem("Export", nil)
	openMenu.ChildMenu = fyne.NewMenu("",
		ExportPNGMenuItem(w),
	)
	return openMenu
}

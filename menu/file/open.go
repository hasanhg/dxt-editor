package file

import "fyne.io/fyne/v2"

func OpenMenuItem(w fyne.Window) *fyne.MenuItem {
	openMenu := fyne.NewMenuItem("Open", nil)
	openMenu.ChildMenu = fyne.NewMenu("",
		OpenDXTMenuItem(w),
		OpenUIFMenuItem(w),
		OpenBMPMenuItem(w),
	)
	return openMenu
}

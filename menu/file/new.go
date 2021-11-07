package file

import (
	"dxt-editor/global"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func NewMenuItem(w fyne.Window) *fyne.MenuItem {
	ctrlN := &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: desktop.ControlModifier}
	w.Canvas().AddShortcut(ctrlN, func(shortcut fyne.Shortcut) {
		onNew(w)
	})

	return fyne.NewMenuItem("New\t(Ctrl + N)", func() {
		onNew(w)
	})
}

func onNew(w fyne.Window) {
	if global.Dialog != nil {
		return
	}
	global.DXTFile = nil
	global.UIFFile = nil
	w.SetTitle(global.CustomTitle)
	w.SetContent(widget.NewLabel(""))
	w.Content().Refresh()
}

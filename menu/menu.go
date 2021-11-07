package menu

import (
	"dxt-editor/menu/file"

	"fyne.io/fyne/v2"
)

func MainMenu(w fyne.Window) *fyne.MainMenu {
	return fyne.NewMainMenu(
		file.FileMenu(w),
	)
}

package menu

import (
	"dxt-editor/menu/file"
	"dxt-editor/menu/tools"

	"fyne.io/fyne/v2"
)

func MainMenu(w fyne.Window) *fyne.MainMenu {
	return fyne.NewMainMenu(
		file.FileMenu(w),
		tools.ToolsMenu(w),
	)
}

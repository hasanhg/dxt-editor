package main

import (
	"dxt-editor/global"
	"dxt-editor/menu"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow(global.CustomTitle)
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(512, 512))
	w.SetMainMenu(menu.MainMenu(w))
	w.ShowAndRun()
}

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	myApp := app.NewWithID("net.postcert.gopom")
	myWindow := myApp.NewWindow("Pomodoro Timer")

	pomodoroTimer := createPomodoroTimer(myApp)

	tabs := container.NewAppTabs(
		container.NewTabItem("Timer", pomodoroTab(pomodoroTimer)),
		container.NewTabItem("Settings", settingsTab(pomodoroTimer)),
	)

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(300, 200))
	myWindow.ShowAndRun()
}

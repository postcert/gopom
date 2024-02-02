package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	myApp := app.NewWithID("net.postcert.gopom")
	myWindow := myApp.NewWindow("Pomodoro Timer")

	pomoconfig := createPomoConfig(myApp)

	timerTab := container.NewTabItem("Timer", pomodoroTab(pomoconfig))
	settingsTab := container.NewTabItem("Settings", settingsTab(pomoconfig))
	tabs := container.NewAppTabs(
		timerTab,
		settingsTab,
	)
	tabs.Select(settingsTab)

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(300, 200))
	myWindow.ShowAndRun()
}

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/sanity-io/litter"
)

func main() {
	myApp := app.NewWithID("net.postcert.gopom")
	myWindow := myApp.NewWindow("Pomodoro Timer")

	pomoconfig := createPomoConfig(myApp)
	litter.Dump(pomoconfig)

	pomodoroTimer := createPomodoroTimer(pomoconfig)
	litter.Dump(pomodoroTimer)

	timerTab := container.NewTabItem("Timer", pomodoroTab(pomodoroTimer))
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

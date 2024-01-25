package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func settingsTab(pomodoroTimer *PomodoroTimer) *fyne.Container {
	settingsForm := widget.NewForm(
		widget.NewFormItem("Work Duration", widget.NewEntryWithData(pomodoroTimer.workDurationBinding)),
		widget.NewFormItem("Break Duration", widget.NewEntryWithData(pomodoroTimer.breakDurationBinding)),
		widget.NewFormItem("Work Iterations", widget.NewEntryWithData(pomodoroTimer.workIterationsBinding)),
	)

	return container.NewVBox(
		settingsForm,
	)
}

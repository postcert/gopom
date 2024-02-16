package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TODO: Use Binding listeners to update this page
func pomodoroTab(pomoconfig *PomoConfig) *fyne.Container {
	pomoTimer := pomoconfig.Timer

	timerLabel := widget.NewLabelWithData(pomoTimer.DurationString)
	progressBar := widget.NewProgressBarWithData(pomoTimer.ProgressFloat)
	startButton := widget.NewButton("Start", nil)
	stopButton := widget.NewButton("Stop", nil)
	resetButton := widget.NewButton("Reset", nil)

	startButton.OnTapped = func() {
		pomoTimer.Start()
	}

	stopButton.OnTapped = func() {
		pomoTimer.Stop()
	}

	resetButton.OnTapped = func() {
		pomoTimer.Stop()
		pomoTimer.Reset()
	}

	buttons := container.NewHBox(startButton, stopButton, resetButton)
	content := container.NewVBox(timerLabel, progressBar, buttons)
	return content
}

package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const pomodoroDuration = 25 * time.Minute

func pomodoroTab(pomoconfig *PomoConfig) *fyne.Container {
	timerLabel := widget.NewLabel(fmt.Sprintf("%s:00", pomoconfig.Timer.workDuration))
	progressBar := widget.NewProgressBar()
	startButton := widget.NewButton("Start", nil)
	stopButton := widget.NewButton("Stop", nil)
	resetButton := widget.NewButton("Reset", nil)

	var timer *time.Timer
	var ticker *time.Ticker

	startButton.OnTapped = func() {
		remaining := pomodoroDuration
		timerLabel.SetText(formatDuration(remaining))
		progressBar.SetValue(0)

		timer = time.NewTimer(pomodoroDuration)
		ticker = time.NewTicker(time.Second)

		go func() {
			for {
				select {
				case <-ticker.C:
					remaining -= time.Second
					timerLabel.SetText(formatDuration(remaining))
					progress := 1.0 - (remaining.Seconds() / pomodoroDuration.Seconds())
					progressBar.SetValue(progress)
				case <-timer.C:
					timerLabel.SetText("Time's up!")
					progressBar.SetValue(1)
					ticker.Stop()
					return
				}
			}
		}()
	}

	stopButton.OnTapped = func() {
		if timer != nil {
			timer.Stop()
		}
		if ticker != nil {
			ticker.Stop()
		}
	}

	resetButton.OnTapped = func() {
		if timer != nil {
			timer.Stop()
		}
		if ticker != nil {
			ticker.Stop()
		}
		timerLabel.SetText("25:00")
		progressBar.SetValue(0)
	}
	buttons := container.NewHBox(startButton, stopButton, resetButton)
	content := container.NewVBox(timerLabel, progressBar, buttons)
	return content
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

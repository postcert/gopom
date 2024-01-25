package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const pomodoroDuration = 25 * time.Minute

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Pomodoro Timer")

	timerLabel := widget.NewLabel("25:00")
	progress := widget.NewProgressBar()
	startButton := widget.NewButton("Start", nil)
	stopButton := widget.NewButton("Stop", nil)
	resetButton := widget.NewButton("Reset", nil)

	var timer *time.Timer
	var ticker *time.Ticker

	startButton.OnTapped = func() {
		remaining := pomodoroDuration
		timerLabel.SetText(formatDuration(remaining))
		circularTimer.SetProgress(0)

		timer = time.NewTimer(pomodoroDuration)
		ticker = time.NewTicker(time.Second)

		go func() {
			for {
				select {
				case <-ticker.C:
					remaining -= time.Second
					timerLabel.SetText(formatDuration(remaining))
					progress := 1 - remaining.Seconds()/pomodoroDuration.Seconds()
					circularTimer.SetProgress(progress)
				case <-timer.C:
					timerLabel.SetText("Time's up!")
					circularTimer.SetProgress(1)
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
		circularTimer.SetProgress(0)
	}

	buttons := container.NewHBox(startButton, stopButton, resetButton)
	content := container.NewVBox(timerLabel, circularTimer, buttons)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 300))
	myWindow.ShowAndRun()
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

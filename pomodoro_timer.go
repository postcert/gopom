package main

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type PomodoroTimer struct {
	workDurationBinding   binding.String
	breakDurationBinding  binding.String
	workIterationsBinding binding.String

	workDuration  binding.String
	breakDuration binding.String

	workDurationInt  int
	breakDurationInt int

	timer *time.Timer
}

func createPomodoroTimer(app fyne.App) *PomodoroTimer {
	workDurationPref := app.Preferences().IntWithFallback("workDuration", 25)
	breakDurationPref := app.Preferences().IntWithFallback("breakDuration", 5)
	workIterationsPref := app.Preferences().IntWithFallback("workIterations", 4)

	workDuration := strconv.Itoa(workDurationPref)
	breakDuration := strconv.Itoa(breakDurationPref)
	workIterations := strconv.Itoa(workIterationsPref)

	return &PomodoroTimer{
		workDurationBinding:   binding.BindString(&workDuration),
		breakDurationBinding:  binding.BindString(&breakDuration),
		workIterationsBinding: binding.BindString(&workIterations),

		workDurationInt:  workDurationPref,
		breakDurationInt: breakDurationPref,
	}
}

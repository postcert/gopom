package main

import (
	"time"

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

	Timer *time.Timer
}

func createPomodoroTimer(timingConfig *TimingConfig) *PomodoroTimer {
	workDurationBinding := binding.IntToString(timingConfig.WorkDuration)
	breakDurationBinding := binding.IntToString(timingConfig.BreakDuration)
	workIterationsBinding := binding.IntToString(timingConfig.WorkIterations)

	return &PomodoroTimer{
		workDurationBinding:   workDurationBinding,
		breakDurationBinding:  breakDurationBinding,
		workIterationsBinding: workIterationsBinding,
	}
}

func (timer *PomodoroTimer) UpdateConfig(timingConfig *TimingConfig) {
	timer.workDurationBinding = binding.IntToString(timingConfig.WorkDuration)
	timer.breakDurationBinding = binding.IntToString(timingConfig.BreakDuration)
	timer.workIterationsBinding = binding.IntToString(timingConfig.WorkIterations)
}

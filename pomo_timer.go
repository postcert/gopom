package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/sirupsen/logrus"
)

const (
	StateWork = iota
	StateBreak
	StateLongBreak
)

const (
	StatusStopped = iota
	StatusPaused
	StatusRunning
)

type PomodoroTimer struct {
	Config *TimingConfig

	State         int
	Duration      *time.Duration
	TotalDuration time.Duration
	Timer         *time.Timer
	Ticker        *time.Ticker

	DurationString binding.String
	ProgressFloat  binding.Float
}

func createPomodoroTimer(timingConfig *TimingConfig) *PomodoroTimer {
	timer := &PomodoroTimer{
		Config: timingConfig,
	}
	timer.Init()

	return timer
}

func (timer *PomodoroTimer) Init() {
	timer.Reset()
}

func (timer *PomodoroTimer) Stop() {
	if timer.Timer != nil {
		timer.Timer.Stop()
	}
	if timer.Ticker != nil {
		timer.Ticker.Stop()
	}
}

func (timer *PomodoroTimer) Reset() {
	timer.Stop()

	switch timer.State {
	case StateWork:
		workDuration, error := timer.Config.WorkDuration.Get()
		if error != nil {
			logrus.WithError(error).Errorf("Failed to query bound workDuration")
		}
		duration := time.Duration(workDuration) * time.Minute
		timer.TotalDuration = duration
		timer.Duration = &duration
	case StateBreak:
		breakDuration, error := timer.Config.BreakDuration.Get()
		if error != nil {
			logrus.WithError(error).Errorf("Failed to query bound breakDuration")
		}
		duration := time.Duration(breakDuration) * time.Minute
		timer.TotalDuration = duration
		timer.Duration = &duration
	case StateLongBreak:
		longBreakDuration, error := timer.Config.LongBreakDuration.Get()
		if error != nil {
			logrus.WithError(error).Errorf("Failed to query bound longBreakDuration")
		}
		duration := time.Duration(longBreakDuration) * time.Minute
		timer.TotalDuration = duration
		timer.Duration = &duration
	}

	if timer.DurationString == nil {
		timer.DurationString = binding.NewString()
	}
	timer.DurationString.Set(formatDuration(*timer.Duration))

	if timer.ProgressFloat == nil {
		timer.ProgressFloat = binding.NewFloat()
	}
	timer.ProgressFloat.Set(0)
}

func (timer *PomodoroTimer) UpdateConfig(timingConfig *TimingConfig) {
	timer.Config = timingConfig
}

func (timer *PomodoroTimer) CalculateProgress() {
	progress := 1.0 - float64(*timer.Duration)/float64(timer.TotalDuration)
	timer.ProgressFloat.Set(progress)
}

func (timer *PomodoroTimer) Start() {
	timer.Timer = time.NewTimer(*timer.Duration)
	timer.Ticker = time.NewTicker(time.Second)

	logrus.Debug("PomodoroTimer.Start")
	go func() {
		for {
			select {
			case <-timer.Ticker.C:
				*timer.Duration -= time.Second
				timer.DurationString.Set(formatDuration(*timer.Duration))
				timer.CalculateProgress()
				progress, _ := timer.ProgressFloat.Get()
				logrus.WithFields(logrus.Fields{
					"progress": progress,
				}).Debug("Tick")
			case <-timer.Timer.C:
				logrus.Debug("Timer")
				timer.ProgressFloat.Set(1)
				timer.DurationString.Set(formatDuration(0 * time.Second))
				timer.Ticker.Stop()
				return
			}
		}
	}()
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

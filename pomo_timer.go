package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/sirupsen/logrus"
)

type TimerState int
type TimerStatus int

const (
	StateWork TimerState = iota
	StateBreak
	StateLongBreak
)

const (
	StatusStopped TimerStatus = iota
	StatusPaused
	StatusRunning
)

type PomodoroTimer struct {
	Config *TimingConfig

	State  TimerState
	Status TimerStatus

	Duration       *time.Duration
	TotalDuration  time.Duration
	WorkIterations int
	Timer          *time.Timer
	Ticker         *time.Ticker

	DurationString binding.String
	ProgressFloat  binding.Float

	EventChannel chan TimerEvent
}

func createPomodoroTimer(timingConfig *TimingConfig, eventChannel chan TimerEvent) *PomodoroTimer {
	timer := &PomodoroTimer{
		Config: timingConfig,

		State:  StateWork,
		Status: StatusStopped,

		EventChannel: eventChannel,
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

func intToDuration(i int) time.Duration {
	if debug {
		// High Speed Testing
		return time.Duration(i) * time.Second
	} else {
		return time.Duration(i) * time.Minute
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
		duration := intToDuration(workDuration)
		timer.TotalDuration = duration
		timer.Duration = &duration
	case StateBreak:
		breakDuration, error := timer.Config.BreakDuration.Get()
		if error != nil {
			logrus.WithError(error).Errorf("Failed to query bound breakDuration")
		}
		duration := intToDuration(breakDuration)
		timer.TotalDuration = duration
		timer.Duration = &duration
	case StateLongBreak:
		longBreakDuration, error := timer.Config.LongBreakDuration.Get()
		if error != nil {
			logrus.WithError(error).Errorf("Failed to query bound longBreakDuration")
		}
		duration := intToDuration(longBreakDuration)
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

	if timer.WorkIterations == 0 {
		workIterations, error := timer.Config.WorkIterations.Get()
		if error != nil {
			logrus.WithError(error).Errorf("Failed to query bound workIterations")
		}
		timer.WorkIterations = workIterations
	}
}

func (timer *PomodoroTimer) TransitionState() {
	switch timer.State {
	case StateWork:
		timer.WorkIterations--

		if timer.WorkIterations == 0 {
			logrus.WithField("workIterations", timer.WorkIterations).Debug("TransitionState: StateWork: Transitioning to StateLongBreak")
			timer.State = StateLongBreak
			// TODO: Reset work iterations
		} else {
			logrus.WithField("workIterations", timer.WorkIterations).Debug("TransitionState: StateWork: Transitioning to StateBreak")
			timer.State = StateBreak
		}

	case StateBreak:
		logrus.Debug("TransitionState: StateBreak: Transitioning to StateWork")
		timer.State = StateWork

	case StateLongBreak:
		logrus.Debug("TransitionState: StateLongBreak: Transitioning to StateWork")
		timer.State = StateWork
	}
	timer.Reset()
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
				logrus.Debug("Timer Complete")
				timer.ProgressFloat.Set(1)
				timer.DurationString.Set(formatDuration(0 * time.Second))
				timer.Ticker.Stop()
				timer.EventChannel <- TimerEvent{
					Type: TimerCompleteEvent,
				}
				timer.TransitionState()
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

package main

import (
	"fyne.io/fyne/v2/data/binding"

	"github.com/sirupsen/logrus"
)

type TimingConfig struct {
	WorkDuration   binding.Int
	BreakDuration  binding.Int
	WorkIterations binding.Int
	AutoStartNext  binding.Bool
}

func newDefaultTimingConfig() TimingConfig {
	workDuration := binding.NewInt()
	workDuration.Set(WorkDurationDefault)
	breakDuration := binding.NewInt()
	breakDuration.Set(BreakDurationDefault)
	workIterations := binding.NewInt()
	workIterations.Set(WorkIterationsDefault)
	autoStartNext := binding.NewBool()
	autoStartNext.Set(AutoStartNextDefault)

	return TimingConfig{
		WorkDuration:   workDuration,
		BreakDuration:  breakDuration,
		WorkIterations: workIterations,
		AutoStartNext:  autoStartNext,
	}
}

func (config TimingConfig) intList() []int {
	intList := make([]int, PrefMappingCount)
	workDuration, error := config.WorkDuration.Get()
	if error != nil {
		logrus.WithError(error).Errorf("Failed to query bound workDuration")
	}

	breakDuration, error := config.BreakDuration.Get()
	if error != nil {
		logrus.WithError(error).Errorf("Failed to query bound breakDuration")
	}

	workIterations, error := config.WorkIterations.Get()
	if error != nil {
		logrus.WithError(error).Errorf("Failed to query bound workIterations")
	}

	autostartNext, error := config.AutoStartNext.Get()
	if error != nil {
		logrus.WithError(error).Errorf("Failed to query bound autostartNext value")
	}

	intList[WorkDurationPrefIndex] = workDuration
	intList[BreakDurationPrefIndex] = breakDuration
	intList[WorkIterationsPrefIndex] = workIterations
	intList[AutoStartNextPrefIndex] = btoi(autostartNext)

	return intList
}

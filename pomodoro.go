package main

import (
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"golang.org/x/exp/maps"

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

var timingConfigDefaults = []int{WorkIterationsDefault, WorkDurationDefault, BreakDurationDefault, btoi(AutoStartNextDefault)}

type PomoConfig struct {
	TimingConfigs  map[string]TimingConfig
	SelectedConfig string
	Timer          *PomodoroTimer

	fyneApp fyne.App
}

func createPomoConfig(app fyne.App) *PomoConfig {
	logger := logrus.WithFields(logrus.Fields{
		"function": "createPomoConfig",
	})

	timingConfigNames := app.Preferences().StringListWithFallback(timingConfigKey, DefaultTimingConfigsNames)
	logger.Debug("timingConfigNames: ", timingConfigNames)

	if len(timingConfigNames) == 0 {
		timingConfigNames = append(timingConfigNames, DefaultTimingConfigName)
	}

	sort.Strings(timingConfigNames)

	timingConfigs := LoadTimingConfigPreferences(app, timingConfigNames)
	previousTimingConfigName := LoadPreviousTimingConfigName(app, timingConfigNames)
	previousTimingConfig := LoadTimingConfig(previousTimingConfigName, timingConfigNames, timingConfigs)

	timer := createPomodoroTimer(previousTimingConfig)

	return &PomoConfig{
		TimingConfigs: timingConfigs,
		Timer:         timer,

		fyneApp: app,
	}
}

func LoadTimingConfigPreferences(app fyne.App, timingConfigNames []string) map[string]TimingConfig {
	logger := logrus.WithFields(logrus.Fields{
		"function": "LoadTimingConfigPreferences",
	})

	timingConfigs := make(map[string]TimingConfig)

	for _, configName := range timingConfigNames {
		timingConfigPref := app.Preferences().IntListWithFallback(configName, timingConfigDefaults)

		workDuration := binding.NewInt()
		workDurationPref := getIndexWithFallback(WorkDurationPrefIndex, WorkDurationDefault, timingConfigPref)
		error := workDuration.Set(workDurationPref)
		if error != nil {
			logger.WithError(error).Errorf("Error setting workDurationPref: %d for timingConfig: %s", workDurationPref, configName)
		}

		breakDuration := binding.NewInt()
		breakDurationPref := getIndexWithFallback(BreakDurationPrefIndex, BreakDurationDefault, timingConfigPref)
		error = breakDuration.Set(breakDurationPref)
		if error != nil {
			logger.WithError(error).Errorf("Error setting breakDurationPref: %d for timingConfig: %s", breakDurationPref, configName)
		}

		workIterations := binding.NewInt()
		workIterationsPref := getIndexWithFallback(WorkIterationsPrefIndex, WorkIterationsDefault, timingConfigPref)
		error = workIterations.Set(workIterationsPref)
		if error != nil {
			logger.WithError(error).Errorf("Error setting workIterationsPref: %d for timingConfig: %s", workIterationsPref, configName)
		}

		autoStartNext := binding.NewBool()
		autoStartNextPref := getIndexWithFallback(AutoStartNextPrefIndex, btoi(AutoStartNextDefault), timingConfigPref)
		error = autoStartNext.Set(itob(autoStartNextPref))
		if error != nil {
			logger.WithError(error).Errorf("Error setting autoStartNextPref: %d for timingConfig: %s", autoStartNextPref, configName)
		}

		timingConfigs[configName] = TimingConfig{
			WorkDuration:   workDuration,
			BreakDuration:  breakDuration,
			WorkIterations: workIterations,
			AutoStartNext:  autoStartNext,
		}
	}

	return timingConfigs
}

func LoadPreviousTimingConfigName(app fyne.App, timingConfigNames []string) string {
	fallbackConfigName := timingConfigNames[0]
	previousConfig := app.Preferences().StringWithFallback(prevTimingConfigKey, fallbackConfigName)

	return previousConfig
}

func LoadTimingConfig(requestedConfig string, timingConfigNames []string, timingConfigs map[string]TimingConfig) *TimingConfig {
	fallbackConfigName := timingConfigNames[0]
	timingConfig, ok := timingConfigs[requestedConfig]
	if !ok {
		// Default to first timingConfigName if previousTimingConfig is not found
		timingConfig = timingConfigs[fallbackConfigName]
		logrus.Errorf("LoadTimingConfig: loading requestedConfig: %s failed, fallbackConfigName: %s", requestedConfig, fallbackConfigName)
	}

	return &timingConfig
}

func (config *PomoConfig) Save() {
	for configName, timingConfig := range config.TimingConfigs {
		config.fyneApp.Preferences().SetIntList(configName, timingConfig.intList())
	}

	timingConfigNames := maps.Keys(config.TimingConfigs)
	logrus.Debug("PomoConfig.Save: timingConfigNames: ", timingConfigNames)
	config.fyneApp.Preferences().SetStringList(timingConfigKey, timingConfigNames)
}

func (config *PomoConfig) NewTimingConfig(configName string) {
	timingConfig := newDefaultTimingConfig()

	config.TimingConfigs[configName] = timingConfig
}

func (config *PomoConfig) DeleteTimingConfig(configName string) {
	_, ok := config.TimingConfigs[configName]
	if ok {
		delete(config.TimingConfigs, configName)
		config.fyneApp.Preferences().RemoveValue(configName)
	}
}

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

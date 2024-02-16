package main

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

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

		longBreakDuration := binding.NewInt()
		longBreakDurationPref := getIndexWithFallback(LongBreakDurationPrefIndex, LongBreakDurationDefault, timingConfigPref)
		error = longBreakDuration.Set(longBreakDurationPref)
		if error != nil {
			logger.WithError(error).Errorf("Error setting longBreakDurationPref: %d for timingConfig: %s", longBreakDurationPref, configName)
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
			WorkDuration:      workDuration,
			BreakDuration:     breakDuration,
			LongBreakDuration: longBreakDuration,
			WorkIterations:    workIterations,
			AutoStartNext:     autoStartNext,
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

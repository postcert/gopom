package main

import (
	"sort"

	"fyne.io/fyne/v2"
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
	timingConfigs := make(map[string]TimingConfig)

	for _, configName := range timingConfigNames {
		timingConfigs[configName] = newTimingConfigFromPrefs(configName, app.Preferences())
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
	// TODO: Check into prefBoundValues and see if manually saving is necessary
	//
	// for configName, timingConfig := range config.TimingConfigs {
	// 	config.fyneApp.Preferences().SetIntList(configName, timingConfig.intList())
	// }

	timingConfigNames := maps.Keys(config.TimingConfigs)
	logrus.Debug("PomoConfig.Save: timingConfigNames: ", timingConfigNames)
	config.fyneApp.Preferences().SetStringList(timingConfigKey, timingConfigNames)
}

func (config *PomoConfig) NewTimingConfig(configName string) {
	timingConfig := newDefaultTimingConfig(configName, config.fyneApp.Preferences())

	config.TimingConfigs[configName] = timingConfig
}

func (config *PomoConfig) DeleteTimingConfig(configName string) {
	_, ok := config.TimingConfigs[configName]
	if ok {
		delete(config.TimingConfigs, configName)
		config.fyneApp.Preferences().RemoveValue(configName)
	}
}

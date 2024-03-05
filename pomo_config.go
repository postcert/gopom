package main

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type PomoConfig struct {
	TimingConfigs  map[string]TimingConfig
	SelectedConfig binding.String
	Timer          *PomodoroTimer
	TimerChannel   chan TimerEvent

	fyneApp fyne.App
}

func createPomoConfig(app fyne.App) *PomoConfig {
	timingConfigs := getTimingConfigsFromPrefs(app.Preferences())
	selectedConfigBinding, selectedConfig := getSelectedTimingConfigFromPrefs(app.Preferences(), maps.Keys(timingConfigs))

	selectedTimingConfig := LoadTimingConfig(selectedConfig, timingConfigs)

	timerChannel := make(chan TimerEvent)
	timer := createPomodoroTimer(selectedTimingConfig, timerChannel)

	config := &PomoConfig{
		TimingConfigs:  timingConfigs,
		SelectedConfig: selectedConfigBinding,
		TimerChannel:   timerChannel,
		Timer:          timer,

		fyneApp: app,
	}

	config.HandleTimerEvents()

	return config
}

func (config *PomoConfig) HandleTimerEvents() {
	go func() {
		for event := range config.TimerChannel {
			switch event.Type {
			case TimerCompleteEvent:
				logrus.Debug("PomoConfig.HandleTimerEvents: TimerCompleteEvent")
			case TimerPauseEvent:
				logrus.Debug("PomoConfig.HandleTimerEvents: TimerPauseEvent")
			case TimerStopEvent:
				logrus.Debug("PomoConfig.HandleTimerEvents: TimerStopEvent")
			}
		}
	}()
}

func getTimingConfigsFromPrefs(preferences fyne.Preferences) map[string]TimingConfig {
	timingConfigNames := preferences.StringListWithFallback(timingConfigKey, DefaultTimingConfigsNames)
	if len(timingConfigNames) == 0 {
		timingConfigNames = append(timingConfigNames, DefaultTimingConfigName)
	}
	sort.Strings(timingConfigNames)

	timingConfigs := make(map[string]TimingConfig)
	for _, configName := range timingConfigNames {
		timingConfigs[configName] = newTimingConfigFromPrefs(configName, preferences)
	}
	return timingConfigs
}

func getSelectedTimingConfigFromPrefs(preferences fyne.Preferences, timingConfigNames []string) (binding.String, string) {
	selectedConfig := binding.BindPreferenceString(selectedTimingConfigKey, preferences)
	selectedConfigStr, err := selectedConfig.Get()
	if err != nil || selectedConfigStr == "" || !slices.Contains(timingConfigNames, selectedConfigStr) {
		selectedConfig.Set(timingConfigNames[0])
		selectedConfigStr = timingConfigNames[0]
	}

	return selectedConfig, selectedConfigStr
}

func LoadTimingConfig(requestedConfig string, timingConfigs map[string]TimingConfig) *TimingConfig {
	fallbackConfigName := maps.Keys(timingConfigs)[0]

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

func (config *PomoConfig) StartTimer() {
	config.Timer.Start()
}

func (config *PomoConfig) StopTimer() {
	config.Timer.Stop()
}

func (config *PomoConfig) ResetTimer() {
	config.Timer.Reset()
}

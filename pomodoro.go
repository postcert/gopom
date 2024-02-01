package main

import (
	"fmt"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"golang.org/x/exp/maps"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type timingConfig struct {
	WorkDuration   binding.Int
	BreakDuration  binding.Int
	WorkIterations binding.Int
	AutoStartNext  binding.Bool
}

func newDefaultTimingConfig() timingConfig {
	workDuration := binding.NewInt()
	workDuration.Set(WorkDurationDefault)
	breakDuration := binding.NewInt()
	breakDuration.Set(BreakDurationDefault)
	workIterations := binding.NewInt()
	workIterations.Set(WorkIterationsDefault)
	autoStartNext := binding.NewBool()
	autoStartNext.Set(AutoStartNextDefault)

	return timingConfig{
		WorkDuration:   workDuration,
		BreakDuration:  breakDuration,
		WorkIterations: workIterations,
		AutoStartNext:  autoStartNext,
	}
}

func (config timingConfig) intList() []int {
	intList := make([]int, PrefMappingCount)
	workDuration, error := config.WorkDuration.Get()
	if error != nil {
		log.WithError(error).Errorf("Failed to query bound workDuration")
	}

	breakDuration, error := config.BreakDuration.Get()
	if error != nil {
		log.WithError(error).Errorf("Failed to query bound breakDuration")
	}

	workIterations, error := config.WorkIterations.Get()
	if error != nil {
		log.WithError(error).Errorf("Failed to query bound workIterations")
	}

	autostartNext, error := config.AutoStartNext.Get()
	if error != nil {
		log.WithError(error).Errorf("Failed to query bound autostartNext value")
	}

	intList[WorkDurationPrefIndex] = workDuration
	intList[BreakDurationPrefIndex] = breakDuration
	intList[WorkIterationsPrefIndex] = workIterations
	intList[AutoStartNextPrefIndex] = btoi(autostartNext)

	return intList
}

var timingConfigDefaults = []int{WorkIterationsDefault, WorkDurationDefault, BreakDurationDefault, btoi(AutoStartNextDefault)}

type PomoConfig struct {
	TimingConfigs map[string]timingConfig

	fyneApp fyne.App
}

func createPomoConfig(app fyne.App) *PomoConfig {
	logger := log.WithFields(logrus.Fields{
		"function": "createPomoConfig",
	})

	timingConfigs := make(map[string]timingConfig)
	timingConfigNames := app.Preferences().StringListWithFallback(timingConfigKey, DefaultTimingConfigsNames)
	logger.Debug("timingConfigNames: ", timingConfigNames)

	if len(timingConfigNames) == 0 {
		timingConfigNames = append(timingConfigNames, DefaultTimingConfigName)
	}

	sort.Strings(timingConfigNames)

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

		timingConfigs[configName] = timingConfig{
			WorkDuration:   workDuration,
			BreakDuration:  breakDuration,
			WorkIterations: workIterations,
			AutoStartNext:  autoStartNext,
		}
	}

	return &PomoConfig{
		TimingConfigs: timingConfigs,

		fyneApp: app,
	}
}

func (config *PomoConfig) Save() {
	for configName, timingConfig := range config.TimingConfigs {
		config.fyneApp.Preferences().SetIntList(configName, timingConfig.intList())
	}

	timingConfigNames := maps.Keys(config.TimingConfigs)
	log.Debug("PomoConfig.Save: timingConfigNames: ", timingConfigNames)
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

	config *PomoConfig
	Timer  *time.Timer
}

func createPomodoroTimer(config *PomoConfig) *PomodoroTimer {
	workDurationBinding := binding.IntToString(config.TimingConfigs[DefaultTimingConfigName].WorkDuration)
	breakDurationBinding := binding.IntToString(config.TimingConfigs[DefaultTimingConfigName].BreakDuration)
	workIterationsBinding := binding.IntToString(config.TimingConfigs[DefaultTimingConfigName].WorkIterations)

	return &PomodoroTimer{
		workDurationBinding:   workDurationBinding,
		breakDurationBinding:  breakDurationBinding,
		workIterationsBinding: workIterationsBinding,

		config: config,
	}
}

func (pomodoroTimer *PomodoroTimer) save() {
	fmt.Println("pomodoroTimer.save")
}

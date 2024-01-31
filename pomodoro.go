package main

import (
	"fmt"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"golang.org/x/exp/maps"

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
	// TODO: Handle binding.get() errors
	workDuration, _ := config.WorkDuration.Get()
	breakDuration, _ := config.BreakDuration.Get()
	workIterations, _ := config.WorkIterations.Get()
	autostartNext, _ := config.AutoStartNext.Get()

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
	timingConfigs := make(map[string]timingConfig)
	timingConfigNames := app.Preferences().StringListWithFallback(timingConfigKey, DefaultTimingConfigsNames)
	log.Debug("createPomoConfig: timingConfigNames: ", timingConfigNames)

	if len(timingConfigNames) == 0 {
		timingConfigNames = append(timingConfigNames, DefaultTimingConfigName)
	}

	// If DefaultTimingConfigName is not in timingConfigNames, add it
	// if !slices.Contains(timingConfigNames, DefaultTimingConfigName) {
	// 	timingConfigNames = append(timingConfigNames, DefaultTimingConfigName)
	// }

	sort.Strings(timingConfigNames)

	for _, configName := range timingConfigNames {
		timingConfigPref := app.Preferences().IntListWithFallback(configName, timingConfigDefaults)

		// TODO: Handle binding.Set() errors
		workDuration := binding.NewInt()
		workDurationPref := getIndexWithFallback(WorkDurationPrefIndex, WorkDurationDefault, timingConfigPref)
		workDuration.Set(workDurationPref)

		breakDuration := binding.NewInt()
		breakDurationPref := getIndexWithFallback(BreakDurationPrefIndex, BreakDurationDefault, timingConfigPref)
		breakDuration.Set(breakDurationPref)

		workIterations := binding.NewInt()
		workIterationsPref := getIndexWithFallback(WorkIterationsPrefIndex, WorkIterationsDefault, timingConfigPref)
		workIterations.Set(workIterationsPref)

		autoStartNext := binding.NewBool()
		autoStartNextPref := getIndexWithFallback(AutoStartNextPrefIndex, btoi(AutoStartNextDefault), timingConfigPref)
		autoStartNext.Set(itob(autoStartNextPref))

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

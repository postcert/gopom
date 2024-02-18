package main

import (
	"fmt"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/sirupsen/logrus"
)

type TimingConfig struct {
	WorkDuration      binding.Int
	BreakDuration     binding.Int
	LongBreakDuration binding.Int
	WorkIterations    binding.Int
	AutoStartNext     binding.Bool
}

var timingConfigDefaults = map[string]interface{}{
	"WorkDuration":      WorkDurationDefault,
	"BreakDuration":     BreakDurationDefault,
	"LongBreakDuration": LongBreakDurationDefault,
	"WorkIterations":    WorkIterationsDefault,
	"AutoStartNext":     AutoStartNextDefault,
}

func newTimingConfigFromPrefs(configName string, preferences fyne.Preferences) TimingConfig {
	// workDurationPrefKey := fmt.Sprintf(WorkDurationKey, configName)
	// workDuration := binding.BindPreferenceInt(workDurationPrefKey, preferences)
	//
	// breakDurationPrefKey := fmt.Sprintf(BreakDurationKey, configName)
	// breakDuration := binding.BindPreferenceInt(breakDurationPrefKey, preferences)
	//
	// longBreakDurationPrefKey := fmt.Sprintf(LongBreakDurationKey, configName)
	// longBreakDuration := binding.BindPreferenceInt(longBreakDurationPrefKey, preferences)
	//
	// workIterationsPrefKey := fmt.Sprintf(WorkIterationsKey, configName)
	// workIterations := binding.BindPreferenceInt(workIterationsPrefKey, preferences)
	//
	// autoStartNextPrefKey := fmt.Sprintf(AutoStartNextKey, configName)
	// autoStartNext := binding.BindPreferenceBool(autoStartNextPrefKey, preferences)

	timingConfig := TimingConfig{}

	tConfigValue := reflect.ValueOf(&timingConfig).Elem()
	tConfigType := tConfigValue.Type()

	intType := reflect.TypeOf((*binding.Int)(nil)).Elem()
	boolType := reflect.TypeOf((*binding.Bool)(nil)).Elem()

	for i := 0; i < tConfigValue.NumField(); i++ {
		field := tConfigValue.Field(i)
		fieldType := field.Type()
		fieldName := tConfigType.Field(i).Name

		prefKey := fmt.Sprintf("%s_%s", configName, fieldName)

		if field.Type().Implements(intType) {
			if !field.CanSet() {
				continue
			}
			field.Set(reflect.ValueOf(binding.BindPreferenceInt(prefKey, preferences)))
		} else if field.Type().Implements(boolType) {
			if !field.CanSet() {
				continue
			}
			field.Set(reflect.ValueOf(binding.BindPreferenceBool(prefKey, preferences)))
		} else {
			logrus.Errorf("Unsupported type: %s", fieldType)
		}
	}

	return timingConfig
}

func newDefaultTimingConfig(configName string, preferences fyne.Preferences) TimingConfig {
	timingConfig := newTimingConfigFromPrefs(configName, preferences)

	timingConfig.WorkDuration.Set(WorkDurationDefault)
	timingConfig.BreakDuration.Set(BreakDurationDefault)
	timingConfig.LongBreakDuration.Set(LongBreakDurationDefault)
	timingConfig.WorkIterations.Set(WorkIterationsDefault)
	timingConfig.AutoStartNext.Set(AutoStartNextDefault)

	return timingConfig
}

func deleteTimingConfig(configName string, preferences fyne.Preferences) {
	for key := range timingConfigDefaults {
		prefKey := fmt.Sprintf("%s_%s", configName, key)
		preferences.RemoveValue(prefKey)
	}
}

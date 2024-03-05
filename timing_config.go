package main

import (
	"fmt"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"github.com/sirupsen/logrus"
)

type TimingConfig struct {
	WorkDuration      binding.Int
	BreakDuration     binding.Int
	LongBreakDuration binding.Int
	WorkIterations    binding.Int
	AutoStartNext     binding.Bool

	WorkFinishedSound      binding.String
	BreakFinishedSound     binding.String
	PlayWorkFinishedSound  binding.Bool
	PlayBreakFinishedSound binding.Bool
}

var timingConfigDefaults = map[string]interface{}{
	"WorkDuration":      WorkDurationDefault,
	"BreakDuration":     BreakDurationDefault,
	"LongBreakDuration": LongBreakDurationDefault,
	"WorkIterations":    WorkIterationsDefault,
	"AutoStartNext":     AutoStartNextDefault,

	"WorkFinishedSound":      WorkFinishedSoundDefault,
	"BreakFinishedSound":     BreakFinishedSoundDefault,
	"PlayWorkFinishedSound":  PlayWorkFinishedSoundDefault,
	"PlayBreakFinishedSound": PlayBreakFinishedSoundDefault,
}

func newTimingConfigFromPrefs(configName string, preferences fyne.Preferences) TimingConfig {
	timingConfig := TimingConfig{}

	// Not the cleanest use of reflection due to the binding.* being Interfaces
	// But it hopefully avoids bugs in boilerplate code
	tConfigValue := reflect.ValueOf(&timingConfig).Elem()
	tConfigType := tConfigValue.Type()

	// Getting the interface types because "typeOf(binding.Int)" was failing
	// Fyne Preferences only supports these 4 types so at least all bases are covered
	intType := reflect.TypeOf((*binding.Int)(nil)).Elem()
	boolType := reflect.TypeOf((*binding.Bool)(nil)).Elem()
	stringType := reflect.TypeOf((*binding.String)(nil)).Elem()
	floatType := reflect.TypeOf((*binding.Float)(nil)).Elem()

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
		} else if field.Type().Implements(stringType) {
			if !field.CanSet() {
				continue
			}
			field.Set(reflect.ValueOf(binding.BindPreferenceString(prefKey, preferences)))
		} else if field.Type().Implements(floatType) {
			if !field.CanSet() {
				continue
			}
			field.Set(reflect.ValueOf(binding.BindPreferenceFloat(prefKey, preferences)))
		} else {
			logrus.Errorf("Unsupported type when Binding to a preference: %s", fieldType)
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

	timingConfig.WorkFinishedSound.Set(WorkFinishedSoundDefault)
	timingConfig.BreakFinishedSound.Set(BreakFinishedSoundDefault)
	timingConfig.PlayWorkFinishedSound.Set(PlayWorkFinishedSoundDefault)
	timingConfig.PlayBreakFinishedSound.Set(PlayBreakFinishedSoundDefault)

	return timingConfig
}

func deleteTimingConfig(configName string, preferences fyne.Preferences) {
	for key := range timingConfigDefaults {
		prefKey := fmt.Sprintf("%s_%s", configName, key)
		preferences.RemoveValue(prefKey)
	}
}

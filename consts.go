package main

const (
	WorkIterationsDefault    = 4
	WorkDurationDefault      = 25
	BreakDurationDefault     = 5
	LongBreakDurationDefault = 15
	AutoStartNextDefault     = true

	WorkFinishedSoundDefault      = "tbd"
	BreakFinishedSoundDefault     = "tbd"
	PlayWorkFinishedSoundDefault  = true
	PlayBreakFinishedSoundDefault = true
)

const (
	timingConfigKey         = "timingConfigs"
	selectedTimingConfigKey = "selectedTimingConfig"
)

const DefaultTimingConfigName = "Default"

var DefaultTimingConfigsNames = []string{DefaultTimingConfigName}

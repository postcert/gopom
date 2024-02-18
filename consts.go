package main

const (
	WorkIterationsPrefIndex = iota
	WorkDurationPrefIndex
	BreakDurationPrefIndex
	LongBreakDurationPrefIndex
	AutoStartNextPrefIndex
	PrefMappingCount
)

const (
	WorkIterationsDefault    = 4
	WorkDurationDefault      = 25
	BreakDurationDefault     = 5
	LongBreakDurationDefault = 15
	AutoStartNextDefault     = true
)

const (
	timingConfigKey     = "timingConfigs"
	prevTimingConfigKey = "prevTimingConfig"
)

const DefaultTimingConfigName = "Default"

var DefaultTimingConfigsNames = []string{DefaultTimingConfigName}

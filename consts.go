package main

const (
	WorkIterationsPrefIndex = iota
	WorkDurationPrefIndex
	BreakDurationPrefIndex
	AutoStartNextPrefIndex
	PrefMappingCount
)

const (
	WorkIterationsDefault = 4
	WorkDurationDefault   = 25
	BreakDurationDefault  = 5
	AutoStartNextDefault  = true
)

const timingConfigKey = "timingConfigs"

const DefaultTimingConfigName = "Default"

var DefaultTimingConfigsNames = []string{DefaultTimingConfigName}

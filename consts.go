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

const (
	timingConfigKey     = "timingConfigs"
	prevTimingConfigKey = "prevTimingConfig"
)

const DefaultTimingConfigName = "Default"

var (
	DefaultTimingConfigsNames = []string{DefaultTimingConfigName}
	timingConfigDefaults      = []int{WorkIterationsDefault, WorkDurationDefault, BreakDurationDefault, btoi(AutoStartNextDefault)}
)

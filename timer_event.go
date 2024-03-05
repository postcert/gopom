package main

type TimerEventType int

const (
	UndefinedEvent TimerEventType = iota
	TimerCompleteEvent
	TimerPauseEvent
	TimerStopEvent
)

type TimerEvent struct {
	Type    TimerEventType
	Payload interface{}
}

package timeouts

import "time"

func TwoMinuteWatch() *int64 {
	timeout := int64(60 * 2)
	return &timeout
}

func ThreeMinuteWatch() *int64 {
	timeout := int64(60 * 3)
	return &timeout
}

func TenMinuteWatch() *int64 {
	timeout := int64(60 * 10)
	return &timeout
}

func TwentyMinuteWatch() *int64 {
	timeout := int64(60 * 20)
	return &timeout
}

func ThirtyMinuteWatch() *int64 {
	timeout := int64(60 * 30)
	return &timeout
}

const (
	FiveHundredMillisecond = 500 * time.Millisecond
	FiveSecond             = 5 * time.Second
	TenSecond              = 10 * time.Second
	OneMinute              = 1 * time.Minute
	TwoMinute              = 2 * time.Minute
	FiveMinute             = 5 * time.Minute
	TenMinute              = 10 * time.Minute
	FifteenMinute          = 15 * time.Minute
	ThirtyMinute           = 30 * time.Minute
)

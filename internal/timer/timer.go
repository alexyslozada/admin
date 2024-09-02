package timer

import "time"

type Timer interface {
	Now() time.Time
}

type RealTimer struct{}

func NewRealTimer() RealTimer {
	return RealTimer{}
}

func (RealTimer) Now() time.Time {
	return time.Now()
}

type FakeTimer struct {
	NowFn func() time.Time
}

func NewFakeTimer(nowFn func() time.Time) FakeTimer {
	return FakeTimer{NowFn: nowFn}
}

func (f FakeTimer) Now() time.Time {
	return f.NowFn()
}

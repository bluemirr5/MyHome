package main

import (
	"time"
)

type Scheduler struct {
	Period                                                  time.Duration
	YearTick, DayTick, HourTick, MinTick, SecTick, NsecTick int
	MonthTick                                               time.Month
	Runner                                                  func()
}

func NewScheduler(period time.Duration) *Scheduler {
	return &Scheduler{Period: period, YearTick: -1, MonthTick: -1, DayTick: -1, HourTick: -1, MinTick: -1, SecTick: -1, NsecTick: -1}
}

func (s *Scheduler) Run() {
	go s.runningRoutine()
}

func (s *Scheduler) runningRoutine() {
	ticker := s.updateTicker()
	for {
		<-ticker.C
		if s.Runner != nil {
			go s.Runner()
		} else {
			panic("Set runner")
		}
		ticker = s.updateTicker()
	}
}

func (s *Scheduler) updateTicker() *time.Ticker {
	var y, d, h, m, ss, n int = s.YearTick, s.DayTick, s.HourTick, s.MinTick, s.SecTick, s.NsecTick
	mm := s.MonthTick
	if s.YearTick == -1 {
		y = time.Now().Year()
	}
	if s.MonthTick == -1 {
		mm = time.Now().Month()
	}
	if s.DayTick == -1 {
		d = time.Now().Day()
	}
	if s.HourTick == -1 {
		h = time.Now().Hour()
	}
	if s.MinTick == -1 {
		m = time.Now().Minute()
	}
	if s.SecTick == -1 {
		ss = time.Now().Second()
	}
	if s.NsecTick == -1 {
		n = time.Now().Nanosecond()
	}
	nextTick := time.Date(y, mm, d, h, m, ss, n, time.Local)
	diff := nextTick.Sub(time.Now())
	if diff <= 0 {
		nextTick = nextTick.Add(s.Period)
	}
	diff = nextTick.Sub(time.Now())
	return time.NewTicker(diff)
}

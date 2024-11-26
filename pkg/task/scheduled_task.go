package task

import "time"

type ScheduledTask struct {
	Task              Task
	SubmissionTime    time.Time
	LastScheduledTime time.Time
	Policy            TaskSchedulePolicy
}

type TaskSchedulePolicy interface {
}

type DelayedTaskPolicy struct {
	Delay time.Duration
}

type FixedRateTaskPolicy struct {
	Rate time.Duration
}

type OneShotTaskPolicy struct {
}

func (t *ScheduledTask) GetScheduledTime() time.Time {
	_, ok := t.Policy.(OneShotTaskPolicy)
	if ok {
		return t.SubmissionTime
	}

	pl, ok := t.Policy.(DelayedTaskPolicy)
	if ok {
		return t.SubmissionTime.Add(pl.Delay)
	}

	dpl, ok := t.Policy.(FixedRateTaskPolicy)
	if ok {
		defaultTime := time.Time{}
		if t.LastScheduledTime == defaultTime {
			return t.SubmissionTime.Add(dpl.Rate)
		} else {
			return t.LastScheduledTime.Add(dpl.Rate)
		}
	}

	return time.Time{}
}

package main

import (
	"time"

	"go.uber.org/cadence/workflow"
)

func newCommonActivityOptions() workflow.ActivityOptions {
	return workflow.ActivityOptions{
		StartToCloseTimeout:    15 * time.Minute,
		ScheduleToStartTimeout: 15 * time.Minute,
		ScheduleToCloseTimeout: 30 * time.Minute,
	}
}

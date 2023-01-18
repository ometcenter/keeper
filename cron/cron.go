package cron

import "github.com/robfig/cron/v3"

var CronScheduler *cron.Cron

var MapScheduker map[string]cron.EntryID

type SetScheduleCronMessage struct {
	TypeSchedule string
	Name         string
	JobID        string
	Cron         string
	JSONParam    string
}

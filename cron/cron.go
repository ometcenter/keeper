package cron

import (
	"time"

	"github.com/robfig/cron/v3"
)

var CronScheduler *cron.Cron

var MapScheduker map[string]cron.EntryID

type SetScheduleCronMessage struct {
	TypeSchedule string
	Name         string
	JobID        string
	Cron         string
	JSONParam    string
}

type JobsCronBackGroundTask struct {
	ID   int
	Job  string
	Next time.Time
	Prev time.Time
}

type CronTaskFillRedis struct {
	Name        string
	Description string
	JobID       string
}

// func (t *CronTaskFillRedis) Run() {
// 	store.FillDataToRedisGorutine()
// 	fmt.Println("Заполнение кеша редис с крон СТАРТ")
// 	// if err != nil {
// 	// 	log.Impl.Errorf("Ошибка Перезаполнение авторизации в Redis, делати: %s \n")
// 	// }
// }

func InitCron() {

	// Хорошая статья с описанием https://stackoverflow.com/questions/68343512/how-to-find-a-particular-running-cron-jobs-with-github-com-robfig-cron

	CronScheduler = cron.New(cron.WithChain())

	// _, err := CronScheduler.AddJob("30 23 * * *", &CronTaskFillRedis{"Перезаполнение авторизации в Redis", "", ""}) // "30 7 * * *" = 7:30
	// if err != nil {
	// 	log.Impl.Errorf("Ошибка добавления Cron задания: %s", err.Error())
	// }

	CronScheduler.Start()

}

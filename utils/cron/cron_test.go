package cron

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"testing"
	"time"
)

var task = func(job gocron.Job) {
	fmt.Println(1111)
	fmt.Println(job.NextRun())
}

func TestCron(t *testing.T) {
	s := gocron.NewScheduler(time.Local)
	s.SingletonModeAll()
	_, _ = s.Cron("*/1 * * * *").Tag("aaa").DoWithJobDetails(task)
	//_, _ = s.CronWithSeconds("*/5 * * * * *").DoWithJobDetails(task)
	s.StartAsync()

	for true {
		time.Sleep(time.Second * 30)
		err := s.RunByTag("aaa")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
func TestGetCorrespondPath(t *testing.T) {
	ts := time.Now().UnixNano() / int64(time.Millisecond)
	s, err := getCorrespondPath(ts)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)
}
func TestReFreshToken(t *testing.T) {
	var job gocron.Job
	ReFreshToken(job)
}

func TestAutoCollect(t *testing.T) {
	var job gocron.Job
	AutoCollect(job)
}

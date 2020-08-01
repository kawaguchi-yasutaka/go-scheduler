package scheduler

import (
"errors"
"sync"
"time"
)


var (
	errorTooManyArguments = errors.New("too many arguments")
	errorInvalidType = errors.New("invalid type")
)

type scheduler struct {
	taskScheduleWg *sync.WaitGroup
	done           chan bool
}

func NewScheduler() scheduler{
	return scheduler{
		taskScheduleWg: &sync.WaitGroup{},
		done:           make(chan bool),
	}
}

func (s scheduler) Wait() {
	s.taskScheduleWg.Wait()
	<- s.done
}

func (s scheduler) Close() {
	close(s.done)
}


func (s scheduler) Every(interval ...int64) worker{
	switch len(interval) {
	case 0:
		return worker{scheduler: s,schedule: every{interval: interval[0]}}
	case 1:
		return worker{scheduler: s,schedule: every{interval: interval[0]}}
	default:
		return worker{err: errorTooManyArguments}
	}
}

func (s scheduler) FromCronFormat(cron string) worker {
	//空白で、トリムする
	// 中身５つあるか
	//- , number
	// *
	// number-number
	// number,number,number
	// number,number-number

	//間隔
	// */number
	// number-number/number
	//* * * * *
	//日付と曜日が設定されている場合は弾く日宇町ありそう
}

type cronType []int64

func convertCronType(type string, cron string) cronType{
}



type cron struct {
	m string
	h string
	dom string
	mon string
	dow string
	lastRunTime time.Time
}

type every struct {
	 unit time.Duration
	 interval int64
}


type worker struct{
	scheduler scheduler
	schedule schedule
	err error
}

type schedule interface {
	nextRun() time.Duration
}


func (e every) nextRun() time.Duration{
	return time.Duration(e.interval) * e.unit
}


func (w worker)Run(f func()) error{
	if w.err != nil {
		return w.err
	}
	w.scheduler.taskScheduleWg.Add(1)
	go func() {
		w.scheduler.taskScheduleWg.Done()
		for {
			next := w.schedule.nextRun()
			select {
			case <-time.After(next):
				go f()
			case <-w.scheduler.done:
				return
			}
		}
	}()
	return nil
}

func (w worker)Second() worker{
	v, ok := w.schedule.(every)
	if !ok {
		w.err = errorInvalidType
	}
	v.unit = time.Second
	w.schedule = v
	return w
}

func (w worker)Minute() worker{
	v, ok := w.schedule.(every)
	if !ok {
		w.err = errorInvalidType
	}
	v.unit = time.Minute
	w.schedule = v
	return w
}

func (w worker)Hour() worker{
	v, ok := w.schedule.(every)
	if !ok {
		w.err = errorInvalidType
	}
	v.unit = time.Hour
	w.schedule = v
	return w
}

func (w worker)Day() worker{
	v, ok := w.schedule.(every)
	if !ok {
		w.err = errorInvalidType
	}
	v.unit = 24 * time.Hour
	w.schedule = v
	return w
}
package services

//read docs of cron
import (
	"github.com/robfig/cron/v3"
)

// why we doing this? why there is need of scheduler?
type CronService struct {
	Scheduler *cron.Cron
}

// what is cron.New()?
func NewCronService() *CronService {
	return &CronService{
		Scheduler: cron.New(),
	}
}

// what does this function do?
func (cs *CronService) Start() {
	cs.Scheduler.Start()
}

// what does this function do?
func (cs *CronService) Stop() {
	cs.Scheduler.Stop()
}

// what does this function do?
func (cs *CronService) AddFunc(spec string, cmd func()) {
	cs.Scheduler.AddFunc(spec, cmd)
}

// what does this function do? -> this allows more complex job definitions for example stateful jobs.
func (cs *CronService) AddJob(spec string, cmd cron.Job) {
	cs.Scheduler.AddJob(spec, cmd)
}

// what does this function do?
func (cs *CronService) RemoveJob(id cron.EntryID) {
	cs.Scheduler.Remove(id)
}

// what does this function do?
func (cs *CronService) Entries() []cron.Entry {
	return cs.Scheduler.Entries()
}

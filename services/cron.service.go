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

func (cs *CronService) Start() {
	cs.Scheduler.Start()
}

func (cs *CronService) Stop() {
	cs.Scheduler.Stop()
}

func (cs *CronService) AddFunc(spec string, cmd func()) {
	cs.Scheduler.AddFunc(spec, cmd)
}

func (cs *CronService) AddJob(spec string, cmd cron.Job) {
	cs.Scheduler.AddJob(spec, cmd)
}

func (cs *CronService) RemoveJob(id cron.EntryID) {
	cs.Scheduler.Remove(id)
}

func (cs *CronService) Entries() []cron.Entry {
	return cs.Scheduler.Entries()
}

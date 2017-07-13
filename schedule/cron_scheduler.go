package schedule

import "gopkg.in/robfig/cron.v2"

type CronScheduler struct {
	*cron.Cron
}

func NewCronScheduler() *CronScheduler {
	c := cron.New()
	return &CronScheduler{Cron: c}
}

func (c *CronScheduler) AddFunc(spec string, cmd func()) (EntryID, error) {
	entryID, err := c.Cron.AddFunc(spec, cmd)
	return EntryID(entryID), err
}
func (c *CronScheduler) Remove(id EntryID) {
	c.Cron.Remove(cron.EntryID(id))
}

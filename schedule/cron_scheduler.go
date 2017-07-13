package schedule

import (
	"sort"
	"time"

	"github.com/robfig/cron"
)

type CronScheduler struct {
	stop    chan struct{}
	add     chan *Entry
	remove  chan *Entry
	entries []*Entry
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		stop:    make(chan struct{}),
		add:     make(chan *Entry, 10),
		remove:  make(chan *Entry, 10),
		entries: []*Entry{},
	}
}

type Entry struct {
	id       int
	schedule cron.Schedule
	next     time.Time
	cmd      func()
}

type byTime []*Entry

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	if s[i].next.IsZero() {
		return false
	}
	if s[j].next.IsZero() {
		return true
	}
	return s[i].next.Before(s[j].next)
}

func (c *CronScheduler) AddFunc(id int, spec string, cmd func()) error {
	schedule, err := cron.Parse(spec)
	if err != nil {
		return err
	}

	now := time.Now()
	next := schedule.Next(now)

	c.add <- &Entry{
		id:       id,
		schedule: schedule,
		next:     next,
		cmd:      cmd,
	}

	return nil
}

func (c *CronScheduler) Start() {
Loop:
	for {
		sort.Sort(byTime(c.entries))

		select {
		case <-c.stop:
			break Loop
		case e := <-c.add:
			c.entries = append(c.entries, e)
		case e := <-c.remove:
			c.removeEntry(e)
		default:
			c.checker()
		}
	}
}

func (c *CronScheduler) checker() {
	now := time.Now()
	for _, e := range c.entries {
		if e.next.After(now) || e.next.IsZero() {
			break
		}
		go e.cmd()
		e.next = e.schedule.Next(now)
	}
}

func (c *CronScheduler) removeEntry(e *Entry) {
	found := false
	foundID := 0

	for i, entry := range c.entries {
		if entry == e {
			found = true
			foundID = i
			break
		}
	}

	if found {
		c.entries = append(c.entries[:foundID], c.entries[foundID+1:]...)
	}
}

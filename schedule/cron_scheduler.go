package schedule

import (
	"sort"
	"time"

	"github.com/robfig/cron"
)

type CronScheduler struct {
	stop    chan struct{}
	add     chan *Entry
	remove  chan EntryID
	entries []*Entry
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		stop:    make(chan struct{}),
		add:     make(chan *Entry, 10),
		remove:  make(chan EntryID, 10),
		entries: []*Entry{},
	}
}

type Entry struct {
	id       EntryID
	schedule cron.Schedule
	next     time.Time
	cmd      func(id EntryID)
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

func (c *CronScheduler) AddFunc(id EntryID, spec string, cmd func(id EntryID)) error {
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
		case id := <-c.remove:
			c.removeEntryByID(id)
		default:
			c.checker()
		}
	}
}

func (c *CronScheduler) Stop() {
	c.stop <- struct{}{}
}

func (c *CronScheduler) checker() {
	now := time.Now()
	for _, e := range c.entries {
		if e.next.After(now) || e.next.IsZero() {
			break
		}
		go e.cmd(e.id)
		e.next = e.schedule.Next(now)
	}
}

func (c *CronScheduler) removeEntryByID(id EntryID) {
	found := false
	foundID := 0

	for i, entry := range c.entries {
		if entry.id == id {
			found = true
			foundID = i
			break
		}
	}

	if found {
		c.entries = append(c.entries[:foundID], c.entries[foundID+1:]...)
	}
}

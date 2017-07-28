package alarm

import (
	"log"
	"os/exec"
	"strings"

	"github.com/jvikstedt/alarmy/schedule"
	"github.com/jvikstedt/alarmy/store"
)

type Executor struct {
	store  store.Store
	logger *log.Logger
}

func NewExecutor(store store.Store, logger *log.Logger) *Executor {
	return &Executor{
		store:  store,
		logger: logger,
	}
}

func (e *Executor) Execute(entryID schedule.EntryID) {
	job, err := e.store.Job().GetOne(int(entryID))
	if err != nil {
		e.logger.Printf("Err getting job: %d %v", entryID, err)
		return
	}

	out, err := exec.Command("sh", "-c", job.Cmd).Output()

	if err != nil {
		e.logger.Printf("Err executing command for job: %d command: %s, %v", entryID, job.Cmd, err)
		return
	}

	e.logger.Printf("Job %d finished with output %s\n", job.ID, strings.TrimSpace(string(out)))
}

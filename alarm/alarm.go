package alarm

import (
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jvikstedt/alarmy/model"
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

	result := make(map[string]interface{})
	err = json.Unmarshal(out, &result)
	if err != nil {
		e.logger.Printf("Err unmarshalling result for job: %d err: %v", entryID, err)
		return
	}

	for _, t := range job.Triggers {
		field := result[t.FieldName]
		switch v := field.(type) {
		case float64:
			// Check if integer
			if v == float64(int64(v)) {
				e.handleAsInt(job, t, v)
			}

		case string:
		case bool:
		default:
			e.logger.Printf("Unknown type %T for job %d\n", v, job.ID)
		}
	}

	e.logger.Printf("Job %d finished with output %s\n", job.ID, strings.TrimSpace(string(out)))
}

func (e *Executor) handleAsInt(job model.Job, t model.Trigger, value float64) bool {
	target, err := strconv.Atoi(t.Target)
	if err != nil {
		e.logger.Printf("%v", err)
		return false
	}
	actual := int(value)

	switch t.TriggerType {
	case model.TriggerEqual:
		if target != actual {
			e.logger.Printf("Expected %d but got %d for job %s.%s\n", target, actual, job.Name, t.FieldName)
			return true
		}
	case model.TriggerMoreThan:
		if actual > target {
			e.logger.Printf("TriggerMoreThan %d was more than %d for job %s.%s\n", actual, target, job.Name, t.FieldName)
			return true
		}
	default:
		e.logger.Printf("Invalid TriggerType %v for int type\n", t.TriggerType)
	}

	return false
}

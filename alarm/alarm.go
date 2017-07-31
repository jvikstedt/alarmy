package alarm

import (
	"encoding/json"
	"fmt"
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

	resultSet := TriggerResultSet{}

	for _, t := range job.Triggers {
		triggerResult := TriggerResult{}

		field := result[t.FieldName]
		switch v := field.(type) {
		case float64:
			// Check if integer
			if v == float64(int64(v)) {
				e.handleAsInt(job, t, v, &triggerResult)
			}

		case string:
		case bool:
		default:
			triggerResult.Err = fmt.Errorf("Unknown type %T for job %d\n", v, job.ID)
		}

		resultSet = append(resultSet, triggerResult)
	}

	e.logger.Printf("Job %d finished with output %s\n", job.ID, strings.TrimSpace(string(out)))
}

func (e *Executor) handleAsInt(job model.Job, t model.Trigger, value float64, triggerResult *TriggerResult) {
	target, err := strconv.Atoi(t.Target)
	if err != nil {
		triggerResult.Err = err
	}
	actual := int(value)

	switch t.TriggerType {
	case model.TriggerEqual:
		if target != actual {
			triggerResult.Err = fmt.Errorf("Expected %d but got %d for job %s.%s", target, actual, job.Name, t.FieldName)
		}
	case model.TriggerMoreThan:
		if actual > target {
			triggerResult.Err = fmt.Errorf("TriggerMoreThan %d was more than %d for job %s.%s", actual, target, job.Name, t.FieldName)
		}
	default:
		triggerResult.Err = fmt.Errorf("Unknown TriggerType")
	}
}

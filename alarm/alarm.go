package alarm

import (
	"fmt"

	"github.com/jvikstedt/alarmy/schedule"
)

type Executor struct {
}

func (e *Executor) Execute(entryID schedule.EntryID) {
	fmt.Println("Hello")
}

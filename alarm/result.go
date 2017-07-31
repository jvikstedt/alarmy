package alarm

import "github.com/jvikstedt/alarmy/model"

type TriggerResult struct {
	Trigger model.Trigger
	Err     error
}

func (t TriggerResult) HasError() bool {
	return t.Err != nil
}

type TriggerResultSet []TriggerResult

func (trs TriggerResultSet) HasErrors() bool {
	for _, v := range trs {
		if v.HasError() {
			return true
		}
	}
	return false
}

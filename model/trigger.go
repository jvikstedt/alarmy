package model

const (
	TriggerEqual    = 0
	TriggerLessThan = 1
	TriggerMoreThan = 2
)

type TriggerType int

var TriggerTypes = []string{
	"Equal",
	"Less Than",
	"More Than",
}

type Trigger struct {
	FieldName string
	Target    string
	TriggerType
}

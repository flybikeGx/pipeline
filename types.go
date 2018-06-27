package pipeline

import (
	"time"
)

type Limit struct {
	timeout time.Duration
}

type Step struct {
	fs     func(interface{}) interface{}
	limits Limit
}

func NewStep(f func(interface{}) interface{}, limits Limit) *Step {
	step := &Step{f, limits}
	return step
}

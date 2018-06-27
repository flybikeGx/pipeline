package pipeline

import (
	"fmt"

	"github.com/flybikeGx/easy-timeout/timelimit"
)

type Pipeline struct {
	steps  []*Step
	chs    []chan interface{}
	errChs []chan error
	run    bool
}

func New() *Pipeline {
	return &Pipeline{
		steps: make([]*Step, 0),
		chs:   make([]chan interface{}, 0),
		run:   true,
	}
}
func (p *Pipeline) Add(step *Step) int {
	p.steps = append(p.steps, step)
	return len(p.steps) - 1
}
func (p *Pipeline) Step(i int) *Step {
	return p.steps[i]
}
func (p *Pipeline) Remove(i int) *Step {
	sp := p.steps[i]
	sp1 := p.steps[:i]
	if i < len(p.steps)-1 {
		sp2 := p.steps[i+1:]
		p.steps = append(sp1, sp2...)
	} else {
		p.steps = sp1
	}
	return sp
}
func (p *Pipeline) Run() {
	//fmt.Println("steps:", len(p.steps))
	p.run = true
	p.chs = append(p.chs, make(chan interface{}))
	p.errChs = append(p.errChs, make(chan error))
	for i := 0; i < len(p.steps); i++ {
		p.chs = append(p.chs, make(chan interface{}))
		p.errChs = append(p.errChs, make(chan error))
		//fmt.Println("len chs", len(p.chs))
		go p.loop(i)
	}
}

func (p *Pipeline) loop(j int) {
	for p.run {
		//fmt.Println("step :", j, "is running")
		select {
		case err := <-p.errChs[j]:
			p.errChs[j+1] <- err
		case arg := <-p.chs[j]:
			var rtn interface{}
			var err error = nil
			ok := timelimit.Run(p.steps[j].limits.timeout, func() {
				defer func() {
					if err0 := recover(); err0 != nil {
						err = err0.(error)
					}
				}()
				rtn = p.steps[j].fs(arg)
			})
			switch {
			case ok && err == nil:
				p.chs[j+1] <- rtn
			case err != nil:
				p.errChs[j+1] <- err
			case !ok:
				p.errChs[j+1] <- fmt.Errorf("timeout: %d", j)
			}
		}
	}
}
func (p *Pipeline) InputChan() chan interface{} {
	return p.chs[0]
}
func (p *Pipeline) OutputChan() chan interface{} {
	return p.chs[len(p.chs)-1]
}
func (p *Pipeline) ErrorChan() chan error {
	return p.errChs[len(p.errChs)-1]
}
func (p *Pipeline) Stop() {
	p.run = false
}

package pipeline

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestPipeline(t *testing.T) {
	pipe := New()
	type input = string
	type output1 struct {
		i int
	}
	type output2 struct {
		b bool
	}
	pipe.Add(&Step{func(i interface{}) interface{} {
		rtn := output1{}
		var err error
		rtn.i, err = strconv.Atoi(i.(input))
		if err != nil {
			panic(err)
		}
		return rtn
	}, Limit{timeout: 200 * time.Millisecond}})

	pipe.Add(&Step{func(i interface{}) interface{} {
		rtn := output2{}
		rtn.b = i.(output1).i < 10
		return rtn
	}, Limit{timeout: 200 * time.Millisecond}})
	pipe.Run()
	chin := pipe.InputChan()
	chout := pipe.OutputChan()
	cherr := pipe.ErrorChan()

	go func() {
		for {
			select {
			case o := <-chout:
				fmt.Println("result:", o.(output2).b)
			case e := <-cherr:
				fmt.Println("result:", e.Error())
			}
		}
	}()

	go func() {
		chin <- "hello"
		//fmt.Println("send hello")
		chin <- "20"
		//fmt.Println("send 20")
		chin <- "5"
		//fmt.Println("send 5")

	}()
	time.Sleep(1 * time.Second)
}

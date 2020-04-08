package spruce

import (
	"time"
)

type clockTask struct {
	c []*clock
	e int64
}
type clock struct {
	Name string
	Fn   func()
}

func NewClockTask(e int64) *clockTask {
	return &clockTask{e: e, c: make([]*clock, 0)}
}
func (c *clockTask) NewClock(name string, fn func()) {
	c.c = append(c.c, &clock{
		Name: name,
		Fn:   fn,
	})

}
func (c *clockTask) Start() {
	for {
		for _, v := range c.c {
			if v != nil {
				v.Fn()
			}
		}
		time.Sleep(time.Second * time.Duration(c.e))
	}
}

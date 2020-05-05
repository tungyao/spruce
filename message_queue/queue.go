package message_queue

import (
	"fmt"
)

type Body struct {
	Next    *Body
	Key     []byte
	Value   interface{}
	Check   bool
	Operate func() error
}
type Queue struct {
	Body   *Body
	Length int
}

func InitQueue() *Queue {
	return &Queue{}
}
func (q *Queue) InQueue(body *Body) {

	if q.Body == nil {
		q.Body = body
		return
	}
	//var x = q.Body.Next
	for q.Body.Next != nil && q.Body.Next.Check == true {
		fmt.Println(q.Body.Next)
		q.Body.Next = q.Body.Next.Next
	}
	q.Length += 1
	q.Body.Next = body

}
func (q *Queue) OutQueue() *Body {
	var get = q.Body
	q.Body = q.Body.Next
	q.Length -= 1
	return get
}
func FrontQueue() {

}
func EmptyQueue() {

}

package test

import (
	"testing"
)
import "../../spruce"

func TestRing(t *testing.T) {
	spruce.CreateHash(1024)
	r := spruce.NewMessage()
	//for {
	r.Push(r.MSG("a", "b", "now", "23131"))
	//}
}

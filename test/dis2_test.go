package test

import (
	"../../spruce"
	"fmt"
	"testing"
)

func TestDIS1(t *testing.T) {
	//spruce.StartSpruceDistributed(spruce.Config{Test: ":80"})
	a := spruce.CreateHash(512)
	a.Set("hello", "world", 0)
	a.Set("helloa", "worlda", 0)
	a.Set("hellob", "worldb", 0)
	fmt.Print(a.Get("*"))
}

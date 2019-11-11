package test

import (
	"../../spruce"
	"testing"
)

func TestDIS1(t *testing.T) {
	//spruce.StartSpruceDistributed(spruce.Config{Test: ":80"})
	t.Log(spruce.New().Get("hello"))
}

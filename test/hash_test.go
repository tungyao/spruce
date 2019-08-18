package test

import (
	"../core/shash"
	"testing"
)

func TestHash(t *testing.T) {

	weights := make(map[string]int)
	weights["192.168.0.246:11212"] = 1
	weights["192.168.0.247:11212"] = 2
	weights["192.168.0.249:11212"] = 1
	ring := shash.NewWithWeights(weights)
	s, _ := ring.GetNode("key")
	t.Log(s)
}

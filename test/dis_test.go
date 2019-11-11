package test

import (
	"../../spruce"
	"testing"
)

func TestDIS1(t *testing.T) {
	d := spruce.StartSpruceDistributed(spruce.Config{Test: ":80"})
	t.Log(d.Get("hello"))
}
func TestDIS2(t *testing.T) {
	spruce.StartSpruceDistributed(spruce.Config{Test: ":81"})
}
func TestDIS3(t *testing.T) {
	spruce.StartSpruceDistributed(spruce.Config{Test: ":83"})

}

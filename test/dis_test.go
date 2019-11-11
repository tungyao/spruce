package test

import (
	"../../spruce"
	"testing"
)

func TestDIS2(t *testing.T) {
	spruce.StartSpruceDistributed(spruce.Config{Test: ":81"})
}
func TestDIS3(t *testing.T) {
	spruce.StartSpruceDistributed(spruce.Config{Test: ":83"})
}

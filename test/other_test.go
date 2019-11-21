package test

import (
	"testing"
)
import "../../spruce"

func TestOther(t *testing.T) {
	s := spruce.Encrypt([]byte("hello wolrd"))
	t.Log(string(spruce.Decrypt(s)))
}

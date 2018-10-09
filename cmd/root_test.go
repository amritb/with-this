package cmd

import (
	"testing"
)

// Dummy test
func Test_Execute(t *testing.T) {
	toRun := "echo with-this"
	prog, _, res, err := execute(toRun)

	if prog != "echo" {
		t.Error("wrong prog detected")
	}

	if string(res[:]) != "with-this\n" {
		t.Error("wrong result")
	}

	if err != nil {
		t.Error(err)
	}
}
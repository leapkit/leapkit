package main_test

import (
	"os/exec"
	"testing"
)

func TestVersion(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
	}

	expected := "Kit version: (devel)\n"
	if string(out) != expected {
		t.Errorf("Expected %v, got %v", expected, string(out))
	}
}

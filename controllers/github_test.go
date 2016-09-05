package controllers

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestExecShell(t *testing.T) {
	rel, _ := os.Getwd()
	cfg := map[string]string{
		"sh":   path.Join(filepath.Dir(rel), "shell/gitpull.sh"),
		"work": filepath.Dir(rel),
	}
	if err := execShell(cfg); err != nil {
		t.Error(err)
	}
}

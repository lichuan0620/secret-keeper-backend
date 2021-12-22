package service

import "testing"

func TestBuild(t *testing.T) {
	if _, err := Build(nil); err != nil {
		t.Error(err)
	}
}

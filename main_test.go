package main

import "testing"

// Test for NewMsgHistory
func TestNewMsgHistory(t *testing.T) {
	mh := NewMsgHistory()
	if mh == nil {
		t.Errorf("NewMsgHistory() failed")
	}
}

package main

import "testing"

func TestInsertMatch(t *testing.T) {
	f := PrefixFilter{root: newNode()}
	consumer := &Consumer{}
	f.Insert("user", consumer)
	if len(f.Match("user.create")) == 0 {
		t.Errorf("Match not found")
	}
}

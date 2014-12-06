package main

// import (
// 	"testing"
// )

// var candidates = []Candidate{
// 	NewCandidate("zero"),
// 	NewCandidate("one"),
// 	NewCandidate("two"),
// 	NewCandidate("three"),
// }

// func moveViewCursor(view *View, dir string) {
// 	for _, d := range dir {
// 		switch d {
// 		case 'd':
// 			view.Down()
// 		case 'u':
// 			view.Up()
// 		}
// 	}
// }

// var testsIndex = []struct {
// 	height   int
// 	expected int
// 	dir      string
// }{
// 	{height: 3, expected: 0, dir: ""},
// 	{height: 3, expected: 1, dir: "d"},
// 	{height: 1, expected: 0, dir: "d"},
// 	{height: 3, expected: 0, dir: "du"},
// 	{height: 3, expected: 0, dir: "u"},
// }

// func TestPickerViewIndex(t *testing.T) {

// 	for _, test := range testsIndex {
// 		view := NewPicker(candidates, test.height).Answer("")
// 		moveViewCursor(view, test.dir)
// 		if got := view.Index(); got != test.expected {
// 			t.Errorf("Expected index %v, got %v", test.height, test.dir, test.expected, got)
// 		}

// 	}
// }

// var testsSelected = []struct {
// 	query    string
// 	dir      string
// 	expected string
// }{
// 	{query: "two", dir: "", expected: "two"},
// 	{query: "two", dir: "dd", expected: "two"},
// 	{query: "blah", dir: "", expected: ""},
// }

// func TestPickerViewSelected(t *testing.T) {
// 	for _, test := range testsSelected {
// 		view := NewPicker(candidates, 3).Answer(test.query)
// 		moveViewCursor(view, test.dir)
// 		if s := view.Selected(); s != test.expected {
// 			t.Errorf("Expected %q, got %q", test.expected, s)
// 		}
// 	}
// }

// func TestPickerViewHandleWhenCandidatesAreFewerThanHeight(t *testing.T) {
// 	view := NewPicker(candidates, 5).Answer("")
// 	if lr := len(view.Rows); lr != 4 {
// 		t.Errorf("Expected %v rows, got", lr)
// 	}
// 	if view.Height != 5 {
// 		t.Errorf("Expected %v rows, got", view.Height)
// 	}
// }

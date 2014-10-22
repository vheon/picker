package main

import (
	. "github.com/onsi/gomega"
	"testing"
)

var candidates = []Candidate{
	NewCandidate("zero"),
	NewCandidate("one"),
	NewCandidate("two"),
	NewCandidate("three"),
}

func TestPicker_SelectFirstCandidateByDefault(t *testing.T) {
	RegisterTestingT(t)

	view := NewPicker(candidates, 3).Answer("")
	Expect(view.Index()).To(Equal(0))
}

func TestPicker_CanSelectCandidateDown(t *testing.T) {
	RegisterTestingT(t)

	view := NewPicker(candidates, 3).Answer("")
	view.Down()
	Expect(view.Index()).To(Equal(1))
}

func TestPicker_DoNotSelectCandidatesDownOverTheEdge(t *testing.T) {
	RegisterTestingT(t)

	view := NewPicker(candidates, 1).Answer("")
	view.Down()
	Expect(view.Index()).To(Equal(0))
}

func TestPicker_CanSelectCandidateUp(t *testing.T) {
	RegisterTestingT(t)

	view := NewPicker(candidates, 3).Answer("")
	view.Down().Up()
	Expect(view.Index()).To(Equal(0))
}

func TestPicker_DoNotSelectCandidatesUpOverTheEdge(t *testing.T) {
	RegisterTestingT(t)

	view := NewPicker(candidates, 3).Answer("")
	view.Up()
	Expect(view.Index()).To(Equal(0))
}

func TestPicker_SortTheRightAnswerForAQuery(t *testing.T) {
	RegisterTestingT(t)
	view := NewPicker(candidates, 3).Answer("two")
	Expect(view.Selected()).To(Equal("two"))
}

func TestPicker_DontShowCandidatesWithScoreZero(t *testing.T) {
	RegisterTestingT(t)

	view := NewPicker(candidates, 3).Answer("two")
	view.Down().Down()
	Expect(view.Selected()).To(Equal("two"))
}

func TestPicker_ReturnAValidViewWhenNoGoodCandidatesAreAvailable(t *testing.T) {
	RegisterTestingT(t)

	view := NewPicker(candidates, 3).Answer("blah")
	Expect(view).To(Equal(&View{height: 3, lines: []string{}, query: "blah"}))
}

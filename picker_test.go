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
	view.Down()
	view.Up()
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
	view.Down()
	view.Down()
	Expect(view.Selected()).To(Equal("two"))
}

func TestPicker_ReturnAValidViewWhenNoGoodCandidatesAreAvailable(t *testing.T) {
	RegisterTestingT(t)

	view := NewPicker(candidates, 3).Answer("blah")
	Expect(view).To(Equal(&View{Height: 3, Rows: []string{}, Query: "blah", prompt: "> "}))
}

func TestPicker_HandleWhenCandidatesAreFewerThanHeight(t *testing.T) {
	RegisterTestingT(t)

	view := NewPicker(candidates, 5).Answer("")
	Expect(len(view.Rows)).To(Equal(4))
	Expect(view.Height).To(Equal(5))
}

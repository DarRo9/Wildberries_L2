package main

import (
	"math"
	"reflect"
	"testing"
)

func TestProcessingFieldsCorrectInput(t *testing.T) {
	s := "-2,5,7-9,12-"
	expected := [][]int{{1, 2}, {5, 5}, {7, 9}, {12, math.MaxInt}}

	parsed, err := processingFields(s)

	if err != nil || !reflect.DeepEqual(parsed, expected) {
		t.Logf("processingFields(%s) = %v, %v, expected: %v, nil", s, parsed, err, expected)
	}
}

func TestProcessingFieldsIncorrectInputDecreasingRange(t *testing.T) {
	s := "5-2"

	parsed, err := processingFields(s)

	if err == nil {
		t.Logf("processingFields(%s) = %v, nil, expected: nil, error", s, parsed)
	}
}

func TestProcessingFieldsIncorrectInputFormat(t *testing.T) {
	s := "hello"

	parsed, err := processingFields(s)

	if err == nil {
		t.Logf("processingFields(%s) = %v, nil, expected: nil, error", s, parsed)
	}
}

func TestProcessingFieldsIncorrectInputComma(t *testing.T) {
	s := ","

	parsed, err := processingFields(s)

	if err == nil {
		t.Logf("processingFields(%s) = %v, nil, expected: nil, error", s, parsed)
	}
}

func TestProcessingFieldsIncorrectInputPreComma(t *testing.T) {
	s := ",2"

	parsed, err := processingFields(s)

	if err == nil {
		t.Logf("processingFields(%s) = %v, nil, expected: nil, error", s, parsed)
	}
}

func TestProcessingFieldsIncorrectInputPostComma(t *testing.T) {
	s := "2,"

	parsed, err := processingFields(s)

	if err == nil {
		t.Logf("processingFields(%s) = %v, nil, expected: nil, error", s, parsed)
	}
}

func TestindexInSectionTrue(t *testing.T) {
	i := 8
	segs := [][]int{{7, 9}}

	res := indexInSection(i, segs)

	if !res {
		t.Logf("indexInSection(%d, %v) = %v, expected: true", i, segs, res)
	}
}

func TestindexInSectionFalse(t *testing.T) {
	i := 5
	segs := [][]int{{7, 9}}

	res := indexInSection(i, segs)

	if res {
		t.Logf("indexInSection(%d, %v) = %v, expected: true", i, segs, res)
	}
}
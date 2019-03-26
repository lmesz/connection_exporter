package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringInStrings(t *testing.T) {
	assert.False(t, StringInSlice("appletree", []string{"peartree", "lemontree"}))
	assert.True(t, StringInSlice("appletree", []string{"appletree", "peartree", "lemontree"}))
}

func TestParseNetstatsResult(t *testing.T) {
	ParseNetstatsResult([]string{"tcp        0      0 127.0.0.1:44076         127.0.0.1:9911          TIME_WAIT   -"})
	//expectedResult := map[string]{"almafa"}map[string]float64{"almafa": 53.0}
	//assert.Equal(t, result, expectedResult)
}

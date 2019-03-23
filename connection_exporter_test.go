package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringInStrings(t *testing.T) {
	assert.False(t, StringInSlice("appletree", []string{"peartree", "lemontree"}))
	assert.True(t, StringInSlice("appletree", []string{"appletree", "peartree", "lemontree"}))
}

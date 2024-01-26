package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTerminal(t *testing.T) {
	nvim := New()
	assert.NotNil(t, nvim)
}

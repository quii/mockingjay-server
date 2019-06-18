package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsURL(t *testing.T) {
	t.Run("url", func(t *testing.T) {
		assert.True(t, isURL("http://quii.dev"))
	})

	t.Run("some file path", func(t *testing.T) {
		assert.False(t, isURL("examples/example.yaml"))
	})

	t.Run("filepath with a leading slash", func(t *testing.T) {
		assert.False(t, isURL("/fakes/example.yaml"))
	})
}

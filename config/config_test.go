package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	args := []string{
		"watcher", "./incoming",
	}

	output, err := Parse(args)
	assert.NoError(t, err)

	assert.Equal(t, output.Path, "./incoming")
	assert.Equal(t, output.Interval, 3*time.Second)
	assert.False(t, output.Recursive)
	assert.Equal(t, output.OutputFormat, Console)
}

func TestCustomConfig(t *testing.T) {
	args := []string{
		"watcher", "-output=json", "-interval=5s", "./test-path",
	}

	output, err := Parse(args)
	assert.NoError(t, err)

	assert.Equal(t, output.Path, "./test-path")
	assert.Equal(t, output.Interval, 5*time.Second)
	assert.False(t, output.Recursive)
	assert.Equal(t, output.OutputFormat, JSON)
}

func TestCutomConfig_WithoutEqualSign(t *testing.T) {
	args := []string{
		"watcher", "-output", "json", "-interval", "5s", "./test-path",
	}

	output, err := Parse(args)
	assert.NoError(t, err)

	assert.Equal(t, output.Path, "./test-path")
	assert.Equal(t, output.Interval, 5*time.Second)
	assert.False(t, output.Recursive)
	assert.Equal(t, output.OutputFormat, JSON)
}

func TestCustomConfig_Mixed(t *testing.T) {
	args := []string{
		"watcher", "-output", "json", "-interval=5s", "--recursive", "./test-path",
	}

	output, err := Parse(args)
	assert.NoError(t, err)

	assert.Equal(t, output.Path, "./test-path")
	assert.Equal(t, output.Interval, 5*time.Second)
	assert.True(t, output.Recursive)
	assert.Equal(t, output.OutputFormat, JSON)
}

func TestConfigError_WrongArg(t *testing.T) {
	args := []string{
		"watcher", "-output", "yaml", "-interval=5s", "--recursive", "./test-path",
	}

	output, err := Parse(args)
	assert.Error(t, err)
	assert.Nil(t, output)
}

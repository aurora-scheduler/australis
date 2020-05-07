package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJob(t *testing.T) {
	_, err := UnmarshalJob("../test/hello_world.yaml")
	assert.NoError(t, err)
}

func TestUnmarshalUpdate(t *testing.T) {
	_, err := UnmarshalUpdate("../test/update_hello_world.yaml")
	assert.NoError(t, err)
}


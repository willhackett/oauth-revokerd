package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DefaultEnvironment(t *testing.T) {
	assert := assert.New(t)

	// This is unlikely to be useful but it's a valid configuration
	config, err := Load()

	assert.NoError(err)
	assert.Equal(8080, config.Port)
}

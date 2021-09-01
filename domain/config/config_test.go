package config

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GenerateFinalTargetPath(t *testing.T) {
	config := Config{
		Path:     "/.*Blub",
		Target:   "http://targetUrl",
		Response: Response{},
		Id:       uuid.UUID{},
	}

	result := config.GenerateFinalTargetPath("/bliBlub")
	assert.Equal(t, "http://targetUrl/bliBlub", result)
}

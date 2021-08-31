package config

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GenerateFinalTargetPath(t *testing.T) {
	config := Config{
		Path:     "/.*Blub",
		Target:   "/Bla",
		Response: Response{},
		Id:       uuid.UUID{},
	}

	result := config.GenerateFinalTargetPath("/bliBlub")
	assert.Equal(t, "/Bla/bliBlub", result)
}

package config

import "github.com/google/uuid"

type Repository interface {
	Save(config Config) error
	Find(id uuid.UUID) *Config
	FindAll() []Config
	Delete(id uuid.UUID) error
	DeleteAll() error
}

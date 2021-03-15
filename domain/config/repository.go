package config

type Repository interface {
	Save(config Config) error
	Find(path string) *Config
	FindAll() []Config
	Delete(path string) error
	DeleteAll() error
}

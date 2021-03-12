package config

type Repository interface {
	Save(config Config) error
	Find(path string) *Config
	Delete(path string) error
}

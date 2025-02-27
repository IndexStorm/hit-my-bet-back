package config

type Environment int

const (
	EnvironmentLocal Environment = 0
	EnvironmentStage Environment = 1
)

type DefaultEnvironment struct {
	Value Environment `env:"ENVIRONMENT,notEmpty" envDefault:"0"`
}

package config

type Environment string

const (
	EnvironmentProduction Environment = "prod"
)

func GetEnvironmentName() Environment {
	return EnvironmentProduction
}

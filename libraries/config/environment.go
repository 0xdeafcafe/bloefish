package config

type Environment string

const (
	EnvironmentProduction Environment = "prod"
	EnvironmentUnitTest   Environment = "unittest"
)

func GetEnvironmentName() Environment {
	return EnvironmentProduction
}

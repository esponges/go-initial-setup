package configuration

type Configuration struct {
	Port string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port: "8080",
	}
}

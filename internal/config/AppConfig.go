package config

type AppConfig struct {
	Host string
	Port int
}

func NewConfig() (AppConfig, error) {

	cfg := AppConfig{Host: "0.0.0.0", Port: 7002}
	return cfg, nil

}

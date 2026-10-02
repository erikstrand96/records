package config

type AppConfig struct {
	Host string
	Port string
}

func NewConfig(host string, port string) (AppConfig, error) {

	cfg := AppConfig{Host: host, Port: port}
	return cfg, nil

}

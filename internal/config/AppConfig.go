package config

type AppConfig struct {
	Host   string
	Port   string
	DBPath string
}

func NewConfig(host, port, dbPath string) (AppConfig, error) {

	cfg := AppConfig{Host: host, Port: port, DBPath: dbPath}
	return cfg, nil

}

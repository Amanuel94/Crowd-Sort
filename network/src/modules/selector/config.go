package selector

type Config struct {
	alg string
}

func NewConfig() *Config {
	return &Config{
		alg: "BEMS",
	}
}

func (cfg *Config) WithAlg(algName string) {
	cfg.alg = algName
}

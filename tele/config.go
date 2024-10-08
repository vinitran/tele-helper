package tele

type Config struct {
	Name  string
	Proxy string
}

func (cfg Config) isValid() bool {
	if cfg.Name == "" {
		return false
	}
	return true
}

func (cfg Config) useProxy() bool {
	if cfg.Proxy == "" {
		return false
	}
	return true
}

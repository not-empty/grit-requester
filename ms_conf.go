package gritrequester

import "fmt"

type MSAuthConf struct {
	Token   string
	Secret  string
	Context string
	BaseUrl string
}

type ConfigProvider interface {
	Get(service string) (MSAuthConf, error)
}

type StaticConfig map[string]MSAuthConf

func (s StaticConfig) Set(service string, conf MSAuthConf) {
	s[service] = conf
}

func (s StaticConfig) Get(service string) (MSAuthConf, error) {
	var conf MSAuthConf
	if len(s) == 0 {
		return conf, fmt.Errorf("config map is empty")
	}

	conf, ok := s[service]
	if !ok {
		return conf, fmt.Errorf("config not found for service: %s", service)
	}
	return conf, nil
}

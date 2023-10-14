package util

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Conf struct {
	ProjectName    string `yaml:"project_name"`
	ProjectAbsPath string `yaml:"project_abs_path"`
	Config         []config
	ProjectMode    string `yaml:"project_mode"`
}

type config struct {
	Name           string `yaml:"name"`
	RegisteredFile string `yaml:"registered_file"`
	Path           string `yaml:"path"`
	ParentName     string `yaml:"parent_name"`
}

// InitConf 初识化配置
func InitConf(path string) (*Conf, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := new(Conf)
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// InitConf2 初识化配置，通过配置数据
func InitConf2(data string) (*Conf, error) {
	cfg := new(Conf)
	if err := yaml.Unmarshal([]byte(data), cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

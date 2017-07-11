package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	MachineId int64 `json:"machine_id"`
	CentorId  int64 `json:"centor_id"`
}

func ParseConfigFile(filename string) (*Config, error) {
	cfg := new(Config)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

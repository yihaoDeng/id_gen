package config

import (
	"testing"
)

func Test_ParseConfigFile(t *testing.T) {
	cfg, err := ParseConfigFile("test.json")
	if err != nil {
		t.Fail()
		t.Error("error")
		return
	}
	t.Log(cfg.MachineId, "\t", cfg.CentorId)

}

package config

import "testing"

func TestLoad(t *testing.T) {
	conf := Load()
	if conf.Database.Host == "" {
		t.Error("Expected Database.Host to be set, got empty string")
	}
	if conf.Database.Port == 0 {
		t.Error("Expected Database.Port to be set, got 0")
	}
	if conf.Database.User == "" {
		t.Error("Expected Database.User to be set, got empty string")
	}
	if conf.Database.Password == "" {
		t.Error("Expected Database.Password to be set, got empty string")
	}
	if conf.Database.DbName == "" {
		t.Error("Expected Database.DbName to be set, got empty string")
	}
	if conf.Database.SslMode == "" {
		t.Error("Expected Database.SslMode to be set, got empty string")
	}
}

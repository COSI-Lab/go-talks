package main

import (
	"io"
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	// Addresses & ports to listen on
	Listen []string `toml:"listen"`
	// Meeting password
	Password string `toml:"password"`
	// Database file name
	Database string `toml:"database"`
	// Trusted subnets
	Subnets []string `toml:"subnets"`
}

func ParseConfig(r io.Reader) Config {
	var config Config
	_, err := toml.NewDecoder(r).Decode(&config)
	if err != nil {
		log.Fatal("[FATAL] Could not parse config file:", err)
	}
	return config
}

func DefaultConfig() Config {
	return Config{
		Listen:   []string{"localhost:5000"},
		Password: "conway",
		Database: "talks.db",
		Subnets:  []string{"128.153.144.0/23", "2605:6480:c051::1/48"},
	}
}

func (c *Config) Validate() {
	if len(c.Listen) == 0 {
		log.Fatal("[FATAL] No listen addresses specified")
	}
	if c.Password == "" {
		log.Fatal("[FATAL] No password specified")
	}
	if c.Database == "" {
		log.Fatal("[FATAL] No database specified")
	}
}

func (c *Config) Network() Networks {
	return NewNetworks(c.Subnets)
}

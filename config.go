package main

import (
	"io"
	"log"

	"github.com/BurntSushi/toml"
)

// Config is the configuration for go-talks
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

// ParseConfig parses a TOML config file and returns a Config struct
func ParseConfig(r io.Reader) Config {
	var config Config
	_, err := toml.NewDecoder(r).Decode(&config)
	if err != nil {
		log.Fatal("[FATAL] Could not parse config file:", err)
	}
	return config
}

// Validate checks a Config for some common errors
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

// Network creates a Networks struct from the []string
// of subnets in the Config
func (c *Config) Network() Networks {
	return NewNetworks(c.Subnets)
}

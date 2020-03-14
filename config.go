package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type configInternal struct {
	Listen   string  `json:"listen"`
	User     string  `json:"user"`
	Password string  `json:"password"`
	Timeout  float64 `json:"timeout"`
}

type Config struct {
	Listen   string
	User     string
	Password string
	Timeout  time.Duration
}

func NewConfig(filename string) *Config {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	config := configInternal{
		Listen:  ":9091",
		Timeout: 10,
	}
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatalf("Error: %s", err)
	}
	return config.export()
}

func (c *configInternal) export() *Config {
	return &Config{
		Listen:   c.Listen,
		User:     c.User,
		Password: c.Password,
		Timeout:  time.Duration(c.Timeout * float64(time.Second)),
	}
}

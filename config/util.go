package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Write to config file.
func Write(filepath string, config Config) error {
	json, err := json.Marshal(config)
	if err != nil {
		fmt.Println("fail to create a json")
		fmt.Println(err)
		return err
	}
	if err := os.WriteFile(filepath, json, os.ModePerm); err != nil {
		fmt.Println("fail to create the config file")
		fmt.Println(err)
		return err
	}
	return nil
}

func Read(filepath string) (*Config, error) {
	con, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	cfg := Config{}
	if err := json.Unmarshal(con, &cfg); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &cfg, nil
}
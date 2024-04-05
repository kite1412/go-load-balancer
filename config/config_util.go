package config

import (
	"os"

	"my.go/load-balancer/lberror"
)

// return the absolute path of the lb config file.
func LBConfigAbs() (string, error) {
	absPath := os.Getenv("lb-config")
	if absPath == "" {
		return "", lberror.ConfigFileError("can't find the absolute path of the lb config file.")
	}	
	if _, err := os.Stat(absPath); err != nil {
		return "", lberror.ConfigFileError("config file is not exist.")
	}
	return absPath, nil
}
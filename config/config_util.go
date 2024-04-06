package config

import (
	"os"
	"strconv"
	"strings"

	"my.go/loadbalancer/lberror"
)

// return the absolute path of the lb config file.
// the env variable which this func look up to is 'lb-config'.
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

// return the port lower and upper limit.
// the env variable which this func look up to is 'll-ul'
// with value being lower limit and upper limit splitted with comma,
// e.g. "'lowerlimit','upperlimit'".
func LBLowUpLimitPort() (ll, ul int, err error) {
	llul := os.Getenv("ll-ul")
	if llul == "" {
		return -1, -1, lberror.ConfigFileError("can't find port lower and upper limit")
	}
	s := strings.Split(llul, ",")
	ll, lErr := strconv.Atoi(s[0])
	ul, uErr := strconv.Atoi(s[1])

	if lErr != nil || uErr != nil {
		return -1, -1, lberror.ConfigFileError("can't parse the given variable")
	}

	return ll, ul, nil
}
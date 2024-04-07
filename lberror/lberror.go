// App-specific errors.
package lberror

// Error when generating port.
type generatePortError struct {
	m string
}

func (e generatePortError) Error() string {
	return e.m
}

func GeneratePortError(message string) error {
	return generatePortError{m: message}
}

// Error when trying to access lb config file.
type configFileError struct {
	m string
}

func (e configFileError) Error() string {
	return e.m
}

func ConfigFileError(message string) error {
	return configFileError{m: message}
}

// Error when no backend available to handle request.
type noBackendFoundError struct {
	m string
}

func (e noBackendFoundError) Error() string {
	return e.m
}

func NoBackendFoundError(message string) error {
	return noBackendFoundError{m: message}
}
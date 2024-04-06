// app-specifif errors.
package lberror

// error when generating port.
type generatePortError struct {
	m string
}

func (e generatePortError) Error() string {
	return e.m
}

func GeneratePortError(message string) error {
	return generatePortError{m: message}
}

// error when trying to access lb config file.
type configFileError struct {
	m string
}

func (e configFileError) Error() string {
	return e.m
}

func ConfigFileError(message string) error {
	return configFileError{m: message}
}

// error when no backend available to handle request.
type noBackendFoundError struct {
	m string
}

func (e noBackendFoundError) Error() string {
	return e.m
}

func NoBackendFoundError(message string) error {
	return noBackendFoundError{m: message}
}

// error when accessing nil-value state.
type stateNotInstantiatedError struct {
	m string
}

func (e stateNotInstantiatedError) Error() string {
	return e.m
}

func StateNotInstantiatedError(message string) error {
	return stateNotInstantiatedError{m: message}
}
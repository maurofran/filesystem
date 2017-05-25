package filesystem

import "fmt"

// PluginError is the error for plugins
type PluginError interface {
	error
	Method() string
}

type pluginError struct {
	message string
	method  string
}

func (e pluginError) Method() string {
	return e.method
}

func (e pluginError) Error() string {
	return fmt.Sprintf(e.message, e.method)
}

// IsPluginError will check if provided error is a plugin error.
func IsPluginError(err error) bool {
	_, ok := err.(PluginError)
	return ok
}

func pluginNotFoundError(method string) PluginError {
	return &pluginError{message: "No plugin found for method %s", method: method}
}

// PathError is the error if provided path is not valid.
type PathError interface {
	error
	Path() Path
}

type pathError struct {
	message string
	path    Path
}

// Path is the path of error.
func (e pathError) Path() Path {
	return e.path
}

func (e pathError) Error() string {
	return fmt.Sprintf(e.message, e.path)
}

// IsPathError will check if provided error is an invalid path error.
func IsPathError(err error) bool {
	_, ok := err.(PathError)
	return ok
}

func invalidPathError(path Path) PathError {
	return pathError{"Path %s as an invalid prefix", path}
}

// MountError is the error returned when a mount already exists.
type MountError interface {
	error
	Prefix() string
}

type mountError struct {
	message string
	prefix  string
}

// Prefix is the prefix for which the mount already exists.
func (e mountError) Prefix() string {
	return e.prefix
}

func (e mountError) Error() string {
	return fmt.Sprintf(e.message, e.prefix)
}

// IsMountError will check if provided error is a mount exists error.
func IsMountError(err error) bool {
	_, ok := err.(MountError)
	return ok
}

func mountNotFoundError(prefix string) MountError {
	return mountError{"Mount prefix %s does not exists", prefix}
}

func mountExistsError(prefix string) MountError {
	return mountError{"Mount prefix %s already exists", prefix}
}

// FileNotFoundError is the error raised when a file was not found.
type FileNotFoundError interface {
	Path() Path
}

// IsFileNotFound will check if file is not found
func IsFileNotFound(err error) bool {
	_, ok := err.(FileNotFoundError)
	return ok
}

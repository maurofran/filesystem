package plugins

import (
	"errors"

	"github.com/maurofran/filesystem"
)

// ForceCopy will copy a file, forcing overwriting.
type ForceCopy struct {
	plugin
}

// Method is the name of the method to be used to invoke the plugin.
func (p *ForceCopy) Method() string {
	return "ForceCopy"
}

// Handle the invocation of plugin.
func (p *ForceCopy) Handle(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return false, errors.New("path and newpath arguments are required")
	}
	path, ok := args[0].(filesystem.Path)
	if !ok {
		return false, errors.New("path must be an instance of filesystem.Path")
	}
	newPath, ok := args[1].(filesystem.Path)
	if !ok {
		return false, errors.New("newPath must be an instance of filesystem.Path")
	}
	deleted, err := p.fs.Delete(newPath)
	if err != nil {
		if filesystem.IsFileNotFound(err) {
			deleted = true
		} else {
			return false, err
		}
	}
	if deleted {
		err := p.fs.Copy(path, newPath)
		return true, err
	}
	return false, nil
}
